package rpc

import (
	"context"
	"time"

	"github.com/clubo-app/packages/utils"
	rg "github.com/clubo-app/protobuf/relation"
	"github.com/clubo-app/relation-service/datastruct"
	"github.com/segmentio/ksuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s relationServer) FavorParty(ctx context.Context, req *rg.FavorPartyRequest) (*rg.FavoriteParty, error) {
	_, err := ksuid.Parse(req.PartyId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Party id")
	}

	_, err = ksuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid User id")
	}

	fp, err := s.fp.FavorParty(ctx, datastruct.FavoriteParty{
		UserId:      req.UserId,
		PartyId:     req.PartyId,
		FavoritedAt: time.Now(),
	})
	if err != nil {
		return nil, utils.HandleError(err)
	}

	return fp.ToGRPCFavoriteParty(), nil
}
