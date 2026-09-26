package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type stubLogger struct{}

func (stubLogger) Info(string, ...zap.Field)  {}
func (stubLogger) Warn(string, ...zap.Field)  {}
func (stubLogger) Error(string, ...zap.Field) {}
func (stubLogger) Sync() error                { return nil }

func TestSplitBrokersTrimsEmptyBrokerEntries(t *testing.T) {
	assert.Equal(t, []string{"broker-1:9092", "broker-2:9092"}, splitBrokers(" broker-1:9092, ,broker-2:9092 "))
	assert.Empty(t, splitBrokers(" , , "))
}

func TestNewUserEventConsumer(t *testing.T) {
	consumer, err := NewUserEventConsumer(
		"broker-1:9092, broker-2:9092",
		"user-group",
		"user-client",
		"user.created.v1",
		"user.updated.v1",
		"user.deleted.v1",
		&UserService{},
		stubLogger{},
	)
	require.NoError(t, err)
	require.NotNil(t, consumer)
	assert.Len(t, consumer.readers, 3)
	assert.Equal(t, "user.created.v1", consumer.topicUserCreated)
	assert.Equal(t, "user.updated.v1", consumer.topicUserUpdated)
	assert.Equal(t, "user.deleted.v1", consumer.topicUserDeleted)

	consumer, err = NewUserEventConsumer("broker-1:9092", "", "user-client", "user.created.v1", "user.updated.v1", "user.deleted.v1", &UserService{}, stubLogger{})
	assert.ErrorContains(t, err, "kafka group id is required")
	assert.Nil(t, consumer)

	consumer, err = NewUserEventConsumer(" , , ", "user-group", "user-client", "user.created.v1", "user.updated.v1", "user.deleted.v1", &UserService{}, stubLogger{})
	require.NoError(t, err)
	assert.Nil(t, consumer)
}

func TestUserEventConsumerNoopsForEmptyState(t *testing.T) {
	var consumer *UserEventConsumer
	assert.NoError(t, consumer.Close())
	consumer = &UserEventConsumer{}
	assert.NoError(t, consumer.Close())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	consumer = &UserEventConsumer{}
	consumer.Start(ctx)
}
