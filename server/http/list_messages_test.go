package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gharsallahmoez/palindrome/model"
	svc "github.com/gharsallahmoez/palindrome/server/http"
	"github.com/stretchr/testify/assert"
)

// TestListMessageHandler tests ListMessageHandler function.
func TestListMessageHandler(t *testing.T) {
	// Mock database list function for valid request
	dbMock := &DatabaseMock{
		ListMessagesFunc: func(ctx context.Context) ([]model.Message, error) {
			return []model.Message{
				{
					ID:           "1",
					Content:      "test message",
					IsPalindrome: false,
				},
			}, nil
		},
	}

	service := svc.NewMessageService(dbMock)

	var expectedBody = []svc.MessageResponse{{
		ID:           "1",
		Content:      "test message",
		IsPalindrome: false},
	}

	t.Run("valid request", func(t *testing.T) {
		t.Parallel()
		// Create a request
		req, err := http.NewRequest("GET", "/messages", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder to record the response
		rr := httptest.NewRecorder()

		// Call the ListMessageHandler method
		handler := http.HandlerFunc(service.ListMessageHandler)
		handler.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code, "Status code should match")

		// Check the response body for valid request
		var response []svc.MessageResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expectedBody, response, "Response body should match")

	})

	t.Run("with failed db operation", func(t *testing.T) {
		t.Parallel()
		dbMock := &DatabaseMock{
			ListMessagesFunc: func(ctx context.Context) ([]model.Message, error) {
				return nil, errors.New("some error")
			},
		}

		service := svc.NewMessageService(dbMock)
		req, err := http.NewRequest("GET", "/messages", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder to record the response
		rr := httptest.NewRecorder()

		// Call the ListMessageHandler method
		handler := http.HandlerFunc(service.ListMessageHandler)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code should match")
	})
}
