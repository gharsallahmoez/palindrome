package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gharsallahmoez/palindrome/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

// GetMessageHandler handles HTTP requests to retrieve message.
func (s *MessageService) GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve id from query.
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "id should not be empty", http.StatusBadRequest)
		return
	}

	message, err := s.database.GetMessage(id, r.Context())
	if err != nil {
		if errors.Is(err, model.ErrMessageNotFound) {
			http.Error(w, "message not found", http.StatusNotFound)
		} else {
			logrus.Errorf(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	httpMessage := mapDomainMessageToSchema(message)

	// Marshal message schema into JSON
	executionsJSON, err := json.Marshal(httpMessage)
	if err != nil {
		errMsg := fmt.Sprintf("Error marshalling schema: %v", err)
		logrus.Errorf(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(executionsJSON)
}

// mapDomainMessageToSchema maps a message model to a http schema.
func mapDomainMessageToSchema(message model.Message) MessageResponse {
	return MessageResponse{
		ID:           message.ID,
		Content:      message.Content,
		IsPalindrome: message.IsPalindrome,
	}
}
