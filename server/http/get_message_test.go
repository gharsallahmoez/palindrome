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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestGetMessageHandler tests GetMessageHandler function.
func TestGetMessageHandler(t *testing.T) {
	// Mock database get function for valid request
	dbMock := &DatabaseMock{
		GetMessageFunc: func(id string, ctx context.Context) (model.Message, error) {
			return model.Message{
				ID:           "1",
				Content:      "test message",
				IsPalindrome: false,
			}, nil
		},
	}

	service := svc.NewMessageService(dbMock)

	// Define test cases
	testCases := []struct {
		Name         string
		ID           string
		ExpectedCode int
		ExpectedBody svc.MessageResponse
	}{
		{
			Name:         "Valid request",
			ID:           "1",
			ExpectedCode: http.StatusOK,
			ExpectedBody: svc.MessageResponse{
				ID:           "1",
				Content:      "test message",
				IsPalindrome: false,
			},
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
			req, err := http.NewRequest("GET", "/messages/"+tc.ID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Set the request variables
			req = mux.SetURLVars(req, map[string]string{"id": tc.ID})

			// Create a response recorder to record the response
			rr := httptest.NewRecorder()

			// Call the GetMessageHandler method
			handler := http.HandlerFunc(service.GetMessageHandler)
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tc.ExpectedCode, rr.Code, "Status code should match")

			// Check the response body for valid request
			if tc.ExpectedCode == http.StatusOK {
				var response svc.MessageResponse
				err = json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.ExpectedBody, response, "Response body should match")
			}
		})
	}

	t.Run("non-exist message", func(t *testing.T) {
		t.Parallel()
		dbMock := &DatabaseMock{
			GetMessageFunc: func(id string, ctx context.Context) (model.Message, error) {
				return model.Message{}, model.ErrMessageNotFound
			},
		}

		service := svc.NewMessageService(dbMock)
		req, err := http.NewRequest("GET", "/messages/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Set the request variables
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		// Create a response recorder to record the response
		rr := httptest.NewRecorder()

		// Call the GetMessageHandler method
		handler := http.HandlerFunc(service.GetMessageHandler)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code, "Status code should match")
	})

	t.Run("with failed db operation", func(t *testing.T) {
		t.Parallel()
		dbMock := &DatabaseMock{
			GetMessageFunc: func(id string, ctx context.Context) (model.Message, error) {
				return model.Message{}, errors.New("some error")
			},
		}

		service := svc.NewMessageService(dbMock)
		req, err := http.NewRequest("GET", "/messages/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Set the request variables
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		// Create a response recorder to record the response
		rr := httptest.NewRecorder()

		// Call the GetMessageHandler method
		handler := http.HandlerFunc(service.GetMessageHandler)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code should match")
	})
}
