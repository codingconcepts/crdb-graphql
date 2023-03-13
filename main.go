package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/codingconcepts/crdb-graphql/resolver"
	"github.com/friendsofgo/graphiql"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	file, err := os.ReadFile("schema.graphql")
	if err != nil {
		log.Fatalf("error reading schema file: %v", err)
	}

	db, err := pgxpool.New(context.Background(), "postgres://root@localhost:26257/defaultdb?sslmode=disable")
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

	resolver := &resolver.Resolver{DB: db}

	schema := graphql.MustParseSchema(string(file), resolver)
	http.Handle("/query", &relay.Handler{Schema: schema})

	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/query")
	if err != nil {
		panic(err)
	}
	http.Handle("/", graphiqlHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
