package datastruct

import "time"

type PartyParticipant struct {
	UserId   string    `db:"user_id"    validate:"required"`
	PartyId  string    `db:"party_id"   validate:"required"`
	JoinedAt time.Time `db:"joined_at"  validate:"required"`
}
