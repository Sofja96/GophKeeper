package interceptors

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	mlogger "github.com/Sofja96/GophKeeper.git/internal/server/logger/mocks"
)

func TestLoggingInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mlogger.NewMockILogger(ctrl)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("key", "value"))

	req := "test request"
	resp := "test response"

	mockLogger.EXPECT().Info("gRPC method %s called", gomock.Any()).Times(1)
	mockLogger.EXPECT().Info("Metadata: %v", gomock.Any()).Times(1)
	mockLogger.EXPECT().Info(
		"gRPC method %s completed successfully, duration: %s, status: %s, size: %d",
		gomock.Any(),
		gomock.Any(),
		codes.OK.String(),
		len(`"test response"`),
	).Times(1)

	interceptor := LoggingInterceptor(mockLogger)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return resp, nil
	}

	start := time.Now()
	result, err := interceptor(ctx, req, &grpc.UnaryServerInfo{FullMethod: "/test.Method"}, handler)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, resp, result)

	assert.True(t, duration > 0)
}

func TestLoggingInterceptor_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mlogger.NewMockILogger(ctrl)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("key", "value"))
	req := "test request"

	gomock.InOrder(
		mockLogger.EXPECT().Info("gRPC method %s called", "/test.Method").Times(1),
		mockLogger.EXPECT().Info("Metadata: %v", map[string][]string{"key": {"value"}}).Times(1),
		mockLogger.EXPECT().Info(
			"gRPC method %s failed with code %s: %s, duration: %s",
			"/test.Method",
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Times(1),
	)

	interceptor := LoggingInterceptor(mockLogger)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, status.Error(codes.Internal, "internal error")
	}

	start := time.Now()
	result, err := interceptor(ctx, req, &grpc.UnaryServerInfo{FullMethod: "/test.Method"}, handler)
	duration := time.Since(start)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, codes.Internal, status.Code(err))

	assert.True(t, duration > 0)
}
