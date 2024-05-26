package http

import (
	"github.com/gharsallahmoez/palindrome/infra/database"
)

// MessageService represents a service for managing messages.
type MessageService struct {
	database database.Database
}

// NewMessageService creates a new instance of MessageService with the provided database.
func NewMessageService(repo database.Database) *MessageService {
	return &MessageService{
		database: repo,
	}
}
