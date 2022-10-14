package main

import (
	"context"
	"net/http"
	"time"

	"go-skeleton-rest-app/internal/db"
	"go-skeleton-rest-app/internal/routes"
	"go-skeleton-rest-app/pkg/utilities"
	"log"
)

var echoPortNumber = "1234"

// @title Go Rest Skeleton Rest APP
// @version 1.0
// @description This is the Go Skeleton Rest APP

// @host localhost:1234
// @BasePath /
func main() {

	// setup Echo to use our golib/logs
	e := setupEcho()

	// Initialize and pass the zero logger to the routes.
	routes.InitRoutes(e)
	db.SetupDatabase()
	// Start server by spinning a goroutine so that it will become non-blocking
	go func() {
		if err := e.Start(":" + echoPortNumber); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down the server")
		}
	}()

	// Gracefully shutdown database, etc services before shutting down Echo server
	// with 10 seconds timeout
	utilities.Graceful(utilities.StopWaitWrapper(func(ctx context.Context) {
		db.DBPool().Close()
		e.Shutdown(ctx)
		log.Println("Go Rest Skeleton Rest APP server shutdown")

	}, time.Duration(10)*time.Second))
}
