package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	v1 "api-gateway/api/helloworld/v1"
	"api-gateway/internal/biz"
)

// MockGreeterUsecase is a mock implementation of GreeterUsecaseInterface for testing
type MockGreeterUsecase struct {
	mock.Mock
}

func (m *MockGreeterUsecase) CreateGreeter(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	args := m.Called(ctx, g)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*biz.Greeter), args.Error(1)
}

func TestGreeterService_SayHello_Success(t *testing.T) {
	mockUsecase := new(MockGreeterUsecase)
	greeterService := NewGreeterService(mockUsecase)

	ctx := context.Background()
	request := &v1.HelloRequest{Name: "World"}

	expectedGreeter := &biz.Greeter{Hello: "World"}
	mockUsecase.On("CreateGreeter", ctx, mock.AnythingOfType("*biz.Greeter")).Return(expectedGreeter, nil)

	response, err := greeterService.SayHello(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Hello World", response.Message)
	mockUsecase.AssertExpectations(t)
}

func TestGreeterService_SayHello_EmptyName(t *testing.T) {
	mockUsecase := new(MockGreeterUsecase)
	greeterService := NewGreeterService(mockUsecase)

	ctx := context.Background()
	request := &v1.HelloRequest{Name: ""}

	expectedGreeter := &biz.Greeter{Hello: ""}
	mockUsecase.On("CreateGreeter", ctx, mock.AnythingOfType("*biz.Greeter")).Return(expectedGreeter, nil)

	response, err := greeterService.SayHello(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Hello ", response.Message)
	mockUsecase.AssertExpectations(t)
}
