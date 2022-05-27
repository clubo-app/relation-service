package datastruct

import "time"

type FriendRelation struct {
	UserId     string
	FriendId   string
	Accepted   bool
	CreatedAt  time.Time
	AcceptedAt time.Time
}
