package models

type (
	StandardJSONResponse struct {
		Code      int         `json:"code"`
		Message   string      `json:"message"`
		RequestID string      `json:"request_id"`
		Data      interface{} `json:"data"`
	}
)
