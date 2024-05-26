package http_test

import (
	"context"
	"errors"
	"github.com/gharsallahmoez/palindrome/model"
	"net/http"
	"net/http/httptest"
	"testing"

	svc "github.com/gharsallahmoez/palindrome/server/http"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestDeleteMessageHandler tests DeleteMessageHandler function.
func TestDeleteMessageHandler(t *testing.T) {
	// Mock database delete function for valid request
	dbMock := &DatabaseMock{
		DeleteMessageFunc: func(id string, ctx context.Context) error {
			return nil
		},
	}

	service := svc.NewMessageService(dbMock)

	// Define test cases
	testCases := []struct {
		Name         string
		ID           string
		ExpectedCode int
	}{
		{
			Name:         "Valid request",
			ID:           "1",
			ExpectedCode: http.StatusNoContent,
		},
		{
			Name:         "Empty ID",
			ID:           "",
			ExpectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Create a request with the ID
			req, err := http.NewRequest("DELETE", "/messages/"+tc.ID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Set the request variables
			req = mux.SetURLVars(req, map[string]string{"id": tc.ID})

			// Create a response recorder to record the response
			rr := httptest.NewRecorder()

			// Call the DeleteMessageHandler method
			handler := http.HandlerFunc(service.DeleteMessageHandler)
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tc.ExpectedCode, rr.Code, "Status code should match")
		})
	}

	t.Run("non-exist message", func(t *testing.T) {
		t.Parallel()
		dbMock := &DatabaseMock{
			DeleteMessageFunc: func(id string, ctx context.Context) error {
				return model.ErrMessageNotFound
			},
		}

		service := svc.NewMessageService(dbMock)
		req, err := http.NewRequest("DELETE", "/messages/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Set the request variables
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		// Create a response recorder to record the response
		rr := httptest.NewRecorder()

		// Call the DeleteMessageHandler method
		handler := http.HandlerFunc(service.DeleteMessageHandler)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code, "Status code should match")
	})

	t.Run("with failed db operation", func(t *testing.T) {
		t.Parallel()
		dbMock := &DatabaseMock{
			DeleteMessageFunc: func(id string, ctx context.Context) error {
				return errors.New("some error")
			},
		}

		service := svc.NewMessageService(dbMock)
		req, err := http.NewRequest("DELETE", "/messages/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Set the request variables
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		// Create a response recorder to record the response
		rr := httptest.NewRecorder()

		// Call the DeleteMessageHandler method
		handler := http.HandlerFunc(service.DeleteMessageHandler)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code should match")
	})
}
