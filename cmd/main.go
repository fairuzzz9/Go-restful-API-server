package main

import (
	"context"
	"net/http"
	"time"

	"go-skeleton-rest-app/internal/routes"
	"go-skeleton-rest-app/pkg/utilities"

	"gitlab.com/pos_malaysia/golib/logs"
)

var echoPortNumber = "1234"

func main() {

	// setup Echo to use our golib/logs
	e, logger := setupEcho()

	// Initialize and pass the zero logger to the routes.
	routes.InitRoutes(e)

	// Start server by spinning a goroutine so that it will become non-blocking
	go func() {
		if err := e.Start(":" + echoPortNumber); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Msg("shutting down the server")
		}
	}()

	// Gracefully shutdown database, etc services before shutting down Echo server
	// with 10 seconds timeout
	utilities.Graceful(utilities.StopWaitWrapper(func(ctx context.Context) {

		logs.Close()
		e.Shutdown(ctx)
		logger.Info().Msg("Echo server shutdown")

	}, 10*time.Second))
}
