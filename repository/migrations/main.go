package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/clubo-app/packages/cqlx"
	"github.com/clubo-app/relation-service/repository/migrations/cql"
	"github.com/scylladb/gocqlx/v2/migrate"
)

func main() {
	ctx := context.Background()

	keyspace, exists := os.LookupEnv("SCYLLA_KEYSPACE")
	if !exists {
		log.Fatalln("scylla keyspace not defined")
	}
	hosts, exists := os.LookupEnv("SCYLLA_HOSTS")
	if !exists {
		log.Fatalln("scylla hosts not defined")
	}
	h := strings.Split(hosts, ",")

	manager := cqlx.NewManager(keyspace, h)

	if err := manager.CreateKeyspace(keyspace); err != nil {
		log.Fatalln(err)
	}

	session, err := manager.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer session.Close()

	if err := migrate.FromFS(ctx, session, cql.Files); err != nil {
		log.Fatal("Migrate: ", err)
	}
}
