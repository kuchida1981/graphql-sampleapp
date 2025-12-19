//go:build ignore

package main

import (
	"context"
	"log"
	"os"
	"time"

	firestoreClient "github.com/kuchida1981/graphql-sampleapp/internal/firestore"
)

type Message struct {
	ID        string    `firestore:"id"`
	Content   string    `firestore:"content"`
	Author    string    `firestore:"author"`
	CreatedAt time.Time `firestore:"createdAt"`
}

func main() {
	ctx := context.Background()

	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		projectID = "demo-project"
	}

	client, err := firestoreClient.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore client: %v", err)
	}
	defer client.Close()

	messages := []Message{
		{
			ID:        "msg1",
			Content:   "Hello, Firestore! This is the first message.",
			Author:    "Alice",
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:        "msg2",
			Content:   "GraphQL and Firestore integration is working!",
			Author:    "Bob",
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        "msg3",
			Content:   "Docker Compose makes local development easy.",
			Author:    "Charlie",
			CreatedAt: time.Now(),
		},
	}

	for _, msg := range messages {
		_, err := client.Collection("messages").Doc(msg.ID).Set(ctx, msg)
		if err != nil {
			log.Fatalf("Failed to create message %s: %v", msg.ID, err)
		}
		log.Printf("Successfully created message: %s", msg.ID)
	}

	log.Printf("Successfully seeded %d messages", len(messages))
}
