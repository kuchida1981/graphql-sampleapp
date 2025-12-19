package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jxpress/graphql-sampleapp/graph"
	firestoreClient "github.com/jxpress/graphql-sampleapp/internal/firestore"
	"github.com/jxpress/graphql-sampleapp/internal/repository"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	ctx := context.Background()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		projectID = "demo-project"
	}

	firestoreConn, err := firestoreClient.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore client: %v", err)
	}
	defer firestoreConn.Close()

	messageRepo := repository.NewFirestoreMessageRepository(firestoreConn)

	resolver := graph.NewResolver(messageRepo)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
