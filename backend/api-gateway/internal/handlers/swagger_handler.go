package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	_ "github.com/swaggo/swag/example/basic/docs"
)

// SwaggerHandler handles API documentation
type SwaggerHandler struct{}

// NewSwaggerHandler creates a new swagger handler
func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// SetupSwagger sets up the swagger documentation routes
func (h *SwaggerHandler) SetupSwagger(router *gin.Engine) {
	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// API documentation index
	router.GET("/docs", h.GetDocsIndex)
}

// GetDocsIndex returns the main documentation page
func (h *SwaggerHandler) GetDocsIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "docs.html", gin.H{
		"title": "API Documentation",
		"services": []map[string]string{
			{
				"name": "Authentication Service",
				"description": "User authentication and registration",
				"endpoints": "POST /api/auth/login, POST /api/auth/register, POST /api/auth/validate",
			},
			{
				"name": "User Service",
				"description": "User profile management",
				"endpoints": "GET /api/users/profile/:id, GET /api/users/list, PUT /api/users/profile/:id",
			},
		},
	})
} 