package service

import (
	"context"

	"github.com/clubo-app/relation-service/datastruct"
)

type FriendRelationService interface {
	CreateFriendRequest(ctx context.Context, fr datastruct.FriendRelation) error
	AcceptFriendRequest(ctx context.Context, uId, fId string) error
	RemoveFriendRelation(ctx context.Context, uId, fId string) error
	GetFriendRelation(ctx context.Context, uId, fId string) (datastruct.FriendRelation, error)
	GetFriendsOfUser(ctx context.Context, uId string, page []byte, limit uint32) ([]datastruct.FriendRelation, []byte, error)
}
