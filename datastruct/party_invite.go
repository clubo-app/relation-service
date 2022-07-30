package datastruct

import "time"

type PartyInvite struct {
	UserId     string    `db:"user_id"     validate:"required"`
	InviterId  string    `db:"inviter_id"  validate:"required"`
	PartyId    string    `db:"party_id"    validate:"required"`
	ValidUntil time.Time `validate:"required"`
}
