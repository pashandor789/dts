// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/build": {
            "put": {
                "description": "Put build.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Put build",
                "parameters": [
                    {
                        "description": "cc",
                        "name": "build",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.Build"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/http.Success"
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    }
                }
            }
        },
        "/builds": {
            "get": {
                "description": "Retrieves a list of builds.",
                "produces": [
                    "application/json"
                ],
                "summary": "Get builds",
                "responses": {
                    "200": {
                        "description": "List of builds",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/types.Build"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "http.Error": {
            "type": "object",
            "properties": {
                "payload": {},
                "status": {
                    "type": "string"
                }
            }
        },
        "http.Success": {
            "type": "object",
            "properties": {
                "payload": {},
                "status": {
                    "type": "string"
                }
            }
        },
        "types.Build": {
            "type": "object",
            "properties": {
                "execute_script": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "init_script": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Task Batching Storage Coordinator API",
	Description:      "HTTP Tabasco",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}