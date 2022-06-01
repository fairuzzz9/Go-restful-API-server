package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
	"track-and-trace-api-server/pkg/usc"
	"track-and-trace-api-server/pkg/utilities"

	"context"

	"github.com/labstack/echo/v4"
	"gitlab.com/pos_malaysia/golib/contextkeys"
	"gitlab.com/pos_malaysia/golib/env"
	"gitlab.com/pos_malaysia/golib/logs"
	"gitlab.com/pos_malaysia/golib/redistools"
)

var (
	uscTrackAndTraceURL = env.Get("USC_URL") + "as2corporate/v2trackntracewebapijson/v1/api/Details?id=%s&Culture=%s"
)

func getAPIURL(id, culture string) string {
	return fmt.Sprintf(uscTrackAndTraceURL, id, culture)
}

func connoteIDKey(id string) string {
	return fmt.Sprintf("connoteID:%s", id)
}

// TrackAndTraceDetails returns the tracking details base on the given ID(tracking number) and Culture(EN/MS)
//
// @Summary Show the tracking details of the given ID(tracking number) and Culture(EN/MS)
// @Description Get the tracking details of a given ID(tracking number) and Culture(EN/MS)
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 string
// @Failure 400 {object} echo.NewHTTPError "connoteID cannot be empty string"
// @Failure 500 {object} echo.NewHTTPError "GetAPIURL req error"
// @Failure 500 {object} echo.NewHTTPError "GetAPIURL resp error"
// @Failure 500 {object} echo.NewHTTPError "redistools.SetJSON error"
// @Router /trackandtracedetails?id={connote_id}&Culture={En/Ms} [get]
func TrackAndTraceDetails(c echo.Context) error {

	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	// pass request ID to child context
	ctx := c.Request().Context()
	ctx = contextkeys.SetContextValue(ctx, contextkeys.CONTEXT_KEY_REQUEST_ID, requestID)

	// Connote ID from path `trackandtracedetails/:id`
	id := c.QueryParam("id")

	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("connoteID cannot be empty string"))
	}

	culture := c.QueryParam("Culture")

	if culture == "" {
		culture = "En" // default to En
	}

	var result string
	// Check if this connote data exist in Redis
	redisErr := redistools.GetJSON(ctx, connoteIDKey(id), &result)

	if redisErr != nil { // no Redis key for this connoteID

		logs.Info().Msg("no Redis key for this connoteID : " + connoteIDKey(id))

		// GET from USC/SDS API gateway
		getAPIURL := getAPIURL(id, culture)
		client := &http.Client{}
		req, err := http.NewRequestWithContext(ctx, "GET", getAPIURL, nil)

		if err != nil {
			logs.Error().Err(err).Send()
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("GetAPIURL req error"))
		}

		accessToken, err := getUSCAccessToken(ctx)
		if err != nil {
			logs.Error().Err(err).Send()
			return err
		}

		req.Header.Add("Authorization", "Bearer "+accessToken)
		resp, err := client.Do(req)

		if err != nil {
			logs.Error().Err(err).Send()
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("GetAPIURL resp error"))
		}
		defer resp.Body.Close()

		var body bytes.Buffer
		_, _ = io.Copy(&body, resp.Body) // io.Copy is safer and faster than ioutil.ReadAll

		result = body.String()

		// create Redis key for this connoteID
		redisErr := redistools.SetJSON(ctx, result, connoteIDKey(id), time.Duration(5)*time.Minute)

		if redisErr != nil {
			logs.Error().Err(err).Send()
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("redistools.SetJSON error"))
		}
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	return c.String(http.StatusOK, result)

}

// getUSCAccessToken returns the access token to USC
func getUSCAccessToken(ctx context.Context) (string, error) {

	const uscRedisKey = "usc:access_token"

	// check if access token is available in Redis
	redisGetStatusCmd := redistools.RedisClient().Get(ctx, uscRedisKey)

	if redisGetStatusCmd.Err() != nil {
		// access token not available in Redis. Get from USC.
		accessToken, ExpiresIn, err := usc.GetAccessToken(ctx)

		if err != nil {
			logs.Error().Err(err).Send()
			return "", echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("usc.GetAccessToken error"))
		}

		// create Redis key for access token
		// rotate by Rot18 before storing into Redis
		redisSetStatusCmd := redistools.RedisClient().Set(ctx, uscRedisKey, utilities.Rot18(accessToken), time.Duration(ExpiresIn)*time.Second)

		if redisSetStatusCmd.Err() != nil {
			logs.Error().Err(err).Send()
			return "", echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("redistools.RedisClient().Set error"))
		}

		return accessToken, nil

	}
	// get access token from Redis
	// rotate by Rot18 before returning the value
	return utilities.Rot18(redisGetStatusCmd.Val()), nil

}
