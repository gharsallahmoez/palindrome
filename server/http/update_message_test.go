package http_test

import (
	"bytes"
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

// TestUpdateMessageHandler tests UpdateMessageHandler function.
func TestUpdateMessageHandler(t *testing.T) {
	// Mock database update function for valid request
	dbMock := &DatabaseMock{
		UpdateMessageFunc: func(id, content string, isPalindrome bool, ctx context.Context) (model.Message, error) {
			return model.Message{
				ID:           id,
				Content:      content,
				IsPalindrome: isPalindrome,
			}, nil
		},
	}

	service := svc.NewMessageService(dbMock)

	// Define test cases
	testCases := []struct {
		Name         string
		ID           string
		RequestBody  []byte
		ExpectedCode int
		ExpectedBody svc.MessageResponse
	}{
		{
			Name: "Valid request",
			ID:   "1",
			RequestBody: []byte(`{
				"content": "updated message"
			}`),
			ExpectedCode: http.StatusOK,
			ExpectedBody: svc.MessageResponse{
				ID:           "1",
				Content:      "updated message",
				IsPalindrome: false,
			},
		},
		{
			Name: "Empty ID",
			ID:   "",
			RequestBody: []byte(`{
				"content": "updated message"
			}`),
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name: "Invalid JSON",
			ID:   "1",
			RequestBody: []byte(`{
				"content": "updated message",`),
			ExpectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Create a request with the request body
			req, err := http.NewRequest("PUT", "/messages/"+tc.ID, bytes.NewBuffer(tc.RequestBody))
			if err != nil {
				t.Fatal(err)
			}

			// Set the request variables
			req = mux.SetURLVars(req, map[string]string{"id": tc.ID})

			// Create a response recorder to record the response
			rr := httptest.NewRecorder()

			// Call the UpdateMessageHandler method
			handler := http.HandlerFunc(service.UpdateMessageHandler)
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
			UpdateMessageFunc: func(id, content string, isPalindrome bool, ctx context.Context) (model.Message, error) {
				return model.Message{}, model.ErrMessageNotFound
			},
		}

		service := svc.NewMessageService(dbMock)
		RequestBody := []byte(`{
			"content": "test message"
		}`)
		req, err := http.NewRequest("PUT", "/messages/1", bytes.NewBuffer(RequestBody))
		if err != nil {
			t.Fatal(err)
		}

		// Set the request variables
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		// Create a response recorder to record the response
		rr := httptest.NewRecorder()

		// Call the UpdateMessageHandler method
		handler := http.HandlerFunc(service.UpdateMessageHandler)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code, "Status code should match")
	})

	t.Run("with failed db operation", func(t *testing.T) {
		t.Parallel()
		dbMock := &DatabaseMock{
			UpdateMessageFunc: func(id, content string, isPalindrome bool, ctx context.Context) (model.Message, error) {
				return model.Message{}, errors.New("some error")
			},
		}

		service := svc.NewMessageService(dbMock)
		RequestBody := []byte(`{
			"content": "test message"
		}`)
		req, err := http.NewRequest("PUT", "/messages/1", bytes.NewBuffer(RequestBody))
		if err != nil {
			t.Fatal(err)
		}

		// Set the request variables
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		// Create a response recorder to record the response
		rr := httptest.NewRecorder()

		// Call the UpdateMessageHandler method
		handler := http.HandlerFunc(service.UpdateMessageHandler)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code should match")
	})
}
