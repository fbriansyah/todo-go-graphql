package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fbriansyah/todo-go-graphql/graph"
	"github.com/fbriansyah/todo-go-graphql/graph/generated"
	"github.com/fbriansyah/todo-go-graphql/internal/auth"
	database "github.com/fbriansyah/todo-go-graphql/internal/pkg/db/mysql"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter() // menggunakan chi router

	router.Use(auth.Middleware())

	database.InitDB()  // database initialization
	database.Migrate() // run migration

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router)) // use chi router here
}
