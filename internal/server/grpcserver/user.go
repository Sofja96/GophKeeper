package grpcserver

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sofja96/GophKeeper.git/internal/models"
	"github.com/Sofja96/GophKeeper.git/internal/server/app"
	"github.com/Sofja96/GophKeeper.git/internal/server/utils"
	"github.com/Sofja96/GophKeeper.git/proto"
)

//todo покрыть тестами и вызывать в тесте аpp и сервис

// authServiceServer реализует gRPC методы для AuthService.
type gophKeeperServer struct {
	proto.UnimplementedGophKeeperServer
	server app.Server
}

// NewGophKeeperServer создает новый экземпляр authServiceServer.
func NewGophKeeperServer(srv app.Server) proto.GophKeeperServer {
	return &gophKeeperServer{
		UnimplementedGophKeeperServer: proto.UnimplementedGophKeeperServer{}, // Инициализация
		server:                        srv,
	}
}

// Register обрабатывает gRPC запрос для регистрации пользователя.
func (s *gophKeeperServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	user := &models.User{
		Username: req.Username,
		Password: req.Password,
	}
	_, err := s.server.GetService().RegisterUser(ctx, user)
	if err != nil {
		if errors.Is(err, utils.ErrUserExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user %s already exists", user.Username)
		}
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}
	return &proto.RegisterResponse{Message: "User registered successfully"}, nil
}

func (s *gophKeeperServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	user := &models.User{
		Username: req.Username,
		Password: req.Password,
	}
	if len(user.Username) == 0 && len(user.Password) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty credentials")
	}
	token, err := s.server.GetService().LoginUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to login: %v", err)
	}

	return &proto.LoginResponse{
		Token:   token,
		Message: "Login successful",
	}, nil
}

// todo пока тестовая функция ее не тестировать
func (s *gophKeeperServer) GetUserData(ctx context.Context, req *proto.GetUserDataRequest) (*proto.GetUserDataResponse, error) {
	// Получаем пользователя из контекста
	user, ok := ctx.Value(models.ContextKeyUser).(string) // user это строка, т.е. имя пользователя
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid user data in context")
	}

	// Реализуем логику для получения данных пользователя
	// Например, можно извлечь информацию из базы данных или использовать username для дальнейших действий
	return &proto.GetUserDataResponse{
		Username: user,
		Message:  "User data fetched successfully",
	}, nil
}
