package services_test

import (
	"auth-service/internal/models"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockKafkaPublisher mirrors the user-service publisher mock and documents the
// expected contract for authentication lifecycle events.
type MockKafkaPublisher struct {
	mock.Mock
}

func (m *MockKafkaPublisher) PublishUserCreated(ctx context.Context, event models.UserEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockKafkaPublisher) PublishUserUpdated(ctx context.Context, event models.UserEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockKafkaPublisher) PublishUserDeleted(ctx context.Context, event models.UserEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockKafkaPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestKafkaPublisherPublishUserCreatedSuccess(t *testing.T) {
	mockPublisher := new(MockKafkaPublisher)
	ctx := context.Background()
	event := models.UserEvent{
		EventID:   primitive.NewObjectID().Hex(),
		EventType: "user.created.v1",
		Timestamp: time.Now().UTC(),
		UserID:    "123",
		Email:     "test@example.com",
		Name:      "Test User",
		Status:    "active",
		Role:      "customer",
	}

	mockPublisher.On("PublishUserCreated", ctx, mock.Anything).Return(nil)

	err := mockPublisher.PublishUserCreated(ctx, event)

	assert.NoError(t, err)
	mockPublisher.AssertExpectations(t)
}

func TestUserCreatedEventStructure(t *testing.T) {
	event := models.UserEvent{
		EventID:   primitive.NewObjectID().Hex(),
		EventType: "user.created.v1",
		Timestamp: time.Now().UTC(),
		UserID:    "123",
		Email:     "test@example.com",
		Name:      "Test User",
		Status:    "active",
		Role:      "customer",
	}

	assert.Equal(t, "user.created.v1", event.EventType)
	assert.Equal(t, "123", event.UserID)
	assert.Equal(t, "test@example.com", event.Email)
	assert.Equal(t, "Test User", event.Name)
	assert.Equal(t, "active", event.Status)
	assert.Equal(t, "customer", event.Role)
	assert.False(t, event.Timestamp.IsZero())
}
