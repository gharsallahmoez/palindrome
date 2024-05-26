package http

import (
	"encoding/json"
	"github.com/gharsallahmoez/palindrome/model"
	"strings"

	"github.com/sirupsen/logrus"
	"net/http"
)

type MessageRequest struct {
	Content string `json:"content"`
}

type MessageResponse struct {
	ID           string `json:"id"`
	Content      string `json:"content"`
	IsPalindrome bool   `json:"is_palindrome"`
}

// CreateMessageHandler handles HTTP requests to create a new message.
func (s *MessageService) CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	httpRequest := MessageRequest{}
	if err := json.NewDecoder(r.Body).Decode(&httpRequest); err != nil {
		logrus.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the content field
	if httpRequest.Content == "" {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	message := model.NewMessage(httpRequest.Content, isPalindrome(httpRequest.Content))

	// Save the message to the database.
	savedMessage, err := s.database.SaveMessage(message, r.Context())
	if err != nil {
		logrus.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Infof("message with id %s created successfully", savedMessage.ID)

	// Build response JSON.
	response := MessageResponse{
		ID:           savedMessage.ID,
		Content:      savedMessage.Content,
		IsPalindrome: savedMessage.IsPalindrome,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		logrus.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(responseJSON)
}

func isPalindrome(s string) bool {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}
	return true
}
