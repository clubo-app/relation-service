package datastruct

import (
	"time"

	rg "github.com/clubo-app/protobuf/relation"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FriendRelation struct {
	UserId      string    `db:"user_id"      validate:"required"`
	FriendId    string    `db:"friend_id"    validate:"required"`
	Accepted    bool      `db:"accepted"`
	RequestedAt time.Time `db:"requested_at" validate:"required"`
	AcceptedAt  time.Time `db:"accepted_at"`
}

func (fr FriendRelation) ToGRPCFriendRelation() *rg.FriendRelation {
	return &rg.FriendRelation{
		UserId:      fr.UserId,
		FriendId:    fr.FriendId,
		Accepted:    fr.Accepted,
		RequestedAt: timestamppb.New(fr.RequestedAt),
		AcceptedAt:  timestamppb.New(fr.AcceptedAt),
	}
}
