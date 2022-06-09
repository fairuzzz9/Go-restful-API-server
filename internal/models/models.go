package models

type (
	StandardJSONResponse struct {
		Code            string      `json:"code"`
		Message         string      `json:"message"`
		ClientRequestID string      `json:"client_request_id"`
		ServerTraceID   string      `json:"server_trace_id"`
		Data            interface{} `json:"data"`
	}
)
