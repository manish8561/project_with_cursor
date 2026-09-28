package handlers

import (
	"net/http"
	"strconv"

	"notification-service/internal/middleware"
	"notification-service/internal/services"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct{ service *services.NotificationService }

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) GetPreference(c *gin.Context) {
	preference, err := h.service.Preference(c.Request.Context(), c.GetString(middleware.UserIDKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load notification preference"})
		return
	}
	c.JSON(http.StatusOK, preference)
}

func (h *NotificationHandler) UpdatePreference(c *gin.Context) {
	var request struct {
		EmailEnabled *bool `json:"emailEnabled" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil || request.EmailEnabled == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "emailEnabled must be provided"})
		return
	}
	preference, err := h.service.UpdatePreference(c.Request.Context(), c.GetString(middleware.UserIDKey), *request.EmailEnabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notification preference"})
		return
	}
	c.JSON(http.StatusOK, preference)
}

func (h *NotificationHandler) History(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	response, err := h.service.History(c.Request.Context(), c.GetString(middleware.UserIDKey), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load notification history"})
		return
	}
	c.JSON(http.StatusOK, response)
}
