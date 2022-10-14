package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4/middleware"
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
