package http

import (
	"context"
	"github.com/gharsallahmoez/palindrome/config"
	"github.com/gharsallahmoez/palindrome/server"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Runner represents a server runner with configurations, services, and an HTTP server.
type Runner struct {
	MessageService *MessageService
	Config         *config.Server
	Server         http.Server
	mux.Router
}

// NewRunner creates a new instance of the server runner.
func NewRunner(conf *config.Server, messageService *MessageService) server.Runner {
	return &Runner{
		MessageService: messageService,
		Config:         conf,
		Server:         http.Server{},
		Router:         mux.Router{},
	}
}

// Start starts the server
func (r *Runner) Start() error {
	muxWithMiddlewares := http.TimeoutHandler(&r.Router, r.Config.Timeout*time.Second, "Timeout!")
	return http.ListenAndServe(":"+r.Config.Port, muxWithMiddlewares)
}

// Stop gracefully shuts down the server.
func (r *Runner) Stop(stopCh chan os.Signal) {
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), r.Config.Timeout*time.Second)
	defer shutdownRelease()
	if err := r.Server.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("failed to shutdown the server : %v", err)
	}
}

// RegisterServices configures the handlers for every route.
func (r *Runner) RegisterServices() {
	// register message APIs
	r.Router.HandleFunc("/messages", r.MessageService.CreateMessageHandler).Methods(http.MethodPost)
	r.Router.HandleFunc("/messages", r.MessageService.ListMessageHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/messages/{id}", r.MessageService.GetMessageHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/messages/{id}", r.MessageService.UpdateMessageHandler).Methods(http.MethodPut)
	r.Router.HandleFunc("/messages/{id}", r.MessageService.DeleteMessageHandler).Methods(http.MethodDelete)
}
