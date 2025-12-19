package repository

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/jxpress/graphql-sampleapp/internal/domain"
	"google.golang.org/api/iterator"
)

type FirestoreMessageRepository struct {
	client *firestore.Client
}

func NewFirestoreMessageRepository(client *firestore.Client) *FirestoreMessageRepository {
	return &FirestoreMessageRepository{
		client: client,
	}
}

func (r *FirestoreMessageRepository) List(ctx context.Context) ([]*domain.Message, error) {
	log.Println("Fetching all messages from Firestore")

	iter := r.client.Collection("messages").OrderBy("createdAt", firestore.Desc).Documents(ctx)
	defer iter.Stop()

	var messages []*domain.Message
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating messages: %v", err)
			return nil, err
		}

		var msg domain.Message
		if err := doc.DataTo(&msg); err != nil {
			log.Printf("Error converting document to Message: %v", err)
			return nil, err
		}

		messages = append(messages, &msg)
	}

	log.Printf("Successfully fetched %d messages", len(messages))
	return messages, nil
}

func (r *FirestoreMessageRepository) GetByID(ctx context.Context, id string) (*domain.Message, error) {
	log.Printf("Fetching message with ID: %s", id)

	doc, err := r.client.Collection("messages").Doc(id).Get(ctx)
	if err != nil {
		log.Printf("Error fetching message %s: %v", id, err)
		return nil, fmt.Errorf("message not found: %w", err)
	}

	var msg domain.Message
	if err := doc.DataTo(&msg); err != nil {
		log.Printf("Error converting document to Message: %v", err)
		return nil, err
	}

	log.Printf("Successfully fetched message: %s", id)
	return &msg, nil
}
