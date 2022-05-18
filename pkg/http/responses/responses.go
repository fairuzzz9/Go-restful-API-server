package responses

import (
	"encoding/json"
	"errors"
	"sync"
)

type JSONResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	SuccessCode = 0
)

var responseMessages = map[int]string{
	SuccessCode: "success",
}

var responseMap = &sync.Map{}

func init() {
	for code, message := range responseMessages {
		if responseData, err := json.Marshal(&JSONResponse{Code: code, Message: message}); err != nil {
			panic(err)
		} else {
			// load into map
			responseMap.Store(code, responseData)
		}
	}
}

// GetReponseMessageByCode returns a JSONResponse struct with the associated message from the input code.
//
// If the input code is not available in the responseMessages map, then GetReponseMessageByCode returns nil with "invalid code" error.
func GetReponseMessageByCode(code int) (*JSONResponse, error) {

	if responseData, ok := responseMap.Load(code); ok {
		result := &JSONResponse{}
		if err := json.Unmarshal(responseData.([]byte), result); err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New("invalid code")
}
