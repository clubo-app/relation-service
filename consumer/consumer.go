package consumer

import (
	"context"
	"log"
	"sync"

	"github.com/clubo-app/packages/stream"
	"github.com/clubo-app/protobuf/events"
	"github.com/clubo-app/relation-service/service"
)

type consumer struct {
	stream stream.Stream
	fs     service.FriendRelationService
	ps     service.FavoriteParty
}

func New(stream stream.Stream, fs service.FriendRelationService) consumer {
	return consumer{stream: stream, fs: fs}
}

func (c consumer) Start() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go c.stream.SubscribeToEvent("relation.friend.created.count", events.FriendCreated{}, c.FriendCreated)
	go c.stream.SubscribeToEvent("relation.friend.removed.count", events.FriendRemoved{}, c.FriendRemoved)
	go c.stream.SubscribeToEvent("relation.party.favorited.count", events.PartyFavorited{}, c.PartyFavorited)
	go c.stream.SubscribeToEvent("relation.party.unfavorited.count", events.PartyUnfavorited{}, c.PartyUnfavorited)

	wg.Wait()
}

func (c consumer) FriendCreated(e *events.FriendCreated) {
	err := c.fs.IncreaseFriendCount(context.Background(), e.UserId)

	if err != nil {
		log.Println("Error increasing Count: ", err)
	}
}

func (c consumer) FriendRemoved(e *events.FriendCreated) {
	err := c.fs.DecreaseFriendCount(context.Background(), e.UserId)

	if err != nil {
		log.Println("Error decreasing Count: ", err)
	}
}

func (c consumer) PartyFavorited(e *events.PartyFavorited) {
	err := c.ps.IncreaseFavoritePartyCount(context.Background(), e.PartyId)

	if err != nil {
		log.Println("Error decreasing Count: ", err)
	}
}

func (c consumer) PartyUnfavorited(e *events.PartyUnfavorited) {
	err := c.ps.DecreaseFavoritePartyCount(context.Background(), e.PartyId)

	if err != nil {
		log.Println("Error decreasing Count: ", err)
	}
}
