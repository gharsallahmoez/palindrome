package model

import (
	"github.com/google/uuid"
	"time"
)

// Message represents the message entity.
type Message struct {
	ID           string
	Content      string
	IsPalindrome bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewMessage creates message model.
func NewMessage(content string, isPalindrome bool) Message {
	now := time.Now()
	message := Message{
		ID:           uuid.NewString(),
		Content:      content,
		IsPalindrome: isPalindrome,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return message
}
