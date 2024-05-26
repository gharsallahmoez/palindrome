package database

import (
	"context"
	"fmt"
	"github.com/gharsallahmoez/palindrome/config"
	"github.com/gharsallahmoez/palindrome/infra/database/in-memory"
	"github.com/gharsallahmoez/palindrome/model"
)

// Database represents an interface for interacting with the messages data storage.
type Database interface {
	// SaveMessage saves a message to the database.
	SaveMessage(message model.Message, ctx context.Context) (model.Message, error)
	// GetMessage retrieves a message from the database.
	GetMessage(id string, ctx context.Context) (model.Message, error)
	// UpdateMessage updates a message in the database.
	UpdateMessage(id string, content string, isPalindrome bool, ctx context.Context) (model.Message, error)
	// DeleteMessage deletes a message from the database.
	DeleteMessage(id string, ctx context.Context) error
	// ListMessages retrieves all messages from the database.
	ListMessages(ctx context.Context) ([]model.Message, error)
}

// Create creates a new instance of a database based on the provided configuration.
func Create(conf config.Database) (Database, error) {
	switch conf.Type {
	case "in-memory":
		return in_memory.NewRepo(), nil
	default:
		return nil, fmt.Errorf(fmt.Sprintf("%s is an unknown database type", conf.Type))
	}
}
