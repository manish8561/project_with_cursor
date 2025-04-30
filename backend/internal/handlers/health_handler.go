package handlers

import (
	"backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health check handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check godoc
// @Summary Check server health
// @Description Returns the server's health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthResponse{
		Status: "ok",
	})
}
