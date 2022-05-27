package rpc

import (
	"context"
	"encoding/base64"

	"github.com/clubo-app/packages/utils"
	rg "github.com/clubo-app/protobuf/relation"
	"github.com/segmentio/ksuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s relationServer) GetFavoritePartiesByUser(ctx context.Context, req *rg.GetFavoritePartiesByUserRequest) (*rg.PagedFavoriteParties, error) {
	_, err := ksuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid User id")
	}
	p, err := base64.URLEncoding.DecodeString(req.NextPage)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Next Page Param")
	}

	fps, p, err := s.fp.GetFavoritePartiesByUser(ctx, req.UserId, p, req.Limit)
	if err != nil {
		return nil, utils.HandleError(err)
	}

	nextPage := base64.URLEncoding.EncodeToString(p)

	var res []*rg.FavoriteParty
	for _, fp := range fps {
		res = append(res, fp.ToGRPCFavoriteParty())
	}

	return &rg.PagedFavoriteParties{FavoriteParties: res, NextPage: nextPage}, nil
}
