// Package main API Gateway.
//
// API Gateway for microservices architecture.
//
//	Schemes: http, https
//	Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- bearer
//
// swagger:meta
package main

import "github.com/swaggo/swag"

// @title API Gateway
// @version 1.0
// @description API Gateway for microservices architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func init() {
	docTemplate := `{
	"swagger": "2.0",
	"info": {
		"description": "API Gateway for microservices architecture",
		"title": "API Gateway",
		"termsOfService": "http://swagger.io/terms/",
		"contact": {
			"name": "API Support",
			"url": "http://www.swagger.io/support",
			"email": "support@swagger.io"
		},
		"license": {
			"name": "Apache 2.0",
			"url": "http://www.apache.org/licenses/LICENSE-2.0.html"
		},
		"version": "1.0"
	},
	"host": "localhost:8080",
	"basePath": "/",
	"schemes": [
		"http",
		"https"
	],
	"consumes": [
		"application/json"
	],
	"produces": [
		"application/json"
	],
	"securityDefinitions": {
		"BearerAuth": {
			"type": "apiKey",
			"name": "Authorization",
			"in": "header",
			"description": "Type \"Bearer\" followed by a space and JWT token."
		}
	}
}`
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate:  docTemplate,
	})
}