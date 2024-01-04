// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Direktorat Pengembangan Teknologi dan Sistem Informasi (DPTSI) - ITS",
            "url": "http://its.ac.id/dptsi",
            "email": "dptsi@its.ac.id"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "security": [
                    {
                        "CSRF Token": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication \u0026 Authorization"
                ],
                "summary": "Rute untuk mendapatkan link login melalui OpenID Connect",
                "responses": {
                    "200": {
                        "description": "Link login berhasil didapatkan",
                        "schema": {
                            "$ref": "#/definitions/responses.GeneralResponse"
                        }
                    },
                    "500": {
                        "description": "Terjadi kesalahan saat menghubungi provider OpenID Connect",
                        "schema": {
                            "$ref": "#/definitions/responses.GeneralResponse"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "delete": {
                "security": [
                    {
                        "Session": []
                    },
                    {
                        "CSRF Token": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication \u0026 Authorization"
                ],
                "summary": "Rute untuk logout",
                "responses": {
                    "200": {
                        "description": "Logout berhasil",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.GeneralResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "code": {
                                            "type": "integer"
                                        },
                                        "data": {
                                            "type": "string"
                                        },
                                        "message": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/auth/user": {
            "get": {
                "security": [
                    {
                        "Session": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication \u0026 Authorization"
                ],
                "summary": "Rute untuk mendapatkan data user yang sedang login",
                "responses": {
                    "200": {
                        "description": "Data user berhasil didapatkan",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.GeneralResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "code": {
                                            "type": "integer"
                                        },
                                        "data": {
                                            "allOf": [
                                                {
                                                    "$ref": "#/definitions/responses.User"
                                                },
                                                {
                                                    "type": "object",
                                                    "properties": {
                                                        "roles": {
                                                            "type": "array",
                                                            "items": {
                                                                "$ref": "#/definitions/responses.Role"
                                                            }
                                                        }
                                                    }
                                                }
                                            ]
                                        },
                                        "message": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/csrf-cookie": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "CSRF Protection"
                ],
                "summary": "Rute dummy untuk set CSRF-TOKEN cookie",
                "responses": {
                    "200": {
                        "description": "Cookie berhasil diset",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.GeneralResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "code": {
                                            "type": "integer"
                                        },
                                        "message": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "responses.GeneralResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 123
                },
                "data": {
                    "type": "object"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "responses.Role": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "mahasiswa"
                },
                "name": {
                    "type": "string",
                    "example": "Mahasiswa"
                },
                "permissions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "responses.User": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "5025201000@student.its.ac.id"
                },
                "name": {
                    "type": "string",
                    "example": "Mahasiswa ITS"
                },
                "picture": {
                    "type": "string",
                    "example": "https://my.its.ac.id/picture/00000000-0000-0000-0000-000000000000"
                },
                "preferred_username": {
                    "type": "string",
                    "example": "5025201000@student.its.ac.id"
                },
                "roles": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/responses.Role"
                    }
                },
                "sub": {
                    "type": "string",
                    "example": "00000000-0000-0000-0000-000000000000"
                }
            }
        }
    },
    "securityDefinitions": {
        "CSRF Token": {
            "type": "apiKey",
            "name": "x-csrf-token",
            "in": "header"
        },
        "Session": {
            "type": "apiKey",
            "name": "akademik_its_ac_id_session",
            "in": "cookie"
        }
    },
    "externalDocs": {
        "description": "Dokumentasi Base Project",
        "url": "http://localhost:8080/doc/project"
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
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
