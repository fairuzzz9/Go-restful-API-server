package routes

import (
	"go-skeleton-rest-app/pkg/handlers"

	"github.com/labstack/echo/v4"
)

// define your URL path here
const (
	PATH_HOME         = "/"
	PATH_HEALTH_CHECK = "/healthcheck"
)

func InitRoutes(e *echo.Echo) {

	// divide each handler to individual .go file.

	e.GET(PATH_HOME, handlers.Home)                // home.go
	e.GET(PATH_HEALTH_CHECK, handlers.HealthCheck) // healthcheck.go

}
