package main

import (
	"log"

	"github.com/clubo-app/packages/stream"
	"github.com/clubo-app/relation-service/config"
	"github.com/clubo-app/relation-service/consumer"
	"github.com/clubo-app/relation-service/repository"
	"github.com/clubo-app/relation-service/rpc"
	"github.com/nats-io/nats.go"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	opts := []nats.Option{nats.Name("Relation Service")}
	stream, err := stream.Connect(c.NATS_CLUSTER, opts)
	if err != nil {
		log.Fatalln(err)
	}
	defer stream.Close()

	cqlx, err := repository.NewDB(c.CQL_KEYSPACE, c.CQL_HOSTS)
	if err != nil {
		log.Fatal(err)
	}
	defer cqlx.Close()

	dao := repository.NewDAO(cqlx)

	fs := dao.NewFriendRelationRepository()

	con := consumer.New(stream, fs)
	go con.Start()

	r := rpc.NewRelationServer(fs, dao.NewFavoritePartyRepository(), stream)
	rpc.Start(r, c.PORT)
}
