package http

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

// ListMessageHandler handles HTTP requests to list messages.
func (s *MessageService) ListMessageHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := s.database.ListMessages(r.Context())
	if err != nil {
		logrus.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httpMessages := make([]MessageResponse, len(messages))

	for index := range messages {
		httpMessages[index] = mapDomainMessageToSchema(messages[index])
	}

	// Marshal message schema into JSON
	messagesJSON, err := json.Marshal(httpMessages)
	if err != nil {
		errMsg := fmt.Sprintf("Error marshalling schema: %v", err)
		logrus.Errorf(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(messagesJSON)
}
