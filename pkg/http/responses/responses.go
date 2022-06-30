package responses

import (
	"encoding/json"
	"errors"
	"go-skeleton-rest-app/internal/models"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"gitlab.com/pos_malaysia/golib/logs"
)

const (
	ClientRequestID = "P-Request-Id"
)

const (
	SuccessCode                 = "S0000" // should return http status code: 200
	SystemErrorCode             = "E1001" // should return http status code: 200
	SDSorUSCRateLimitCode       = "E1002" // should return http status code: 429
	InvalidRequestCode          = "E1003" // should return http status code: 400
	RequestFailedValidationCode = "E1004" // should return http status code: 400
)

var responseMessages = map[string]string{
	SuccessCode:                 "Success",
	SystemErrorCode:             "System Error",
	SDSorUSCRateLimitCode:       "Rate Limit Hit at SDS or USC",
	InvalidRequestCode:          "Invalid Request Format or Data",
	RequestFailedValidationCode: "Request Failed Validation",
}

var responseMap = &sync.Map{}

func init() {
	for code, message := range responseMessages {
		if responseData, err := json.Marshal(&models.StandardJSONResponse{Code: code, Message: message}); err != nil {
			panic(err)
		} else {
			// load into map
			responseMap.Store(code, responseData)
		}
	}
}

// GetReponseMessageByCode returns a StandardJSONResponse struct with the associated message from the input code.
//
// If the input code is not available in the responseMessages map, then GetReponseMessageByCode returns nil with "invalid code" error.
func GetReponseMessageByCode(code string) (*models.StandardJSONResponse, error) {

	if responseData, ok := responseMap.Load(code); ok {
		result := &models.StandardJSONResponse{}
		if err := json.Unmarshal(responseData.([]byte), result); err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New("invalid code")
}

// ResponseWithError returns a HTTP status 400 together with associated error code, client request ID and server trace ID
func ResponseWithError(pRequestID, serverTraceID string, c echo.Context, responseErrorCode string) error {
	reply, err := GetReponseMessageByCode(responseErrorCode)

	if err != nil {
		logs.Error().Err(err).Caller().Msg("P-Request-Id : " + pRequestID + " server trace ID : " + serverTraceID)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set(ClientRequestID, pRequestID)
	c.Response().WriteHeader(http.StatusInternalServerError)
	return c.JSON(http.StatusInternalServerError, reply)
}
