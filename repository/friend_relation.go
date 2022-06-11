package repository

import (
	"context"
	"errors"
	"log"
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
	PartKey: []string{"user_id", "accepted", "friend_id"},
}

type FriendRelationRepository interface {
	CreateFriendRequest(ctx context.Context, fr datastruct.FriendRelation) error
	DeclineFriendRequest(ctx context.Context, uId, fId string) error
	AcceptFriendRequest(ctx context.Context, uId, fId string) error
	RemoveFriendRelation(ctx context.Context, uId, fId string) error
	GetFriendRelation(ctx context.Context, uId, fId string) (datastruct.FriendRelation, error)
	GetFriendsOfUser(ctx context.Context, uId string, page []byte, limit uint64) ([]datastruct.FriendRelation, []byte, error)
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
		Where(qb.In("user_id")).
		Where(qb.Eq("accepted")).
		Where(qb.In("friend_id")).
		ToCql()

	err := r.sess.ContextQuery(ctx, stmt, names).
		BindMap((qb.M{
			"user_id":   []string{uId, fId},
			"friend_id": []string{uId, fId},
			"accepted":  false,
		})).
		ExecRelease()
	if err != nil {
		return err
	}

	return nil

}

func (r *friendRelationRepository) AcceptFriendRequest(ctx context.Context, uId, fId string) error {
	updateStmt, updateNames := qb.
		Update(FRIEND_RELATIONS).
		Where(qb.Eq("user_id")).
		Where(qb.EqNamed("accepted", "old.accepted")).
		Where(qb.Eq("friend_id")).
		If(qb.EqNamed("accepted", "old.accepted")).
		Set("accepted").
		Set("accepted_at").
		ToCql()

	countStmt, countNames := qb.
		Update(FRIEND_COUNT).
		Where(qb.In("user_id")).
		Add("friend_count").
		ToCql()

	createStmt, createNames := qb.
		Insert(FRIEND_RELATIONS).
		Columns(friendRelationMetadata.Columns...).
		ToCql()

	batch, names := qb.
		Batch().
		AddStmtWithPrefix("u", updateStmt, updateNames).
		AddStmtWithPrefix("c", createStmt, createNames).
		AddStmtWithPrefix("a", countStmt, countNames).
		ToCql()

	err := r.sess.ContextQuery(ctx, batch, names).
		BindMap((qb.M{
			"u.user_id":      uId,
			"u.friend_id":    fId,
			"u.accepted":     true,
			"u.old.accepted": false,
			"u.accepted_at":  time.Now(),
			"c.user_id":      fId,
			"c.friend_id":    uId,
			"c.accepted":     true,
			"c.accepted_at":  time.Now(),
			"c.requested_at": time.Now(),
			"a.friend_count": 1,
			"a.user_id":      []string{uId, fId},
		})).
		ExecRelease()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *friendRelationRepository) RemoveFriendRelation(ctx context.Context, uId, fId string) error {
	stmt, names := qb.
		Delete(FRIEND_RELATIONS).
		Where(qb.In("user_id")).
		Where(qb.Eq("accepted")).
		Where(qb.In("friend_id")).
		ToCql()

	err := r.sess.ContextQuery(ctx, stmt, names).
		BindMap((qb.M{
			"user_id":   []string{uId, fId},
			"friend_id": []string{uId, fId},
			"accepted":  true,
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
		Where(qb.In("user_id")).
		Where(qb.In("accepted")).
		Where(qb.In("friend_id")).
		ToCql()

	err = r.sess.
		ContextQuery(ctx, stmt, names).
		BindMap((qb.M{"user_id": []string{fId, uId}, "friend_id": []string{fId, uId}, "accepted": []bool{true, false}})).
		GetRelease(&res)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (r *friendRelationRepository) GetFriendsOfUser(ctx context.Context, uId string, page []byte, limit uint64) (result []datastruct.FriendRelation, nextPage []byte, err error) {
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
	err = iter.Select(&result)
	if err != nil {
		return []datastruct.FriendRelation{}, nil, errors.New("no friends found")
	}

	return result, iter.PageState(), nil
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
