package service

import (
	"context"

	"github.com/clubo-app/relation-service/datastruct"
)

type FriendRelationService interface {
	CreateFriendRequest(ctx context.Context, uId, fId string) error
	DeclineFriendRequest(ctx context.Context, uId, fId string) error
	AcceptFriendRequest(ctx context.Context, uId, fId string) error
	RemoveFriendRelation(ctx context.Context, uId, fId string) error
	GetFriendRelation(ctx context.Context, uId, fId string) (datastruct.FriendRelation, error)
	GetFriends(ctx context.Context, uId string, page []byte, limit uint64) ([]datastruct.FriendRelation, []byte, error)
	GetIncomingFriendRequests(ctx context.Context, uId string, page []byte, limit uint64) ([]datastruct.FriendRelation, []byte, error)
	GetFriendCount(ctx context.Context, uId string) (datastruct.FriendCount, error)
	GetManyFriendCount(ctx context.Context, ids []string) ([]datastruct.FriendCount, error)
}
