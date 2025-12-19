package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	conf := &firebase.Config{ProjectID: projectID}

	var opts []option.ClientOption

	app, err := firebase.NewApp(ctx, conf, opts...)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Firestore client initialized for project: %s", projectID)
	return client, nil
}
