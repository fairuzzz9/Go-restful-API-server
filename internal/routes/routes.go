package routes

import (
	"track-and-trace-api-server/pkg/handlers"

	"github.com/labstack/echo/v4"
)

// define your URL path here
const (
	PATH_HOME            = "/"
	PATH_HEALTH_CHECK    = "/healthcheck"
	PATH_TRACK_AND_TRACE = "/trackandtracedetails"
)

func InitRoutes(e *echo.Echo) {

	// divide each handler to individual .go file.

	// route to handler
	e.GET(PATH_HOME, handlers.Home)                            // home.go
	e.GET(PATH_HEALTH_CHECK, handlers.HealthCheck)             // healthcheck.go
	e.GET(PATH_TRACK_AND_TRACE, handlers.TrackAndTraceDetails) // trackandtrace.go

}
