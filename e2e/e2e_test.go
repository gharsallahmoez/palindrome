//go:build e2e

package e2e_test

import (
	"bytes"
	"encoding/json"
	svc "github.com/gharsallahmoez/palindrome/server/http"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

var host = "http://127.0.0.1:8080"

func TestCreateMessage(t *testing.T) {
	url := host + "/messages"

	tests := []struct {
		name               string
		content            string
		expectedPalindrome bool
	}{
		{
			name:               "non-palindrome message",
			content:            "my message",
			expectedPalindrome: false,
		},
		{
			name:               "palindrome message",
			content:            "A man a plan a canal Panama",
			expectedPalindrome: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := svc.MessageRequest{Content: tt.content}

			// Convert request to JSON
			reqJSON, err := json.Marshal(request)
			assert.NoError(t, err, "Error marshalling JSON")

			// Send POST request
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqJSON))
			assert.NoError(t, err, "Error sending POST request")

			// Check response status
			assert.Equal(t, http.StatusCreated, resp.StatusCode, "Unexpected response status")

			// Convert response
			var mdlResp svc.MessageResponse
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err, "Error reading response body")
			err = json.Unmarshal(body, &mdlResp)
			assert.NoError(t, err, "Error unmarshalling response")

			assert.Equal(t, tt.content, mdlResp.Content, "response content should match")
			assert.NotEmpty(t, mdlResp.ID)
			assert.Equal(t, tt.expectedPalindrome, mdlResp.IsPalindrome)
		})
	}
}

func TestRetrieveMessage(t *testing.T) {
	// create a message
	createURL := host + "/messages"
	request := svc.MessageRequest{Content: "test message"}

	reqJSON, err := json.Marshal(request)
	assert.NoError(t, err, "Error marshalling JSON")

	resp, err := http.Post(createURL, "application/json", bytes.NewBuffer(reqJSON))
	assert.NoError(t, err, "Error sending POST request")
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Unexpected response status")

	var createdMessage svc.MessageResponse
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Error reading response body")
	err = json.Unmarshal(body, &createdMessage)
	assert.NoError(t, err, "Error unmarshalling response")

	// retrieve the created message
	retrieveURL := host + "/messages/" + createdMessage.ID
	resp, err = http.Get(retrieveURL)
	assert.NoError(t, err, "Error sending GET request")

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected response status")

	var retrievedMessage svc.MessageResponse
	body, err = io.ReadAll(resp.Body)
	assert.NoError(t, err, "Error reading response body")
	err = json.Unmarshal(body, &retrievedMessage)
	assert.NoError(t, err, "Error unmarshalling response")

	assert.Equal(t, createdMessage.ID, retrievedMessage.ID, "Retrieved ID should match created ID")
	assert.Equal(t, createdMessage.Content, retrievedMessage.Content, "Retrieved content should match created content")
	assert.Equal(t, createdMessage.IsPalindrome, retrievedMessage.IsPalindrome, "Retrieved palindrome status should match")
}

func TestUpdateMessage(t *testing.T) {
	// create a message
	createURL := host + "/messages"
	request := svc.MessageRequest{Content: "initial message"}

	reqJSON, err := json.Marshal(request)
	assert.NoError(t, err, "Error marshalling JSON")

	resp, err := http.Post(createURL, "application/json", bytes.NewBuffer(reqJSON))
	assert.NoError(t, err, "Error sending POST request")
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Unexpected response status")

	var createdMessage svc.MessageResponse
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Error reading response body")
	err = json.Unmarshal(body, &createdMessage)
	assert.NoError(t, err, "Error unmarshalling response")

	// update the created message
	updateURL := host + "/messages/" + createdMessage.ID
	updateRequest := svc.MessageRequest{Content: "updated message"}

	updateReqJSON, err := json.Marshal(updateRequest)
	assert.NoError(t, err, "Error marshalling JSON")

	req, err := http.NewRequest(http.MethodPut, updateURL, bytes.NewBuffer(updateReqJSON))
	assert.NoError(t, err, "Error creating PUT request")

	client := &http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err, "Error sending PUT request")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected response status")

	var updatedMessage svc.MessageResponse
	body, err = io.ReadAll(resp.Body)
	assert.NoError(t, err, "Error reading response body")
	err = json.Unmarshal(body, &updatedMessage)
	assert.NoError(t, err, "Error unmarshalling response")

	assert.Equal(t, createdMessage.ID, updatedMessage.ID, "Updated ID should match created ID")
	assert.Equal(t, "updated message", updatedMessage.Content, "Updated content should match")
	assert.Equal(t, false, updatedMessage.IsPalindrome, "Updated palindrome status should match")
}

func TestDeleteMessage(t *testing.T) {
	// create a message
	createURL := host + "/messages"
	request := svc.MessageRequest{Content: "message to be deleted"}

	reqJSON, err := json.Marshal(request)
	assert.NoError(t, err, "Error marshalling JSON")

	resp, err := http.Post(createURL, "application/json", bytes.NewBuffer(reqJSON))
	assert.NoError(t, err, "Error sending POST request")
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Unexpected response status")

	var createdMessage svc.MessageResponse
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Error reading response body")
	err = json.Unmarshal(body, &createdMessage)
	assert.NoError(t, err, "Error unmarshalling response")

	// delete the created message
	deleteURL := host + "/messages/" + createdMessage.ID
	req, err := http.NewRequest(http.MethodDelete, deleteURL, nil)
	assert.NoError(t, err, "Error creating DELETE request")

	client := &http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err, "Error sending DELETE request")
	assert.Equal(t, http.StatusNoContent, resp.StatusCode, "Unexpected response status")

	// Try to retrieve the deleted message
	resp, err = http.Get(deleteURL)
	assert.NoError(t, err, "Error sending GET request")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Unexpected response status")
}

func TestListMessages(t *testing.T) {
	listURL := host + "/messages"
	resp, err := http.Get(listURL)
	assert.NoError(t, err, "Error sending GET request")

	var messages []svc.MessageResponse
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Error reading response body")
	err = json.Unmarshal(body, &messages)
	assert.NoError(t, err, "Error unmarshalling response")

	// Create a message
	createURL := host + "/messages"
	request := svc.MessageRequest{Content: "first message"}

	reqJSON, err := json.Marshal(request)
	assert.NoError(t, err, "Error marshalling JSON")

	resp, err = http.Post(createURL, "application/json", bytes.NewBuffer(reqJSON))
	assert.NoError(t, err, "Error sending POST request")
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Unexpected response status")

	// List messages again
	resp, err = http.Get(listURL)
	assert.NoError(t, err, "Error sending GET request")

	body, err = io.ReadAll(resp.Body)
	assert.NoError(t, err, "Error reading response body")
	err = json.Unmarshal(body, &messages)
	assert.NoError(t, err, "Error unmarshalling response")
	assert.NotEmpty(t, messages, "Messages list should not be empty")
}
