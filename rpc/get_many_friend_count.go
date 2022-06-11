package rpc

import (
	"context"

	"github.com/clubo-app/packages/utils"
	rg "github.com/clubo-app/protobuf/relation"
)

func (s relationServer) GetManyFriendCount(ctx context.Context, req *rg.GetManyFriendCountRequest) (*rg.GetManyFriendCountResponse, error) {
	fs, err := s.fs.GetManyFriendCount(ctx, req.UserIds)
	if err != nil {
		return nil, utils.HandleError(err)
	}

	fcMap := make(map[string]uint32, len(fs))
	for _, fc := range fs {
		fcMap[fc.UserId] = uint32(fc.FriendCount)
	}

	return &rg.GetManyFriendCountResponse{FriendCounts: fcMap}, nil
}
