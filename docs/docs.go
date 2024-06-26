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
        "/createBulkUpload": {
            "post": {
                "security": [
                    {
                        "bearerToken": []
                    }
                ],
                "description": "Upload a CSV file containing grocery items and create them in bulk",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Grocery"
                ],
                "summary": "Create grocery items in bulk from a CSV file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "CSV file containing grocery items",
                        "name": "csvfile",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/createItem": {
            "post": {
                "security": [
                    {
                        "bearerToken": []
                    }
                ],
                "description": "Create a new grocery item with product information and an optional image",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Grocery"
                ],
                "summary": "Create a new grocery item",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Image file (.jpg)",
                        "name": "Image",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "name": "Brand",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "Category",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "CountryofOrigin",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "name": "ItemPackageQuantity",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "Manufacturer",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "PackageInformation",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "name": "Price",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "ProductName",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "name": "Vegeterian",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "name": "Weight",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.CreationSuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/deleteItem": {
            "delete": {
                "security": [
                    {
                        "bearerToken": []
                    }
                ],
                "description": "Delete a grocery item by providing the ProductId",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Grocery"
                ],
                "summary": "Delete a grocery item",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ProductId of the grocery item to be deleted",
                        "name": "ProductId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/getItemById": {
            "get": {
                "security": [
                    {
                        "bearerToken": []
                    }
                ],
                "description": "Retrieve a grocery item by providing the DocumentId",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Grocery"
                ],
                "summary": "Get a grocery item by DocumentId",
                "parameters": [
                    {
                        "type": "string",
                        "description": "DocumentId of the grocery item to be retrieved",
                        "name": "ProductId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/listGrocery": {
            "get": {
                "security": [
                    {
                        "bearerToken": []
                    }
                ],
                "description": "Retrieve a list of grocery items with optional filters like productName, priceGreaterThan, priceLessThan, and category. Supports pagination.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Grocery"
                ],
                "summary": "List grocery items with optional filters and pagination",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filter by product name",
                        "name": "productName",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by price greater than this value",
                        "name": "priceGreaterThanOrEqualTo",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by price less than this value",
                        "name": "priceLessThanOrEqualTo",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by price less than this value",
                        "name": "priceEqualTo",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by category",
                        "name": "category",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page number for pagination",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ListGrocerySuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Login to access APIs",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authenticatuion"
                ],
                "summary": "User login",
                "parameters": [
                    {
                        "description": "Credentials",
                        "name": "RequestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Credentials"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponseJWT"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/updateItem": {
            "put": {
                "security": [
                    {
                        "bearerToken": []
                    }
                ],
                "description": "Update a grocery item by providing the ProductId and new product information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Grocery"
                ],
                "summary": "Update a grocery item",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ProductId of the grocery item to be updated",
                        "name": "ProductId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Optional: Updated image file (.jpg)",
                        "name": "Image",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "name": "Brand",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "Category",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "CountryofOrigin",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "name": "ItemPackageQuantity",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "Manufacturer",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "PackageInformation",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "name": "Price",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "ProductName",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "name": "Vegeterian",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "name": "Weight",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UpdateSuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.CreationSuccessResponse": {
            "type": "object",
            "properties": {
                "Product created with id": {
                    "type": "string"
                },
                "productData": {
                    "$ref": "#/definitions/models.Product"
                }
            }
        },
        "models.Credentials": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.ErrorResponse": {
            "type": "object",
            "properties": {
                "Error": {
                    "type": "string"
                }
            }
        },
        "models.ListGrocerySuccessResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "next page": {},
                "number of products": {
                    "type": "integer"
                },
                "page": {
                    "type": "integer"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Product"
                    }
                }
            }
        },
        "models.Product": {
            "type": "object",
            "required": [
                "Brand",
                "Category",
                "CountryofOrigin",
                "ItemPackageQuantity",
                "Manufacturer",
                "PackageInformation",
                "Price",
                "ProductName",
                "Weight"
            ],
            "properties": {
                "Brand": {
                    "type": "string"
                },
                "Category": {
                    "type": "string"
                },
                "CountryofOrigin": {
                    "type": "string"
                },
                "ImageURL": {
                    "type": "string"
                },
                "ItemPackageQuantity": {
                    "type": "integer",
                    "minimum": 1
                },
                "Manufacturer": {
                    "type": "string"
                },
                "PackageInformation": {
                    "type": "string"
                },
                "Price": {
                    "type": "integer",
                    "minimum": 1
                },
                "ProductID": {
                    "type": "string"
                },
                "ProductName": {
                    "type": "string"
                },
                "ThumbnailURL": {
                    "type": "string"
                },
                "Vegeterian": {
                    "type": "boolean"
                },
                "Weight": {
                    "type": "string"
                }
            }
        },
        "models.SuccessResponse": {
            "type": "object",
            "properties": {
                "Success": {
                    "type": "string"
                }
            }
        },
        "models.SuccessResponseJWT": {
            "type": "object",
            "properties": {
                "Success": {
                    "type": "string"
                },
                "Token": {
                    "type": "string"
                }
            }
        },
        "models.UpdateSuccessResponse": {
            "type": "object",
            "properties": {
                "Product updated with id": {
                    "type": "string"
                },
                "productData": {
                    "$ref": "#/definitions/models.Product"
                }
            }
        }
    },
    "securityDefinitions": {
        "bearerToken": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "us-central1-capstone-412502.cloudfunctions.net",
	BasePath:         "",
	Schemes:          []string{"https"},
	Title:            "Grocery Management System",
	Description:      "Grocery Management API (v1) is a RESTful service that allows you to manage grocery data.\nThis API provides a set of endpoints to create, retrieve, update, and delete grocery records. It is designed to streamline the processes of managing grocery information in a grocery shop.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
