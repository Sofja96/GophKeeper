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

type gophKeeperServer struct {
	proto.UnimplementedGophKeeperServer
	server app.Server
}

// NewGophKeeperServer создает новый экземпляр authServiceServer.
func NewGophKeeperServer(srv app.Server) proto.GophKeeperServer {
	return &gophKeeperServer{
		UnimplementedGophKeeperServer: proto.UnimplementedGophKeeperServer{},
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

// Login обрабытвает gRPC запрос для входа пользователя.
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

	userID, err := s.server.GetService().GetUserIDByUsername(ctx, user.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user ID: %v", err)
	}

	return &proto.LoginResponse{
		UserId:  userID,
		Token:   token,
		Message: "Login successful",
	}, nil
}
