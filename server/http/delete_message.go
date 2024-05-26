package http

import (
	"errors"
	"github.com/gharsallahmoez/palindrome/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

// DeleteMessageHandler handles HTTP requests to delete message.
func (s *MessageService) DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve id from query.
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "id should not be empty", http.StatusBadRequest)
		return
	}

	err := s.database.DeleteMessage(id, r.Context())
	if err != nil {
		if errors.Is(err, model.ErrMessageNotFound) {
			http.Error(w, "message not found", http.StatusNotFound)
		} else {
			logrus.Errorf(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
