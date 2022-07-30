package repository

import (
	"context"
	"sync"
	"time"

	"github.com/clubo-app/relation-service/datastruct"
	"github.com/go-playground/validator/v10"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	PARTY_PARTICIPANTS         string = "party_participants"
	PARTY_PARTICIPANTS_BY_USER string = "party_participants_by_user"
	PARTY_INVITES              string = "party_invites"
)

var partyParticipantMetadata = table.Metadata{
	Name:    PARTY_PARTICIPANTS,
	Columns: []string{"user_id", "party_id", "joined_at"},
	PartKey: []string{"party_id", "user_id"},
}

var partyInviteMetadata = table.Metadata{
	Name:    PARTY_INVITES,
	Columns: []string{"user_id", "party_id", "inviter_id"},
	PartKey: []string{"user_id", "party_id"},
}

type PartyParticipantsRepository interface {
	Invite(context.Context, InviteParams) (datastruct.PartyInvite, error)
	Decline(context.Context, UserPartyParams) error
	Accept(context.Context, UserPartyParams) error
	GetUserInvites(context.Context, GetUserInvitesParams) ([]datastruct.PartyInvite, []byte, error)
	Join(context.Context, UserPartyParams) error
	Leave(context.Context, UserPartyParams) error
	GetPartyParticipants(context.Context, GetPartyParticipantsParams) ([]datastruct.PartyParticipant, []byte, error)
}

type partyParticipantRepository struct {
	sess *gocqlx.Session
	val  *validator.Validate
}

type InviteParams struct {
	UserId    string
	InviterId string
	PartyId   string
	ValidFor  time.Duration
}

func (r partyParticipantRepository) Invite(ctx context.Context, params InviteParams) (datastruct.PartyInvite, error) {
	i := datastruct.PartyInvite{
		UserId:     params.UserId,
		InviterId:  params.InviterId,
		PartyId:    params.PartyId,
		ValidUntil: time.Now().Add(params.ValidFor),
	}
	err := r.val.StructCtx(ctx, i)
	if err != nil {
		return datastruct.PartyInvite{}, err
	}

	stmt, names := qb.
		Insert(PARTY_INVITES).
		Unique().
		Columns(partyInviteMetadata.Columns...).
		TTL(params.ValidFor).
		ToCql()

	err = r.sess.
		ContextQuery(ctx, stmt, names).
		BindStruct(i).
		ExecRelease()
	if err != nil {
		return datastruct.PartyInvite{}, err
	}

	return i, nil
}

type UserPartyParams struct {
	UserId  string
	PartyId string
}

func (r partyParticipantRepository) Decline(ctx context.Context, params UserPartyParams) error {
	stmt, names := qb.
		Delete(PARTY_INVITES).
		Where(qb.Eq("user_id")).
		Where(qb.Eq("party_id")).
		ToCql()

	err := r.sess.ContextQuery(ctx, stmt, names).
		BindMap((qb.M{
			"user_id":  params.UserId,
			"party_id": params.PartyId,
		})).
		ExecRelease()
	if err != nil {
		return err
	}

	return nil
}

func (r partyParticipantRepository) Accept(ctx context.Context, params UserPartyParams) error {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	var err error

	go func() {
		defer wg.Done()

		tmp := r.Decline(ctx, params)
		if tmp != nil {
			err = tmp
		}
	}()

	go func() {
		defer wg.Done()

		tmp := r.Join(ctx, params)
		if tmp != nil {
			err = tmp
		}

	}()
	return err
}

type GetUserInvitesParams struct {
	UId   string
	Page  []byte
	Limit int
}

func (r partyParticipantRepository) GetUserInvites(ctx context.Context, params GetUserInvitesParams) (res []datastruct.PartyInvite, nextPage []byte, err error) {
	stmt, names := qb.
		Select(PARTY_INVITES).
		Where(qb.Eq("user_id")).
		ToCql()

	q := r.sess.
		ContextQuery(ctx, stmt, names).
		BindMap((qb.M{
			"user_id": params.UId,
		}))

	q.PageState(params.Page)
	if params.Limit == 0 {
		q.PageSize(20)
	} else {
		q.PageSize(params.Limit)
	}

	iter := q.Iter()
	err = iter.Select(&res)
	if err != nil {
		return []datastruct.PartyInvite{}, nil, status.Error(codes.Internal, "No invites found")
	}

	return res, iter.PageState(), nil
}

func (r partyParticipantRepository) Join(ctx context.Context, params UserPartyParams) error {
	p := datastruct.PartyParticipant{
		UserId:   params.UserId,
		PartyId:  params.PartyId,
		JoinedAt: time.Now(),
	}

	stmt, names := qb.
		Insert(PARTY_PARTICIPANTS).
		Unique().
		Columns(partyParticipantMetadata.Columns...).
		ToCql()

	err := r.sess.
		ContextQuery(ctx, stmt, names).
		BindStruct(p).
		ExecRelease()
	if err != nil {
		return err
	}
	return nil
}

func (r partyParticipantRepository) Leave(ctx context.Context, params UserPartyParams) error {
	stmt, names := qb.
		Delete(PARTY_PARTICIPANTS).
		Where(qb.Eq("user_id")).
		Where(qb.Eq("party_id")).
		ToCql()

	err := r.sess.
		ContextQuery(ctx, stmt, names).
		BindMap((qb.M{
			"user_id":  params.UserId,
			"party_id": params.PartyId,
		})).
		ExecRelease()
	if err != nil {
		return err
	}
	return nil
}

type GetPartyParticipantsParams struct {
	PId   string
	Page  []byte
	Limit int
}

func (r partyParticipantRepository) GetPartyParticipants(ctx context.Context, params GetPartyParticipantsParams) (res []datastruct.PartyParticipant, nextPage []byte, err error) {
	stmt, names := qb.
		Select(PARTY_PARTICIPANTS).
		Where(qb.Eq("party_id")).
		ToCql()

	q := r.sess.
		ContextQuery(ctx, stmt, names).
		BindMap((qb.M{
			"party_id": params.PId,
		}))

	q.PageState(params.Page)
	if params.Limit == 0 {
		q.PageSize(20)
	} else {
		q.PageSize(params.Limit)
	}

	iter := q.Iter()
	err = iter.Select(&res)
	if err != nil {
		return []datastruct.PartyParticipant{}, nil, status.Error(codes.Internal, "No invites found")
	}

	return res, iter.PageState(), nil
}
