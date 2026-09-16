package biz

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGreeterRepo is a mock implementation of GreeterRepo for testing
type MockGreeterRepo struct {
	mock.Mock
}

func (m *MockGreeterRepo) Save(ctx context.Context, g *Greeter) (*Greeter, error) {
	args := m.Called(ctx, g)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Greeter), args.Error(1)
}

func (m *MockGreeterRepo) Update(ctx context.Context, g *Greeter) (*Greeter, error) {
	args := m.Called(ctx, g)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Greeter), args.Error(1)
}

func (m *MockGreeterRepo) FindByID(ctx context.Context, id int64) (*Greeter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Greeter), args.Error(1)
}

func (m *MockGreeterRepo) ListByHello(ctx context.Context, hello string) ([]*Greeter, error) {
	args := m.Called(ctx, hello)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Greeter), args.Error(1)
}

func (m *MockGreeterRepo) ListAll(ctx context.Context) ([]*Greeter, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Greeter), args.Error(1)
}

func TestGreeterUsecase_CreateGreeter_Success(t *testing.T) {
	mockRepo := new(MockGreeterRepo)
	logger := log.DefaultLogger
	usecase := NewGreeterUsecase(mockRepo, logger)

	ctx := context.Background()
	greeter := &Greeter{Hello: "Test"}

	expectedGreeter := &Greeter{Hello: "Test"}
	mockRepo.On("Save", ctx, greeter).Return(expectedGreeter, nil)

	result, err := usecase.CreateGreeter(ctx, greeter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test", result.Hello)
	mockRepo.AssertExpectations(t)
}

func TestGreeterUsecase_CreateGreeter_EmptyHello(t *testing.T) {
	mockRepo := new(MockGreeterRepo)
	logger := log.DefaultLogger
	usecase := NewGreeterUsecase(mockRepo, logger)

	ctx := context.Background()
	greeter := &Greeter{Hello: ""}

	expectedGreeter := &Greeter{Hello: ""}
	mockRepo.On("Save", ctx, greeter).Return(expectedGreeter, nil)

	result, err := usecase.CreateGreeter(ctx, greeter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "", result.Hello)
	mockRepo.AssertExpectations(t)
}

func TestGreeterUsecase_CreateGreeter_RepoError(t *testing.T) {
	mockRepo := new(MockGreeterRepo)
	logger := log.DefaultLogger
	usecase := NewGreeterUsecase(mockRepo, logger)

	ctx := context.Background()
	greeter := &Greeter{Hello: "Test"}

	mockRepo.On("Save", ctx, greeter).Return(nil, assert.AnError)

	result, err := usecase.CreateGreeter(ctx, greeter)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
