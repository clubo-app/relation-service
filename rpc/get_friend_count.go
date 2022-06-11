package rpc

import (
	"context"

	"github.com/clubo-app/packages/utils"
	rg "github.com/clubo-app/protobuf/relation"
	"github.com/segmentio/ksuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s relationServer) GetFriendCount(ctx context.Context, req *rg.GetFriendCountRequest) (*rg.GetFriendCountResponse, error) {
	_, err := ksuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid User id")
	}

	fc, err := s.fs.GetFriendCount(ctx, req.UserId)
	if err != nil {
		return &rg.GetFriendCountResponse{FriendCount: 0}, utils.HandleError(err)
	}

	return &rg.GetFriendCountResponse{FriendCount: uint32(fc.FriendCount)}, nil
}
