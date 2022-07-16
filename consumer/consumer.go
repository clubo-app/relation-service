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
}

func New(stream stream.Stream, fs service.FriendRelationService) consumer {
	return consumer{stream: stream, fs: fs}
}

func (c consumer) Start() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go c.stream.SubscribeToEvent("relation.friend.created.friendCount", events.FriendCreated{}, c.FriendCreated)
	go c.stream.SubscribeToEvent("relation.friend.removed.friendCount", events.FriendRemoved{}, c.FriendRemoved)

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
