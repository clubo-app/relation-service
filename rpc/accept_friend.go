package rpc

import (
	"context"

	"github.com/clubo-app/packages/utils"
	cg "github.com/clubo-app/protobuf/common"
	"github.com/clubo-app/protobuf/events"
	rg "github.com/clubo-app/protobuf/relation"
)

func (s relationServer) AcceptFriend(ctx context.Context, req *rg.AcceptFriendRequest) (*cg.SuccessIndicator, error) {
	err := s.fs.AcceptFriendRequest(ctx, req.UserId, req.FriendId)
	if err != nil {
		return nil, utils.HandleError(err)
	}

	s.stream.PublishEvent(&events.FriendAccepted{
		UserId:   req.UserId,
		FriendId: req.FriendId,
	})

	return &cg.SuccessIndicator{Sucess: true}, nil
}
