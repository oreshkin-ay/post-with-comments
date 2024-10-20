package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/oreshkin/posts/graph"
	"github.com/oreshkin/posts/internal/auth"
	"github.com/oreshkin/posts/internal/comments"
	database "github.com/oreshkin/posts/internal/pkg/db/postgres"
	"github.com/oreshkin/posts/internal/posts"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

const defaultPort = "8080"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(auth.Middleware())

	database.InitDB()
	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Printf("Error closing the database: %v", err)
		}
	}()

	database.Migrate()

	postRepository := &posts.DBPostRepository{}
	commentRepository := &comments.DBCommentRepository{}

	server := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					PostRepository:    postRepository,
					CommentRepository: commentRepository,
				},
			},
		),
	)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
