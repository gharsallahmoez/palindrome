package model_test

import (
	"github.com/gharsallahmoez/palindrome/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMessage(t *testing.T) {
	// input parameters
	content := "test message"
	isPalindrome := false

	// Call the NewMessage function
	message := model.NewMessage(content, isPalindrome)

	// Assertions
	assert.NotEmpty(t, message.ID)
	assert.Equal(t, content, message.Content, "content should match")
	assert.Equal(t, isPalindrome, message.IsPalindrome, "IsPalindrome should match")
}
