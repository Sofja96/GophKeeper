package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sofja96/GophKeeper.git/internal/models"
	amock "github.com/Sofja96/GophKeeper.git/internal/server/app/mocks"
	smock "github.com/Sofja96/GophKeeper.git/internal/server/service/mocks"
	"github.com/Sofja96/GophKeeper.git/internal/server/utils"
	"github.com/Sofja96/GophKeeper.git/proto"
)

type mocks struct {
	app     *amock.MockServer
	service *smock.MockService
}

func TestRegister(t *testing.T) {
	type (
		args struct {
			user *models.User
		}
		mockBehavior func(m *mocks, args args)
	)

	tests := []struct {
		name            string
		req             *proto.RegisterRequest
		args            args
		mockBehavior    mockBehavior
		expectedError   error
		expectedMessage string
	}{
		{
			name: "TestRegisterUserSuccess",
			req: &proto.RegisterRequest{
				Username: "testuser",
				Password: "password123",
			},
			args: args{
				user: &models.User{
					Username: "testuser",
					Password: "password123",
				},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service)
				m.service.EXPECT().RegisterUser(gomock.Any(), args.user).Return(&models.User{
					Username: "testuser",
					Password: "password123",
				}, nil)
			},
			expectedError:   nil,
			expectedMessage: "User registered successfully",
		},
		{
			name: "TestRegisterUserAlreadyExists",
			req: &proto.RegisterRequest{
				Username: "testuser",
				Password: "password123",
			},
			args: args{
				user: &models.User{
					Username: "testuser",
					Password: "password123",
				},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service)
				m.service.EXPECT().RegisterUser(gomock.Any(), args.user).
					Return(nil, utils.ErrUserExists)
			},
			expectedError:   status.Errorf(codes.AlreadyExists, "user testuser already exists"),
			expectedMessage: "",
		},
		{
			name: "TestRegisterUserInternalError",
			req: &proto.RegisterRequest{
				Username: "testuser",
				Password: "password123",
			},
			args: args{
				user: &models.User{
					Username: "testuser",
					Password: "password123",
				},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service)
				m.service.EXPECT().RegisterUser(gomock.Any(), args.user).
					Return(nil, errors.New("failed to hash password"))
			},
			expectedError:   status.Errorf(codes.Internal, "failed to register user: failed to hash password"),
			expectedMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := &mocks{
				app:     amock.NewMockServer(ctrl),
				service: smock.NewMockService(ctrl),
			}
			tt.mockBehavior(m, tt.args)

			server := &gophKeeperServer{
				UnimplementedGophKeeperServer: proto.UnimplementedGophKeeperServer{},
				server:                        m.app,
			}

			resp, err := server.Register(context.Background(), tt.req)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMessage, resp.Message)
			}

		})
	}
}

func TestLogin(t *testing.T) {
	type (
		args struct {
			user *models.User
		}
		mockBehavior func(m *mocks, args args)
	)

	tests := []struct {
		name            string
		req             *proto.LoginRequest
		args            args
		mockBehavior    mockBehavior
		expectedError   error
		expectedToken   string
		expectedMessage string
	}{
		{
			name: "TestLoginUserSuccess",
			req: &proto.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			args: args{
				user: &models.User{
					Username: "testuser",
					Password: "password123",
				},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service).Times(2)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").
					Return(int64(1), nil)
				m.service.EXPECT().LoginUser(gomock.Any(), args.user).Return("Bearer mock_token",
					nil)
			},
			expectedError:   nil,
			expectedToken:   "Bearer mock_token",
			expectedMessage: "Login successful",
		},
		{
			name: "TestLoginUserNotFound",
			req: &proto.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			args: args{
				user: &models.User{
					Username: "testuser",
					Password: "password123",
				},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service)
				m.service.EXPECT().LoginUser(gomock.Any(), args.user).
					Return("", fmt.Errorf("users not found, please to registration"))
			},
			expectedError:   status.Errorf(codes.Unauthenticated, "failed to login: users not found, please to registration"),
			expectedToken:   "",
			expectedMessage: "",
		},
		{
			name: "TestLoginInvalidPassword",
			req: &proto.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			args: args{
				user: &models.User{
					Username: "testuser",
					Password: "wrongpassword",
				},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service)
				m.service.EXPECT().LoginUser(gomock.Any(), args.user).
					Return("", fmt.Errorf("invalid password"))
			},
			expectedError:   status.Errorf(codes.Unauthenticated, "failed to login: invalid password"),
			expectedMessage: "",
		},
		{
			name: "TestLoginEmptyCredentials",
			req: &proto.LoginRequest{
				Username: "",
				Password: "",
			},
			args: args{
				user: &models.User{},
			},
			mockBehavior: func(m *mocks, args args) {
			},
			expectedError: status.Errorf(codes.InvalidArgument, "empty credentials"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := &mocks{
				app:     amock.NewMockServer(ctrl),
				service: smock.NewMockService(ctrl),
			}
			tt.mockBehavior(m, tt.args)

			server := &gophKeeperServer{
				UnimplementedGophKeeperServer: proto.UnimplementedGophKeeperServer{},
				server:                        m.app,
			}

			resp, err := server.Login(context.Background(), tt.req)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMessage, resp.Message)
				assert.Equal(t, tt.expectedToken, resp.Token)
			}

		})
	}
}

func TestNewGophKeeperServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApp := amock.NewMockServer(ctrl)

	server := NewGophKeeperServer(mockApp)

	assert.NotNil(t, server)

	gkServer, ok := server.(*gophKeeperServer)
	assert.True(t, ok, "Expected type *gophKeeperServer")

	assert.Equal(t, mockApp, gkServer.server)
}
