package rpc

import (
	"context"

	"github.com/clubo-app/packages/utils"
	rg "github.com/clubo-app/protobuf/relation"
	"github.com/segmentio/ksuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s relationServer) GetFavoritePartyCount(ctx context.Context, req *rg.GetFavoritePartyCountRequest) (*rg.GetFavoritePartyCountResponse, error) {
	_, err := ksuid.Parse(req.PartyId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Party id")
	}

	fp, err := s.fp.GetfavoritePartyCount(ctx, req.PartyId)
	if err != nil {
		return &rg.GetFavoritePartyCountResponse{FavoriteCount: 0}, utils.HandleError(err)
	}

	return &rg.GetFavoritePartyCountResponse{FavoriteCount: uint32(fp.FavoritePartyCount)}, nil
}
