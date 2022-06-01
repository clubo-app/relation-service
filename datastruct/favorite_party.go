package datastruct

import (
	"time"

	rg "github.com/clubo-app/protobuf/relation"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FavoriteParty struct {
	UserId      string    `json:"user_id"      db:"user_id"      validate:"required"`
	PartyId     string    `json:"party_id"     db:"party_id"     validate:"required"`
	FavoritedAt time.Time `json:"favorited_at" db:"favorited_at" validate:"required"`
}

func (f FavoriteParty) ToGRPCFavoriteParty() *rg.FavoriteParty {
	return &rg.FavoriteParty{
		UserId:      f.UserId,
		PartyId:     f.PartyId,
		FavoritedAt: timestamppb.New(f.FavoritedAt),
	}
}
