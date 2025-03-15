package interceptors

import (
	"context"
	"testing"

	"github.com/Sofja96/GophKeeper.git/internal/models"

	_ "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuthInterceptor(t *testing.T) {
	interceptor := AuthInterceptor()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	t.Run("allows Login and Register endpoints", func(t *testing.T) {
		req := struct{}{}
		ctx := context.Background()
		info := &grpc.UnaryServerInfo{FullMethod: "/UserService/Login"}

		resp, err := interceptor(ctx, req, info, handler)
		require.NoError(t, err)
		assert.Equal(t, "success", resp)

		info.FullMethod = "/UserService/Register"
		resp, err = interceptor(ctx, req, info, handler)
		require.NoError(t, err)
		assert.Equal(t, "success", resp)
	})

	t.Run("returns error if authorization header is missing", func(t *testing.T) {
		req := struct{}{}
		ctx := context.Background()
		info := &grpc.UnaryServerInfo{FullMethod: "/UserService/Protected"}

		_, err := interceptor(ctx, req, info, handler)
		require.Error(t, err)
		assert.Equal(t, "missing metadata", err.Error())
	})

	t.Run("returns error if authorization header format is invalid", func(t *testing.T) {
		req := struct{}{}
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "InvalidToken"))
		info := &grpc.UnaryServerInfo{FullMethod: "/UserService/Protected"}

		_, err := interceptor(ctx, req, info, handler)
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Contains(t, err.Error(), "invalid authorization format")
	})

	t.Run("returns error if token is invalid", func(t *testing.T) {
		req := struct{}{}
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer invalid.token.string"))
		info := &grpc.UnaryServerInfo{FullMethod: "/UserService/Protected"}

		_, err := interceptor(ctx, req, info, handler)
		require.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
		assert.Contains(t, err.Error(), "You must be logged in to access this resource")
	})
	t.Run("returns error if authorization header is missing in metadata", func(t *testing.T) {
		req := struct{}{}

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs())

		info := &grpc.UnaryServerInfo{FullMethod: "/UserService/Protected"}

		_, err := interceptor(ctx, req, info, handler)
		require.Error(t, err)
		assert.Equal(t, "authorization token is required", err.Error())
	})

	t.Run("allows request with valid token", func(t *testing.T) {
		token, err := CreateToken("testuser")
		require.NoError(t, err)

		req := struct{}{}

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			user, ok := ctx.Value(models.ContextKeyUser).(string)
			assert.True(t, ok, "User should be set in context")
			assert.Equal(t, "testuser", user)
			return "success", nil
		}

		info := &grpc.UnaryServerInfo{FullMethod: "/UserService/Protected"}

		resp, err := interceptor(ctx, req, info, handler)
		require.NoError(t, err)
		assert.Equal(t, "success", resp)
	})
}
