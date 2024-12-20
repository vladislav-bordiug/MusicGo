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
        "/addsong": {
            "post": {
                "description": "Add song based on group and song provided as json.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "song"
                ],
                "summary": "Add song",
                "parameters": [
                    {
                        "description": "JSON with group and song",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AddDeleteRequestData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/deletesong": {
            "post": {
                "description": "Delete song based on group and song provided as json.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "song"
                ],
                "summary": "Delete song",
                "parameters": [
                    {
                        "description": "JSON with group and song",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AddDeleteRequestData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/editsong": {
            "post": {
                "description": "Edit song releaseDate, text and link based on group and song provided as json.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "song"
                ],
                "summary": "Edit song text",
                "parameters": [
                    {
                        "description": "JSON with group, song, releaseDate, text, and link",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.EditRequestData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/getdata": {
            "get": {
                "description": "Retrieve songs and their details with pagination based on the page and items and filtration based on group, song, releaseDate, text and link provided as query parameters.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Get all songs and their information with pagination",
                "parameters": [
                    {
                        "type": "integer",
                        "example": 1,
                        "description": "Current page",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "example": 10,
                        "description": "Number of elements on the page",
                        "name": "items",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "\"Muse\"",
                        "description": "Group",
                        "name": "group",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "\"Supermassive Black Hole\"",
                        "description": "Song name",
                        "name": "song",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "\"16.07.2006\"",
                        "description": "Release date in format DD.MM.YYYY",
                        "name": "releaseDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "\"Ooh baby, don't you know I suffer?\\nOoh baby, can you hear me moan?\\nYou caught me under false pretenses\\nHow long before you let me go?\\n\\nOoh\\nYou set my soul alight\\nOoh\\nYou set my soul alight\"",
                        "description": "Song text (multiline allowed)",
                        "name": "text",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "\"https://www.youtube.com/watch?v=Xsp3_a-PMTw\"",
                        "description": "Song link",
                        "name": "link",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AnswerData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/getsongtext": {
            "get": {
                "description": "Retrieve song text with pagination based on the group, song and couplet provided as query parameters.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "song"
                ],
                "summary": "Get songs text with pagination",
                "parameters": [
                    {
                        "type": "string",
                        "example": "\"Muse\"",
                        "description": "Group",
                        "name": "group",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "\"Supermassive Black Hole\"",
                        "description": "Song name",
                        "name": "song",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "example": 1,
                        "description": "Couplet",
                        "name": "couplet",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AnswerCoupletData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.AddDeleteRequestData": {
            "type": "object",
            "required": [
                "group",
                "song"
            ],
            "properties": {
                "group": {
                    "type": "string",
                    "example": "Muse"
                },
                "song": {
                    "type": "string",
                    "example": "Supermassive Black Hole"
                }
            }
        },
        "models.AnswerCoupletData": {
            "type": "object",
            "required": [
                "text"
            ],
            "properties": {
                "text": {
                    "type": "string",
                    "example": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?"
                }
            }
        },
        "models.AnswerData": {
            "type": "object",
            "required": [
                "items"
            ],
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.RowDbData"
                    }
                }
            }
        },
        "models.EditRequestData": {
            "type": "object",
            "required": [
                "group",
                "song"
            ],
            "properties": {
                "group": {
                    "type": "string",
                    "example": "Muse"
                },
                "link": {
                    "type": "string",
                    "example": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
                },
                "releaseDate": {
                    "type": "string",
                    "example": "16.07.2006"
                },
                "song": {
                    "type": "string",
                    "example": "Supermassive Black Hole"
                },
                "text": {
                    "type": "string",
                    "example": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
                }
            }
        },
        "models.RowDbData": {
            "type": "object",
            "required": [
                "group",
                "link",
                "releaseDate",
                "song",
                "text"
            ],
            "properties": {
                "group": {
                    "type": "string",
                    "example": "Muse"
                },
                "link": {
                    "type": "string",
                    "example": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
                },
                "releaseDate": {
                    "type": "string",
                    "example": "16.07.2006"
                },
                "song": {
                    "type": "string",
                    "example": "Supermassive Black Hole"
                },
                "text": {
                    "type": "string",
                    "example": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Go Music",
	Description:      "This is a sample server for music library",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
