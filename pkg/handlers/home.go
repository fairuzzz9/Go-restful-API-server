package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/pos_malaysia/golib/logs"
)

func Home(c echo.Context) error {

	serverTraceID := c.Response().Header().Get(echo.HeaderXRequestID)

	// uncomment these lines to pass context to child function
	//ctx := c.Request().Context()
	//ctx = contextkeys.SetContextValue(ctx, contextkeys.CONTEXT_KEY_SERVER_TRACE_ID, serverTraceID)

	logs.Info().Str("server request ID", serverTraceID).Str("handler", "Home").Send()

	return c.HTML(http.StatusOK, "<h1>home</h1>")
}
