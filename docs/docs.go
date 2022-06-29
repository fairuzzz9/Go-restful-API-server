// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Go Skeleton Rest APP Developer",
            "email": "boojiun@pos.com.my"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/healthcheck": {
            "get": {
                "description": "get the status of server.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Show the status of server.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Request ID generated by client",
                        "name": "P-Request-Id",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.StandardJSONResponse"
                        },
                        "headers": {
                            "string": {
                                "type": "string",
                                "description": "P-Request-Id"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.StandardJSONResponse": {
            "type": "object",
            "properties": {
                "client_request_id": {
                    "type": "string"
                },
                "code": {
                    "type": "string"
                },
                "data": {},
                "message": {
                    "type": "string"
                },
                "server_trace_id": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:1234",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Go Rest Skeleton Rest APP",
	Description:      "This is the Go Skeleton Rest APP",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
