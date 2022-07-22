package repository

import (
	"strings"

	"github.com/clubo-app/packages/cqlx"
	"github.com/go-playground/validator/v10"
	"github.com/scylladb/gocqlx/v2"
)

type dao struct {
	sess *gocqlx.Session
}

func NewDB(keyspace, hosts string) (*gocqlx.Session, error) {
	h := strings.Split(hosts, ",")

	manager := cqlx.NewManager(keyspace, h)

	if err := manager.CreateKeyspace(keyspace); err != nil {
		return nil, err
	}

	session, err := manager.Connect()
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func NewDAO(sess *gocqlx.Session) dao {
	return dao{sess: sess}
}

func (d *dao) NewFriendRelationRepository(val *validator.Validate) FriendRelationRepository {
	return &friendRelationRepository{sess: d.sess, val: val}
}

func (d *dao) NewFavoritePartyRepository(val *validator.Validate) FavoritePartyRepository {
	return &favoritePartyRepository{sess: d.sess, val: val}
}
