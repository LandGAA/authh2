// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/login": {
            "get": {
                "description": "Вход email + пароль",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Вход",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth_internal_entity.User"
                        }
                    },
                    "400": {
                        "description": "Неправильные данные",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Ошибка аутентификации, неверный email или пароль",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка создания токенов",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Введите новый пароль",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Обновление пароля",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth_internal_entity.User"
                        }
                    },
                    "400": {
                        "description": "Некоректные данные",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Пользователь не найден",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Пользователь уже создан или ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "Получение всех пользователей (Без паролей)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Получить всех пользователей",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth_internal_entity.User"
                        }
                    },
                    "404": {
                        "description": "Пустая база данных",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users/email/{email}": {
            "get": {
                "description": "Получить пользователя по его Email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Получить пользователя по Email",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email пользователя",
                        "name": "email",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth_internal_entity.User"
                        }
                    },
                    "404": {
                        "description": "Пользователь не найден",
                        "schema": {
                            "type": "sting"
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "description": "Получишь пользователя по ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Получить пользователя по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID пользователя",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неправильный параметр",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Пользователя нет",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Удаление пользователя по его ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Удалить пользователя по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID пользователя",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пользователь удален",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неправильные параметры",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth_internal_entity.User": {
            "type": "object",
            "properties": {
                "createAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "69",
	Host:             "localhost:8081",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "Auth SQLite",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
