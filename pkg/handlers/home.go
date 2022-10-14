package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {

	serverTraceID := c.Response().Header().Get(echo.HeaderXRequestID)

	// uncomment these lines to pass context to child function
	//ctx := c.Request().Context()
	//ctx = contextkeys.SetContextValue(ctx, contextkeys.CONTEXT_KEY_SERVER_TRACE_ID, serverTraceID)

	log.Println("home")
	log.Println("server trace ID : " + serverTraceID)

	return c.HTML(http.StatusOK, "<h1>home</h1>")
}
