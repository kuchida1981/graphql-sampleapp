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
	"github.com/kuchida1981/graphql-sampleapp/graph"
	firestoreClient "github.com/kuchida1981/graphql-sampleapp/internal/firestore"
	"github.com/kuchida1981/graphql-sampleapp/internal/postgres"
	firestoreRepo "github.com/kuchida1981/graphql-sampleapp/internal/repository/firestore"
	postgresRepo "github.com/kuchida1981/graphql-sampleapp/internal/repository/postgres"
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

	databaseURL := os.Getenv("DATABASE_URL")
	pgConn, err := postgres.NewClient(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL client: %v", err)
	}
	defer pgConn.Close()

	messageRepo := firestoreRepo.NewFirestoreMessageRepository(firestoreConn)
	userRepo := postgresRepo.NewPostgresUserRepository(pgConn)
	weatherAlertMetadataRepo := postgresRepo.NewPostgresWeatherAlertMetadataRepository(pgConn)
	weatherAlertRepo := firestoreRepo.NewFirestoreWeatherAlertRepository(firestoreConn)

	resolver := graph.NewResolver(messageRepo, userRepo, weatherAlertMetadataRepo, weatherAlertRepo)

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
