package http

import (
	"encoding/json"
	"errors"
	"github.com/gharsallahmoez/palindrome/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

// UpdateMessageHandler handles HTTP requests to create a new message.
func (s *MessageService) UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve id from query.
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "id should not be empty", http.StatusBadRequest)
		return
	}

	httpRequest := MessageRequest{}
	if err := json.NewDecoder(r.Body).Decode(&httpRequest); err != nil {
		logrus.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// update the message in the database.
	savedMessage, err := s.database.UpdateMessage(id, httpRequest.Content, isPalindrome(httpRequest.Content), r.Context())
	if err != nil {
		if errors.Is(err, model.ErrMessageNotFound) {
			http.Error(w, "message not found", http.StatusNotFound)
		} else {
			logrus.Errorf(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	logrus.Infof("message with id %s updated successfully", savedMessage.ID)

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
	_, _ = w.Write(responseJSON)
}
