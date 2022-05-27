package repository

import (
	"strings"

	"github.com/clubo-app/packages/cqlx"
	"github.com/scylladb/gocqlx/v2"
)

type Dao interface {
	NewFriendRelationRepository() FriendRelationRepository
	NewFavoritePartyRepository() FavoritePartyRepository
}

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

func NewDAO(sess *gocqlx.Session) Dao {
	return &dao{sess: sess}
}

func (d *dao) NewFriendRelationRepository() FriendRelationRepository {
	return &friendRelationRepository{sess: d.sess}
}

func (d *dao) NewFavoritePartyRepository() FavoritePartyRepository {
	return &favoritePartyRepository{sess: d.sess}
}
