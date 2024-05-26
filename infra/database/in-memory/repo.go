package in_memory

import (
	"context"
	"github.com/gharsallahmoez/palindrome/model"
	"sync"
	"time"
)

// Repo represents an in-memory repository for messages.
type Repo struct {
	messages map[string]model.Message
	mx       sync.Mutex
}

// NewRepo creates a new instance of Repo with an empty map of messages.
func NewRepo() *Repo {
	return &Repo{
		messages: map[string]model.Message{},
	}
}

// SaveMessage saves a message to the database.
func (r *Repo) SaveMessage(message model.Message, _ context.Context) (model.Message, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.messages[message.ID] = message
	return message, nil
}

// GetMessage retrieves a message from the database.
func (r *Repo) GetMessage(id string, _ context.Context) (model.Message, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	msg, exists := r.messages[id]
	if !exists {
		return model.Message{}, model.ErrMessageNotFound
	}
	return msg, nil
}

// UpdateMessage updates a message in the database.
func (r *Repo) UpdateMessage(id string, content string, isPalindrome bool, _ context.Context) (model.Message, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	msg, exists := r.messages[id]
	if !exists {
		return model.Message{}, model.ErrMessageNotFound
	}
	message := model.Message{
		ID:           id,
		Content:      content,
		IsPalindrome: isPalindrome,
		CreatedAt:    msg.CreatedAt,
		UpdatedAt:    time.Now(),
	}
	r.messages[id] = message
	return message, nil
}

// DeleteMessage deletes a message from the database.
func (r *Repo) DeleteMessage(id string, _ context.Context) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	_, exists := r.messages[id]
	if !exists {
		return model.ErrMessageNotFound
	}
	delete(r.messages, id)
	return nil
}

// ListMessages retrieves all messages from the database.
func (r *Repo) ListMessages(_ context.Context) ([]model.Message, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	var messages []model.Message
	for _, m := range r.messages {
		messages = append(messages, m)
	}
	return messages, nil
}
