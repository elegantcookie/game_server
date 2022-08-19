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
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/users": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get all users endpoint",
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": ""
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Create user endpoint",
                "parameters": [
                    {
                        "description": "create user struct",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.CreateUserDTO"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": ""
                    },
                    "400": {
                        "description": ""
                    },
                    "418": {
                        "description": "I'm a teapot",
                        "schema": {
                            "$ref": "#/definitions/apperror.AppError"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Delete user by id endpoint",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "asd",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "400": {
                        "description": ""
                    },
                    "418": {
                        "description": "I'm a teapot",
                        "schema": {
                            "$ref": "#/definitions/apperror.AppError"
                        }
                    }
                }
            },
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Partially update user endpoint",
                "parameters": [
                    {
                        "description": "update user struct",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.UpdateUserDTO"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "400": {
                        "description": ""
                    },
                    "418": {
                        "description": "I'm a teapot",
                        "schema": {
                            "$ref": "#/definitions/apperror.AppError"
                        }
                    }
                }
            }
        },
        "/api/users/auth": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get user by username and password. Needs for authorization",
                "parameters": [
                    {
                        "description": "create user struct",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.CreateUserDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": ""
                    },
                    "418": {
                        "description": "I'm a teapot",
                        "schema": {
                            "$ref": "#/definitions/apperror.AppError"
                        }
                    }
                }
            }
        },
        "/api/users/id/{id}": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get user by user id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": ""
                    },
                    "418": {
                        "description": "I'm a teapot",
                        "schema": {
                            "$ref": "#/definitions/apperror.AppError"
                        }
                    }
                }
            }
        },
        "/api/users/tickets": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tickets"
                ],
                "summary": "Add user by user id and user id",
                "parameters": [
                    {
                        "description": "ticket dto struct",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.TicketDTO"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": ""
                    },
                    "400": {
                        "description": ""
                    },
                    "418": {
                        "description": "I'm a teapot",
                        "schema": {
                            "$ref": "#/definitions/apperror.AppError"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tickets"
                ],
                "summary": "Delete user by user id and user id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "data",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "400": {
                        "description": ""
                    },
                    "418": {
                        "description": "I'm a teapot",
                        "schema": {
                            "$ref": "#/definitions/apperror.AppError"
                        }
                    }
                }
            }
        },
        "/api/users/tickets/free/{id}": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tickets"
                ],
                "summary": "Get user status",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Ticket ID",
                        "name": "ticket_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "404": {
                        "description": ""
                    },
                    "418": {
                        "description": "I'm a teapot",
                        "schema": {
                            "$ref": "#/definitions/apperror.AppError"
                        }
                    }
                }
            }
        },
        "/api/users/username/{username}": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get user by username endpoint",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": ""
                    },
                    "418": {
                        "description": "I'm a teapot",
                        "schema": {
                            "$ref": "#/definitions/apperror.AppError"
                        }
                    }
                }
            }
        },
        "/userapi/heartbeat": {
            "get": {
                "tags": [
                    "Metrics"
                ],
                "summary": "Heartbeat metric",
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "400": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "apperror.AppError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "developer_message": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "user.CreateUserDTO": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "user.GameTickets": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "game_type": {
                    "type": "string"
                },
                "tickets_of_gt": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "user.TicketDTO": {
            "type": "object",
            "properties": {
                "game_type": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "ticket_id": {
                    "type": "string"
                }
            }
        },
        "user.UpdateUserDTO": {
            "type": "object",
            "properties": {
                "has_free_ticket": {
                    "type": "boolean"
                },
                "id": {
                    "type": "string"
                },
                "tickets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/user.GameTickets"
                    }
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
