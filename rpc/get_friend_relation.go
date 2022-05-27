package rpc

import (
	"context"

	"github.com/clubo-app/packages/utils"
	rg "github.com/clubo-app/protobuf/relation"
)

func (s relationServer) GetFriendRelation(ctx context.Context, req *rg.GetFriendRelationRequest) (*rg.FriendRelation, error) {
	fr, err := s.fs.GetFriendRelation(ctx, req.UserId, req.FriendId)
	if err != nil {
		return nil, utils.HandleError(err)
	}

	return fr.ToGRPCFriendRelation(), nil
}
