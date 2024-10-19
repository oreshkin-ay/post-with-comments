package main

import (
	"log"
	"net/http"
	"os"

	"github.com/oreshkin/posts/graph"
	"github.com/oreshkin/posts/internal/auth"
	database "github.com/oreshkin/posts/internal/pkg/db/postgres"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(auth.Middleware())

	database.InitDB()
	defer database.CloseDB()

	database.Migrate()
	// for i := 1; i <= 30; i++ {
	// 	title := fmt.Sprintf("Post Title %d", i)
	// 	content := fmt.Sprintf("This is the content of post number %d. Automatically generated for testing.", i)
	// 	commentsDisabled := rand.Intn(2) == 0

	// 	post := posts.Post{
	// 		Title:            title,
	// 		Content:          content,
	// 		CommentsDisabled: commentsDisabled,
	// 	}

	// 	post.Save()
	// }

	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
