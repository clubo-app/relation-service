package rpc

import (
	"context"

	"github.com/clubo-app/packages/utils"
	rg "github.com/clubo-app/protobuf/relation"
)

func (s relationServer) GetManyFavoritePartyCount(ctx context.Context, req *rg.GetManyFavoritePartyCountRequest) (*rg.GetManyFavoritePartyCountResponse, error) {
	fp, err := s.fp.GetManyfavoritePartyCount(ctx, req.PartyIds)
	if err != nil {
		return nil, utils.HandleError(err)
	}

	fpMap := make(map[string]uint32, len(fp))
	for _, fc := range fp {
		fpMap[fc.PartyId] = uint32(fc.FavoritePartyCount)
	}

	return &rg.GetManyFavoritePartyCountResponse{FavoriteCounts: fpMap}, nil
}
