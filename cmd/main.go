package main

import (
	"context"
	"net/http"
	"time"

	"track-and-trace-api-server/internal/routes"
	"track-and-trace-api-server/pkg/utilities"

	"gitlab.com/pos_malaysia/golib/logs"
)

var echoPortNumber = "1234"

// @title Track And Trace Backend
// @version 1.0
// @description This is the Track And Trace Backend server.

// @contact.name API Developer
// @contact.email boojiun@pos.com.my

// @host localhost:1234
// @BasePath /
// @schemes http
func main() {

	// setup Echo to use our golib/logs
	e := setupEcho()

	// Initialize and pass the zero logger to the routes.
	routes.InitRoutes(e)

	// setup Redis
	redisClient := setupRedis()

	//setup AWS session
	setupAWS()

	// Start server by spinning a goroutine so that it will become non-blocking
	go func() {
		if err := e.Start(":" + echoPortNumber); err != nil && err != http.ErrServerClosed {
			logs.Fatal().Msg("shutting down the server")
		}
	}()

	// Gracefully shutdown database, etc services before shutting down Echo server
	// with 10 seconds timeout
	utilities.Graceful(utilities.StopWaitWrapper(func(ctx context.Context) {

		redisClient.Close()
		logs.Close()
		e.Shutdown(ctx)
		logs.Info().Msg("Track And Trace API server shutdown")

	}, time.Duration(10)*time.Second))
}
