package graph

import "github.com/jxpress/graphql-sampleapp/internal/repository"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	messageRepo repository.MessageRepository
	userRepo    repository.UserRepository
}

func NewResolver(messageRepo repository.MessageRepository, userRepo repository.UserRepository) *Resolver {
	return &Resolver{
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}
