package utilities

import (
	"context"
	"os"
	"os/signal"
	"time"
)

var quit = make(chan os.Signal, 1) // for unit test to work

// Stop individual services before graceful shutdown with timeout
func StopWaitWrapper(stop func(ctx context.Context), wait time.Duration) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		stop(ctx)
	}
}

// Graceful waits for interrupt signal to gracefully shutdown other services such as database, logger and scheduler before
// shutting down the Echo server.
// Use a buffered channel to avoid missing signals as recommended for signal.Notify
func Graceful(stop func()) {

	signal.Notify(quit, os.Interrupt)
	<-quit
	stop()
}
