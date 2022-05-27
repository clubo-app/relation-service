package datastruct

import (
	"time"

	rg "github.com/clubo-app/protobuf/relation"
)

type FriendRelation struct {
	UserId      string    `json:"user_id"      db:"user_id"      validate:"required"`
	FriendId    string    `json:"friend_id"    db:"friend_id"    validate:"required"`
	Accepted    bool      `json:"accepted"     db:"accepted"`
	RequestedAt time.Time `json:"requested_at" db:"requested_at" validate:"required"`
	AcceptedAt  time.Time `json:"accepted_at"  db:"accepted_at"`
}

func (fr FriendRelation) ToGRPCFriendRelation() *rg.FriendRelation {
	return &rg.FriendRelation{
		UserId:      fr.UserId,
		FriendId:    fr.FriendId,
		Accepted:    fr.Accepted,
		RequestedAt: fr.RequestedAt.UTC().Format(time.RFC3339),
		AcceptedAt:  fr.AcceptedAt.UTC().Format(time.RFC3339),
	}
}
