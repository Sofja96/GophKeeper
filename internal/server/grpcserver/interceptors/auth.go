package interceptors

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Sofja96/GophKeeper.git/internal/models"
)

// Claims представляет собой структуру для хранения информации о пользователе в JWT.
type Claims struct {
	jwt.RegisteredClaims
	User string
}

const (
	// JwtSecret используется для подписи токенов JWT.
	JwtSecret = "JWT_SECRET"

	// TokenExp указывает время истечения токена (24 часа).
	TokenExp = time.Hour * 24

	// BearerSchema - схема авторизации, ожидаемая в заголовке Authorization.
	BearerSchema = "Bearer "
)

// CreateToken создает новый JWT токен для пользователя с указанным именем.
func CreateToken(user string) (string, error) {
	claims := Claims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		user,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyToken проверяет и расшифровывает JWT токен.
func VerifyToken(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(JwtSecret), nil
		})
	if err != nil {
		return "", fmt.Errorf("error on parsing token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	fmt.Println("Token is valid")
	return claims.User, nil
}

// AuthInterceptor перехватывает gRPC-запросы и проверяет токен
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if strings.HasSuffix(info.FullMethod, "/Login") || strings.HasSuffix(info.FullMethod, "/Register") {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("missing metadata")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, fmt.Errorf("authorization token is required")
		}

		authHeader := authHeaders[0]
		if !strings.HasPrefix(authHeader, BearerSchema) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid authorization format")
		}

		tokenString := strings.TrimPrefix(authHeader, BearerSchema)

		user, err := VerifyToken(tokenString)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "You must be logged in to access this resource")
		}

		ctx = context.WithValue(ctx, models.ContextKeyUser, user)
		return handler(ctx, req)
	}
}
