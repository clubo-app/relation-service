package repository

import (
	"context"
	"errors"

	"github.com/clubo-app/relation-service/datastruct"
	"github.com/go-playground/validator/v10"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

const (
	FAVORITE_PARTIES string = "favorite_parties"
)

var favoritePartyMetadata = table.Metadata{
	Name:    FAVORITE_PARTIES,
	Columns: []string{"user_id", "party_id", "favorited_at"},
	PartKey: []string{"user_id", "party_id", "favorited_at"},
}
var favoritePartyTable = table.New(favoritePartyMetadata)

type FavoritePartyRepository interface {
	FavorParty(ctx context.Context, fp datastruct.FavoriteParty) (datastruct.FavoriteParty, error)
	DefavorParty(ctx context.Context, uId, pId string) error
	GetFavoritePartiesByUser(ctx context.Context, uId string, page []byte, limit uint32) ([]datastruct.FavoriteParty, []byte, error)
	GetFavorisingUsersByParty(ctx context.Context, pId string, page []byte, limit uint32) ([]datastruct.FavoriteParty, []byte, error)
}

type favoritePartyRepository struct {
	sess *gocqlx.Session
}

func (r *favoritePartyRepository) FavorParty(ctx context.Context, fp datastruct.FavoriteParty) (datastruct.FavoriteParty, error) {
	v := validator.New()
	err := v.Struct(fp)
	if err != nil {
		return datastruct.FavoriteParty{}, err
	}

	stmt, names := qb.
		Insert(FAVORITE_PARTIES).
		Columns(favoritePartyMetadata.Columns...).
		ToCql()

	err = r.sess.
		Query(stmt, names).
		BindStruct(fp).
		ExecRelease()
	if err != nil {
		return datastruct.FavoriteParty{}, err
	}

	return fp, nil
}

func (r *favoritePartyRepository) DefavorParty(ctx context.Context, uId, pId string) error {
	stmt, names := qb.
		Delete(FAVORITE_PARTIES).
		Where(qb.Eq("user_id")).
		Where(qb.Eq("party_id")).
		ToCql()

	err := r.sess.
		Query(stmt, names).
		BindMap((qb.M{"party_id": pId, "user_id": uId})).
		ExecRelease()
	if err != nil {
		return err
	}
	return nil
}

func (r *favoritePartyRepository) GetFavoritePartiesByUser(ctx context.Context, uId string, page []byte, limit uint32) (result []datastruct.FavoriteParty, nextPage []byte, err error) {
	stmt, names := qb.
		Select(FAVORITE_PARTIES).
		Where(qb.Eq("user_id")).
		ToCql()

	q := r.sess.
		Query(stmt, names).
		BindMap((qb.M{"user_id": uId}))
	defer q.Release()

	q.PageState(page)
	if limit == 0 {
		q.PageSize(10)
	} else {
		q.PageSize(int(limit))
	}

	iter := q.Iter()
	err = iter.Select(&result)
	if err != nil {
		return []datastruct.FavoriteParty{}, nil, errors.New("no favorite parties found")
	}

	return result, iter.PageState(), nil
}

func (r *favoritePartyRepository) GetFavorisingUsersByParty(ctx context.Context, pId string, page []byte, limit uint32) (result []datastruct.FavoriteParty, nextPage []byte, err error) {
	stmt, names := qb.
		Select(FAVORITE_PARTIES).
		Where(qb.Eq("party_id")).
		ToCql()

	q := r.sess.
		Query(stmt, names).
		BindMap((qb.M{"party_id": pId}))
	defer q.Release()

	q.PageState(page)
	if limit == 0 {
		q.PageSize(10)
	} else {
		q.PageSize(int(limit))
	}

	iter := q.Iter()
	err = iter.Select(&result)
	if err != nil {
		return []datastruct.FavoriteParty{}, nil, errors.New("no favorite parties found")
	}

	return result, iter.PageState(), nil
}
