package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/pos_malaysia/golib/logs"
)

// For application healthcheck purpose. DO NOT MODIFY OR REMOVE THIS!

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 { Message: "I'm ok" }
// @Router / [get]
func HealthCheck(c echo.Context) error {

	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	// uncomment these lines to pass context to child function
	//ctx := c.Request().Context()
	//ctx = contextkeys.SetContextValue(ctx, contextkeys.CONTEXT_KEY_REQUEST_ID, requestID)

	logs.Info().Str("request ID", requestID).Str("handler", "HealthCheck").Send()

	reply := struct {
		Message string `json:"message"`
	}{
		Message: "I'm ok",
	}

	return c.JSON(http.StatusOK, reply)
}
