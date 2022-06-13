package repository

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/clubo-app/relation-service/datastruct"
	"github.com/go-playground/validator/v10"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

const (
	FRIEND_RELATIONS string = "friend_relations"
	FRIEND_COUNT     string = "friend_count"
)

var friendCountMetadata = table.Metadata{
	Name:    FRIEND_COUNT,
	Columns: []string{"user_id", "friend_count"},
	PartKey: []string{"user_id"},
}
var friendRelationMetadata = table.Metadata{
	Name:    FRIEND_RELATIONS,
	Columns: []string{"user_id", "friend_id", "accepted", "requested_at", "accepted_at"},
	PartKey: []string{"user_id", "friend_id"},
}

type FriendRelationRepository interface {
	CreateFriendRequest(ctx context.Context, fr datastruct.FriendRelation) error
	DeclineFriendRequest(ctx context.Context, uId, fId string) error
	AcceptFriendRequest(ctx context.Context, uId, fId string) error
	RemoveFriendRelation(ctx context.Context, uId, fId string) error
	GetFriendRelation(ctx context.Context, uId, fId string) (datastruct.FriendRelation, error)
	GetFriends(ctx context.Context, uId string, page []byte, limit uint64) ([]datastruct.FriendRelation, []byte, error)
	GetIncomingFriendRequests(ctx context.Context, uId string, page []byte, limit uint64) ([]datastruct.FriendRelation, []byte, error)
	GetFriendCount(ctx context.Context, uId string) (datastruct.FriendCount, error)
	GetManyFriendCount(ctx context.Context, ids []string) ([]datastruct.FriendCount, error)
}

type friendRelationRepository struct {
	sess *gocqlx.Session
}

func (r *friendRelationRepository) CreateFriendRequest(ctx context.Context, fr datastruct.FriendRelation) error {
	v := validator.New()
	err := v.Struct(fr)
	if err != nil {
		return err
	}

	stmt, names := qb.
		Insert(FRIEND_RELATIONS).
		Unique().
		Columns(friendRelationMetadata.Columns...).
		ToCql()

	err = r.sess.
		ContextQuery(ctx, stmt, names).
		BindStruct(fr).
		ExecRelease()
	if err != nil {
		return err
	}

	return nil
}

func (r *friendRelationRepository) DeclineFriendRequest(ctx context.Context, uId, fId string) error {
	stmt, names := qb.
		Delete(FRIEND_RELATIONS).
		Where(qb.Eq("user_id")).
		Where(qb.Eq("friend_id")).
		ToCql()

	err := r.sess.ContextQuery(ctx, stmt, names).
		BindMap((qb.M{
			"user_id":   uId,
			"friend_id": fId,
		})).
		ExecRelease()
	if err != nil {
		return err
	}

	return nil
}

func (r *friendRelationRepository) AcceptFriendRequest(ctx context.Context, uId, fId string) error {
	wg := new(sync.WaitGroup)
	wg.Add(3)

	go func() {
		defer wg.Done()

		uStmt, uNames := qb.
			Update(FRIEND_RELATIONS).
			Where(qb.Eq("user_id")).
			Where(qb.Eq("friend_id")).
			If(qb.EqNamed("accepted", "old.accepted")).
			Set("accepted").
			Set("accepted_at").
			ToCql()

		err := r.sess.
			ContextQuery(ctx, uStmt, uNames).
			BindMap((qb.M{
				"user_id":      uId,
				"friend_id":    fId,
				"old.accepted": false,
				"accepted":     true,
				"accepted_at":  time.Now(),
			})).
			ExecRelease()

		log.Println("Update Friend Request Error: ", err)
	}()

	go func() {
		defer wg.Done()

		cStmt, cNames := qb.
			Insert(FRIEND_RELATIONS).
			Columns(friendRelationMetadata.Columns...).
			ToCql()

		err := r.sess.
			ContextQuery(ctx, cStmt, cNames).
			BindMap((qb.M{
				"user_id":      fId,
				"friend_id":    uId,
				"accepted":     true,
				"accepted_at":  time.Now(),
				"requested_at": time.Now(),
			})).
			ExecRelease()

		log.Println("Create Friend Relation Error: ", err)
	}()

	go func() {
		defer wg.Done()

		countStmt, countNames := qb.
			Update(FRIEND_COUNT).
			Where(qb.In("user_id")).
			Add("friend_count").
			ToCql()

		err := r.sess.
			ContextQuery(ctx, countStmt, countNames).
			BindMap((qb.M{
				"friend_count": 1,
				"user_id":      []string{uId, fId},
			})).
			ExecRelease()

		log.Println("Friend Count Addition Error: ", err)
	}()

	wg.Wait()
	return nil
}

func (r *friendRelationRepository) RemoveFriendRelation(ctx context.Context, uId, fId string) error {
	stmt1, names1 := qb.
		Delete(FRIEND_RELATIONS).
		Where(qb.Eq("user_id")).
		Where(qb.Eq("friend_id")).
		ToCql()

	stmt2, names2 := qb.
		Delete(FRIEND_RELATIONS).
		Where(qb.Eq("user_id")).
		Where(qb.Eq("friend_id")).
		ToCql()

	err := r.sess.ContextQuery(ctx, stmt1, names1).
		BindMap((qb.M{
			"user_id":   uId,
			"friend_id": fId,
		})).
		ExecRelease()
	if err != nil {
		return err
	}

	err = r.sess.ContextQuery(ctx, stmt2, names2).
		BindMap((qb.M{
			"user_id":   fId,
			"friend_id": uId,
		})).
		ExecRelease()
	if err != nil {
		return err
	}

	return nil
}

func (r *friendRelationRepository) GetFriendRelation(ctx context.Context, uId, fId string) (res datastruct.FriendRelation, err error) {
	stmt, names := qb.
		Select(FRIEND_RELATIONS).
		Columns(friendRelationMetadata.Columns...).
		Where(qb.Eq("user_id")).
		Where(qb.Eq("friend_id")).
		ToCql()

	err = r.sess.
		ContextQuery(ctx, stmt, names).
		BindMap((qb.M{
			"user_id":   uId,
			"friend_id": fId,
		})).
		GetRelease(&res)
	if err != nil {
		stmt2, names2 := qb.
			Select(FRIEND_RELATIONS).
			Columns(friendRelationMetadata.Columns...).
			Where(qb.Eq("user_id")).
			Where(qb.Eq("friend_id")).
			ToCql()

		err = r.sess.
			ContextQuery(ctx, stmt2, names2).
			BindMap((qb.M{
				"user_id":   fId,
				"friend_id": uId,
			})).
			GetRelease(&res)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func (r *friendRelationRepository) GetFriends(ctx context.Context, uId string, page []byte, limit uint64) (res []datastruct.FriendRelation, nextPage []byte, err error) {
	stmt, names := qb.
		Select(FRIEND_RELATIONS).
		Where(qb.Eq("user_id")).
		Where(qb.Eq("accepted")).
		ToCql()

	q := r.sess.
		ContextQuery(ctx, stmt, names).
		BindMap((qb.M{
			"user_id":  uId,
			"accepted": true,
		}))
	defer q.Release()

	q.PageState(page)
	if limit == 0 {
		q.PageSize(20)
	} else {
		q.PageSize(int(limit))
	}

	iter := q.Iter()
	err = iter.Select(&res)
	if err != nil {
		return []datastruct.FriendRelation{}, nil, errors.New("no friends found")
	}

	return res, iter.PageState(), nil
}

func (r *friendRelationRepository) GetIncomingFriendRequests(ctx context.Context, uId string, page []byte, limit uint64) (res []datastruct.FriendRelation, nextPage []byte, err error) {
	stmt, names := qb.
		Select(FRIEND_RELATIONS).
		Where(qb.Eq("user_id")).
		Where(qb.Eq("accepted")).
		ToCql()

	q := r.sess.
		ContextQuery(ctx, stmt, names).
		BindMap((qb.M{
			"user_id":  uId,
			"accepted": false,
		}))
	defer q.Release()

	q.PageState(page)
	if limit == 0 {
		q.PageSize(20)
	} else {
		q.PageSize(int(limit))
	}

	iter := q.Iter()
	err = iter.Select(&res)
	if err != nil {
		return []datastruct.FriendRelation{}, nil, errors.New("no friend requests found")
	}

	return res, iter.PageState(), nil
}

func (r *friendRelationRepository) GetFriendCount(ctx context.Context, uId string) (res datastruct.FriendCount, err error) {
	stmt, names := qb.
		Select(FRIEND_COUNT).
		Columns(friendCountMetadata.Columns...).
		Where(qb.Eq("user_id")).
		ToCql()

	err = r.sess.
		ContextQuery(ctx, stmt, names).
		BindMap((qb.M{"user_id": uId})).
		GetRelease(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *friendRelationRepository) GetManyFriendCount(ctx context.Context, ids []string) (res []datastruct.FriendCount, err error) {
	stmt, names := qb.
		Select(FRIEND_COUNT).
		Columns(friendCountMetadata.Columns...).
		Where(qb.In("user_id")).
		ToCql()

	err = r.sess.
		ContextQuery(ctx, stmt, names).
		BindMap((qb.M{"user_id": ids})).
		GetRelease(&res)

	if err != nil {
		return res, err
	}

	return res, nil
}
