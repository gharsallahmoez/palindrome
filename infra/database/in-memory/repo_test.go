package in_memory

import (
	"context"
	"github.com/gharsallahmoez/palindrome/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSaveMessage(t *testing.T) {
	// Create a Repo instance
	repo := NewRepo()

	// Create a message to save
	message := model.NewMessage("test message", false)

	// Save the message
	savedMessage, err := repo.SaveMessage(message, context.Background())

	// Check for errors
	assert.NoError(t, err)
	assert.Equal(t, message, savedMessage)
}

func TestGetMessage(t *testing.T) {
	t.Run("RetrieveExistingMessage", func(t *testing.T) {
		// Create a Repo instance
		repo := NewRepo()

		// Create and save a message
		message := model.NewMessage("test message", false)
		_, err := repo.SaveMessage(message, context.Background())
		if err != nil {
			assert.NoError(t, err)
		}

		// Retrieve the message
		retrievedMessage, err := repo.GetMessage(message.ID, context.Background())

		// Check for errors
		assert.NoError(t, err)
		assert.Equal(t, message, retrievedMessage)
	})

	t.Run("RetrieveNonExistentMessage", func(t *testing.T) {
		// Create a Repo instance
		repo := NewRepo()

		// Try to retrieve a non-existent message
		_, err := repo.GetMessage("non-existent-id", context.Background())
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "message not found")
	})
}

func TestUpdateMessage(t *testing.T) {
	t.Run("UpdateExistingMessage", func(t *testing.T) {
		// Create a Repo instance
		repo := NewRepo()

		// Create and save a message
		message := model.NewMessage("test message", false)
		_, err := repo.SaveMessage(message, context.Background())
		if err != nil {
			assert.NoError(t, err)
		}

		// Update the message
		updatedContent := "updated message"
		updatedMessage, err := repo.UpdateMessage(message.ID, updatedContent, true, context.Background())

		// Check for errors
		assert.NoError(t, err)
		assert.Equal(t, updatedContent, updatedMessage.Content)
		assert.True(t, updatedMessage.IsPalindrome)
		assert.Equal(t, message.ID, updatedMessage.ID)
		assert.NotEqual(t, message.UpdatedAt, updatedMessage.UpdatedAt)
	})

	t.Run("UpdateNonExistentMessage", func(t *testing.T) {
		// Create a Repo instance
		repo := NewRepo()

		// Try to update a non-existent message
		_, err := repo.UpdateMessage("non-existent-id", "new content", false, context.Background())
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "message not found")
	})
}

func TestDeleteMessage(t *testing.T) {
	t.Run("DeleteExistingMessage", func(t *testing.T) {
		// Create a Repo instance
		repo := NewRepo()

		// Create and save a message
		message := model.NewMessage("test message", false)
		_, err := repo.SaveMessage(message, context.Background())
		if err != nil {
			assert.NoError(t, err)
		}

		// Delete the message
		err = repo.DeleteMessage(message.ID, context.Background())

		// Check for errors
		assert.NoError(t, err)

		// Try to retrieve the deleted message
		_, err = repo.GetMessage(message.ID, context.Background())
		assert.Error(t, err)
	})

	t.Run("DeleteNonExistentMessage", func(t *testing.T) {
		// Create a Repo instance
		repo := NewRepo()

		// Try to delete a non-existent message
		err := repo.DeleteMessage("non-existent-id", context.Background())
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "message not found")
	})
}

func TestListMessages(t *testing.T) {
	// Create a Repo instance
	repo := NewRepo()

	// Create and save multiple messages
	message1 := model.NewMessage("message 1", false)
	message2 := model.NewMessage("message 2", true)
	_, err := repo.SaveMessage(message1, context.Background())
	if err != nil {
		assert.NoError(t, err)
	}
	_, err = repo.SaveMessage(message2, context.Background())
	if err != nil {
		assert.NoError(t, err)
	}

	// List all messages
	messages, err := repo.ListMessages(context.Background())
	assert.NoError(t, err)

	// Check if the retrieved list matches the saved messages
	assert.Len(t, messages, 2)
	assert.Contains(t, messages, message1)
	assert.Contains(t, messages, message2)
}
