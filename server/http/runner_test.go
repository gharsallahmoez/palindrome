package http

import (
	"github.com/gharsallahmoez/palindrome/config"
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestStartAndStopServer(t *testing.T) {
	// config
	conf := &config.Server{
		Port:    "8080",
		Timeout: 10,
	}
	messageService := NewMessageService(nil)
	// Create a new Runner instance
	runner := NewRunner(conf, messageService)
	// Start the server
	go func() {
		err := runner.Start()
		assert.NoError(t, err)
	}()
	// Stop the server
	stopCh := make(chan os.Signal, 1)
	go func() {
		runner.Stop(stopCh)
	}()
	// Send SIGINT signal to stop the server
	stopCh <- syscall.SIGINT
	// Wait for a short duration to ensure the server has stopped
	time.Sleep(100 * time.Millisecond)
}

func TestRegisterServices(t *testing.T) {
	// Mock configuration
	conf := &config.Server{
		Port:    "8080",
		Timeout: 10,
	}
	messageService := NewMessageService(nil)
	// Create a new Runner instance
	runner := NewRunner(conf, messageService)
	// Register services
	runner.RegisterServices()
	// Check if handlers are registered
	assert.NotNil(t, messageService.GetMessageHandler)
	assert.NotNil(t, messageService.CreateMessageHandler)
	assert.NotNil(t, messageService.ListMessageHandler)
	assert.NotNil(t, messageService.UpdateMessageHandler)
	assert.NotNil(t, messageService.DeleteMessageHandler)
}
