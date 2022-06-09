package responses

import (
	"encoding/json"
	"errors"
	"go-skeleton-rest-app/internal/models"
	"sync"
)

const (
	SuccessCode = "S0000"
)

var responseMessages = map[string]string{
	SuccessCode: "success",
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
