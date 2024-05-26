package server

import "os"

// Runner represents an interface for managing the server's lifecycle and services.
type Runner interface {
	// Start starts the server.
	Start() error
	// Stop stops the server gracefully.
	Stop(stopCh chan os.Signal)
	// RegisterServices registers services required by the server.
	RegisterServices()
}
