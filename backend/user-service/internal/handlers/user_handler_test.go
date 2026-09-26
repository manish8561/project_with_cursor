package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"user-service/internal/handlers"
	"user-service/internal/middleware"
	"user-service/internal/models"
)

// MockUserService is a mock implementation of UserService for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) ListUsers(page, size int) (*models.UserListResponse, error) {
	args := m.Called(page, size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserListResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(id string, req models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpsertUserProfileFromEvent(event models.UserEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockUserService) DeleteUserProfileFromEvent(event models.UserEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func setupTestRouter(mockUserService *MockUserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Create a mock logger
	logger, _ := zap.NewDevelopment()
	userHandler := handlers.NewUserHandler(mockUserService, logger)

	// Add middleware to set user ID context for testing (bypass auth)
	r.Use(func(c *gin.Context) {
		c.Set(middleware.ContextUserIDKey, "123")
		c.Next()
	})

	api := r.Group("/api/users")
	{
		api.GET("/me", userHandler.GetCurrentUser)
		api.GET("/profile/:id", userHandler.GetUserByID)
		api.GET("/list", userHandler.ListUsers)
		api.PUT("/profile/:id", userHandler.UpdateUser)
		api.DELETE("/profile/:id", userHandler.DeleteUser)
	}

	return r
}

func TestGetCurrentUser_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	r := setupTestRouter(mockUserService)

	expectedUser := &models.User{
		ID:    "123",
		Name:  "testuser",
		Email: "test@example.com",
	}

	mockUserService.On("GetUserByID", "123").Return(expectedUser, nil)

	req, _ := http.NewRequest("GET", "/api/users/me", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUserService.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	r := setupTestRouter(mockUserService)

	expectedUser := &models.User{
		ID:    "123",
		Name:  "testuser",
		Email: "test@example.com",
	}

	mockUserService.On("GetUserByID", "123").Return(expectedUser, nil)

	req, _ := http.NewRequest("GET", "/api/users/profile/123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, response.ID)

	mockUserService.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	r := setupTestRouter(mockUserService)

	updateReq := models.UpdateUserRequest{
		Name:  "updateduser",
		Email: "updated@example.com",
	}

	expectedUser := &models.User{
		ID:    "123",
		Name:  "updateduser",
		Email: "updated@example.com",
	}

	mockUserService.On("UpdateUser", "123", updateReq).Return(expectedUser, nil)

	reqBody, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", "/api/users/profile/123", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Name, response.Name)

	mockUserService.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	r := setupTestRouter(mockUserService)

	mockUserService.On("DeleteUser", "123").Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/users/profile/123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUserService.AssertExpectations(t)
}

func TestListUsers_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	r := setupTestRouter(mockUserService)

	expectedResponse := &models.UserListResponse{
		Users: []models.User{
			{ID: "1", Name: "user1", Email: "user1@example.com"},
			{ID: "2", Name: "user2", Email: "user2@example.com"},
		},
		Page:       1,
		Size:       10,
		TotalCount: 2,
		Total:      2,
	}

	mockUserService.On("ListUsers", 1, 10).Return(expectedResponse, nil)

	req, _ := http.NewRequest("GET", "/api/users/list?page=1&size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.UserListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response.Users))

	mockUserService.AssertExpectations(t)
}

func TestGetUserByID_RejectsMismatchedUserID(t *testing.T) {
	mockUserService := new(MockUserService)
	r := setupTestRouter(mockUserService)

	req, _ := http.NewRequest("GET", "/api/users/profile/other-user", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.JSONEq(t, `{"error":"forbidden"}`, w.Body.String())
}

func TestUpdateUser_RejectsInvalidJSON(t *testing.T) {
	mockUserService := new(MockUserService)
	r := setupTestRouter(mockUserService)

	req, _ := http.NewRequest("PUT", "/api/users/profile/123", bytes.NewBufferString(`{"name":`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestListUsers_DefaultsInvalidPaginationValues(t *testing.T) {
	mockUserService := new(MockUserService)
	r := setupTestRouter(mockUserService)

	expectedResponse := &models.UserListResponse{Users: []models.User{{ID: "123", Name: "testuser", Email: "test@example.com"}}, Page: 1, Size: 10, TotalCount: 1, Total: 1}
	mockUserService.On("ListUsers", 1, 10).Return(expectedResponse, nil)

	req, _ := http.NewRequest("GET", "/api/users/list?page=0&size=9999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUserService.AssertExpectations(t)
}

func TestDeleteUser_RejectsMissingUserIDParam(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	handler := handlers.NewUserHandler(new(MockUserService), logger)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/api/users/profile/", nil)
	c.Set(middleware.ContextUserIDKey, "123")
	c.Params = gin.Params{{Key: "id", Value: ""}}

	handler.DeleteUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"user ID is required"}`, w.Body.String())
}
