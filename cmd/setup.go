package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ziflex/lecho"
	"gitlab.com/pos_malaysia/golib/logs"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4/middleware"
)

var (
	// configure the logger's behaviour here
	logsConfig = logs.ConfigSet{}
	logger     = logs.Configure(logsConfig)
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func setupEcho() *echo.Echo {

	e := echo.New()

	// Setup Echo to use our logger
	e.Logger = lecho.New(logger) // Echo adapter for Zerolog

	e.Validator = &CustomValidator{validator: validator.New()}

	// Setup Echo's middleware
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(5) * time.Second,
	}))
	// log every request
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORS())

	return e
}
