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

	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Printf("Error closing the database: %v", err)
		}
	}()

	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

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

	log.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
