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
	"github.com/stretchr/testify/assert"
)

// TestCreateMessageHandler tests CreateMessageHandler function.
func TestCreateMessageHandler(t *testing.T) {
	// Mock database save function for valid request
	dbMock := &DatabaseMock{
		SaveMessageFunc: func(message model.Message, ctx context.Context) (model.Message, error) {
			message.ID = "1"
			return message, nil
		},
	}

	service := svc.NewMessageService(dbMock)

	// Define test cases
	testCases := []struct {
		Name            string
		RequestBody     []byte
		message         string
		ExpectedCode    int
		ExpectedMessage string
		isPalindrome    bool
	}{
		{
			Name: "Valid request, non-Palindrome content",
			RequestBody: []byte(`{
				"content": "test message"
			}`),
			message:         "test message",
			ExpectedCode:    http.StatusCreated,
			ExpectedMessage: "",
		},
		{
			Name: "Valid request, Palindrome content",
			RequestBody: []byte(`{
				"content": "A man a plan a canal Panama"
			}`),
			message:         "A man a plan a canal Panama",
			ExpectedCode:    http.StatusCreated,
			ExpectedMessage: "",
			isPalindrome:    true,
		},
		{
			Name: "Invalid JSON",
			RequestBody: []byte(`{
				"content": "test message",`),
			ExpectedCode:    http.StatusBadRequest,
			ExpectedMessage: "unexpected EOF\n",
		},
		{
			Name: "Empty content",
			RequestBody: []byte(`{
				"content": ""
			}`),
			ExpectedCode:    http.StatusBadRequest,
			ExpectedMessage: "Content cannot be empty\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Create a request with the request body
			req, err := http.NewRequest("POST", "/messages", bytes.NewBuffer(tc.RequestBody))
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder to record the response
			rr := httptest.NewRecorder()

			// Call the CreateMessageHandler method
			handler := http.HandlerFunc(service.CreateMessageHandler)
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tc.ExpectedCode, rr.Code, "Status code should match")

			// Check the response body if an expected message is provided
			if tc.ExpectedMessage != "" {
				assert.Equal(t, tc.ExpectedMessage, rr.Body.String(), "Response body should match")
				return
			}

			// Validate the response of valid requests
			var response svc.MessageResponse
			err = json.NewDecoder(rr.Body).Decode(&response)
			if err != nil {
				t.Fatal(err)
			}

			// Assert that the response body contains a non-empty ID
			assert.NotEmpty(t, response.ID, "ID should not be empty")
			assert.Equal(t, tc.message, response.Content)
			if tc.isPalindrome {
				assert.True(t, response.IsPalindrome)
			}
		})
	}

	t.Run("with failed db operation", func(t *testing.T) {
		t.Parallel()
		dbMock := &DatabaseMock{
			SaveMessageFunc: func(message model.Message, ctx context.Context) (model.Message, error) {
				return model.Message{}, errors.New("some error")
			},
		}

		service := svc.NewMessageService(dbMock)
		RequestBody := []byte(`{
			"content": "test message"
		}`)
		req, err := http.NewRequest("POST", "/messages", bytes.NewBuffer(RequestBody))
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder to record the response
		rr := httptest.NewRecorder()

		// Call the CreateMessageHandler method
		handler := http.HandlerFunc(service.CreateMessageHandler)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code should match")
	})
}
