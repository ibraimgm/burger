package app

import (
	"context"
)

// Server is the base interface for generic servers.
type Server interface {
	// Start starts the server, asynchronously
	Start(ctx context.Context)

	// Stop sends an async signal to terminate the server
	Stop()

	// Synchronously wait until the server shutdown
	Wait()
}
