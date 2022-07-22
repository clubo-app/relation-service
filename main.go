package main

import (
	"log"

	"github.com/clubo-app/packages/stream"
	"github.com/clubo-app/relation-service/config"
	"github.com/clubo-app/relation-service/consumer"
	"github.com/clubo-app/relation-service/repository"
	"github.com/clubo-app/relation-service/rpc"
	"github.com/go-playground/validator/v10"
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
	val := validator.New()

	fs := dao.NewFriendRelationRepository(val)
	ps := dao.NewFavoritePartyRepository(val)

	con := consumer.New(stream, fs, ps)
	go con.Start()

	r := rpc.NewRelationServer(fs, ps, stream)
	rpc.Start(r, c.PORT)
}
