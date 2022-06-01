package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/ziflex/lecho"
	"gitlab.com/pos_malaysia/golib/env"
	"gitlab.com/pos_malaysia/golib/logs"
	"gitlab.com/pos_malaysia/golib/redistools"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4/middleware"

	"gitlab.com/pos_malaysia/golib/aws"
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

func setupAWS() {
	region := env.Get("AWS_REGION")

	aws.InitAWS(region)
}

func setupRedis() *redis.Client {
	redisHost := env.Get("REDIS_HOST")
	redisPort := env.Get("REDIS_PORT")
	redisUsername := env.Get("REDIS_USERNAME")
	redisPassword := env.Get("REDIS_PASSWORD")
	redisDatabaseNumber, err := strconv.Atoi(env.Get("REDIS_DATABASE_NUMBER"))
	if err != nil {
		logs.Fatal().Msg("REDIS_DATABASE_NUMBER conversion to int error : " + err.Error())
	}

	redisMinimumIdleConns, err := strconv.Atoi(env.Get("REDIS_MIN_IDLE_CONNS"))
	if err != nil {
		logs.Fatal().Msg("REDIS_MIN_IDLE_CONNS conversion to int error : " + err.Error())
	}

	redisTimeout, err := strconv.Atoi(env.Get("REDIS_TIMEOUT"))
	if err != nil {
		logs.Fatal().Msg("REDIS_TIMEOUT conversion to int error : " + err.Error())
	}

	return redistools.InitRedisClient(redisHost, redisPort, redisUsername, redisPassword,
		redisDatabaseNumber, redisMinimumIdleConns, time.Duration(redisTimeout)*time.Second)
}
