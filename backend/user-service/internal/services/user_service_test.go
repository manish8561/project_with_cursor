package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"user-service/internal/models"
)

// MockKafkaPublisher is a mock implementation of KafkaPublisher for testing
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

func TestKafkaPublisher_PublishUserUpdated_Success(t *testing.T) {
	mockPublisher := new(MockKafkaPublisher)

	// Mock the Kafka publisher to expect a call
	ctx := context.Background()
	event := models.UserEvent{
		EventID:   primitive.NewObjectID().Hex(),
		EventType: "user.updated.v1",
		Timestamp: time.Now().UTC(),
		UserID:    "123",
		Email:     "test@example.com",
		Name:      "Test User",
		Status:    "active",
		Role:      "customer",
	}

	mockPublisher.On("PublishUserUpdated", ctx, mock.Anything).Return(nil)

	err := mockPublisher.PublishUserUpdated(ctx, event)

	assert.NoError(t, err)
	mockPublisher.AssertExpectations(t)
}

func TestKafkaPublisher_PublishUserDeleted_Success(t *testing.T) {
	mockPublisher := new(MockKafkaPublisher)

	// Mock the Kafka publisher to expect a call
	ctx := context.Background()
	event := models.UserEvent{
		EventID:   primitive.NewObjectID().Hex(),
		EventType: "user.deleted.v1",
		Timestamp: time.Now().UTC(),
		UserID:    "123",
	}

	mockPublisher.On("PublishUserDeleted", ctx, mock.Anything).Return(nil)

	err := mockPublisher.PublishUserDeleted(ctx, event)

	assert.NoError(t, err)
	mockPublisher.AssertExpectations(t)
}

func TestUserEvent_Structure(t *testing.T) {
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

	assert.Equal(t, "123", event.UserID)
	assert.Equal(t, "test@example.com", event.Email)
	assert.Equal(t, "Test User", event.Name)
	assert.Equal(t, "active", event.Status)
	assert.Equal(t, "customer", event.Role)
}
