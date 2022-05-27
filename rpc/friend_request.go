package rpc

import (
	"context"
	"time"

	"github.com/clubo-app/packages/utils"
	cg "github.com/clubo-app/protobuf/common"
	"github.com/clubo-app/protobuf/events"
	rg "github.com/clubo-app/protobuf/relation"
	"github.com/clubo-app/relation-service/datastruct"
)

func (s relationServer) FriendRequest(ctx context.Context, req *rg.FriendRequestRequest) (*cg.SuccessIndicator, error) {
	fr := datastruct.FriendRelation{
		UserId:      req.FriendId,
		FriendId:    req.UserId,
		Accepted:    false,
		RequestedAt: time.Now(),
	}

	err := s.fs.CreateFriendRequest(ctx, fr)
	if err != nil {
		return nil, utils.HandleError(err)
	}

	s.stream.PublishEvent(&events.FriendRequested{
		UserId:   req.UserId,
		FriendId: req.FriendId,
	})

	return &cg.SuccessIndicator{Sucess: true}, nil
}
