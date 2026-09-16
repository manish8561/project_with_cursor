package data

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"

	"api-gateway/internal/biz"
)

func TestGreeterRepo_Save_Success(t *testing.T) {
	// Create a mock Data instance
	mockData := &Data{}
	logger := log.DefaultLogger
	repo := NewGreeterRepo(mockData, logger)

	ctx := context.Background()
	greeter := &biz.Greeter{Hello: "Test"}

	result, err := repo.Save(ctx, greeter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test", result.Hello)
}

func TestGreeterRepo_Update_Success(t *testing.T) {
	// Create a mock Data instance
	mockData := &Data{}
	logger := log.DefaultLogger
	repo := NewGreeterRepo(mockData, logger)

	ctx := context.Background()
	greeter := &biz.Greeter{Hello: "Updated"}

	result, err := repo.Update(ctx, greeter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated", result.Hello)
}

func TestGreeterRepo_FindByID_Success(t *testing.T) {
	// Create a mock Data instance
	mockData := &Data{}
	logger := log.DefaultLogger
	repo := NewGreeterRepo(mockData, logger)

	ctx := context.Background()
	id := int64(123)

	result, err := repo.FindByID(ctx, id)

	assert.NoError(t, err)
	assert.Nil(t, result) // Current implementation returns nil
}

func TestGreeterRepo_ListByHello_Success(t *testing.T) {
	// Create a mock Data instance
	mockData := &Data{}
	logger := log.DefaultLogger
	repo := NewGreeterRepo(mockData, logger)

	ctx := context.Background()
	hello := "Test"

	result, err := repo.ListByHello(ctx, hello)

	assert.NoError(t, err)
	assert.Nil(t, result) // Current implementation returns nil
}

func TestGreeterRepo_ListAll_Success(t *testing.T) {
	// Create a mock Data instance
	mockData := &Data{}
	logger := log.DefaultLogger
	repo := NewGreeterRepo(mockData, logger)

	ctx := context.Background()

	result, err := repo.ListAll(ctx)

	assert.NoError(t, err)
	assert.Nil(t, result) // Current implementation returns nil
}
