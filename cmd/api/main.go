package main

import (
	"github.com/phyrwork/bogglr/pkg/api"
	"github.com/phyrwork/bogglr/pkg/api/generated"
	"github.com/phyrwork/bogglr/pkg/database"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var (
		db  *database.DB
		err error
	)
	if db, err = database.Open(""); err != nil {
		log.Fatalf("database open error: %v", err)
	} else if err = database.Migrate(db); err != nil {
		log.Fatalf("database migrate error: %v", err)
	}

	resolver := api.Resolver{DB: db}
	config := generated.Config{Resolvers: &resolver}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
