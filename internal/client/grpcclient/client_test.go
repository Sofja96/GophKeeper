package grpcclient

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Sofja96/GophKeeper.git/internal/server/settings"
	"github.com/Sofja96/GophKeeper.git/proto"
	mproto "github.com/Sofja96/GophKeeper.git/proto/mocks"
	"github.com/Sofja96/GophKeeper.git/shared"
)

func TestNewGRPCClient_Success(t *testing.T) {
	certPath := "rootCACert.pem"
	keyPath := "rootCAKey.pem"

	defer os.Remove(certPath)
	defer os.Remove(keyPath)

	err := shared.GenerateCertificate(certPath, keyPath)
	assert.NoError(t, err)

	conf := &settings.Settings{
		PathCert: certPath,
		PathKey:  keyPath,
		Host:     "localhost",
		Port:     "50051",
	}

	client, err := NewGRPCClient(conf)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewGRPCClient_Error(t *testing.T) {
	conf := &settings.Settings{
		PathCert: "path/to/cert",
		PathKey:  "path/to/key",
		Host:     "localhost",
		Port:     "50051",
	}

	client, err := NewGRPCClient(conf)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "failed to read cert file")
}

func TestClient_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)

	mockClient.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(&proto.LoginResponse{Token: "mock-token"}, nil).
		Times(1)

	client := &Client{
		Client: mockClient,
	}

	token, err := client.Login("testuser", "password123")

	assert.NoError(t, err)
	assert.Equal(t, "mock-token", token)
}

func TestClient_Login_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)

	mockClient.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("invalid credentials")).
		Times(1)

	client := &Client{
		Client: mockClient,
	}

	token, err := client.Login("testuser", "password123")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "login failed")
}

func TestClient_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)

	mockClient.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(&proto.RegisterResponse{}, nil).
		Times(1)

	client := &Client{
		Client: mockClient,
	}

	err := client.Register("testuser", "password123")
	assert.NoError(t, err)
}

func TestClient_Register_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)

	mockClient.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("registration failed")).
		Times(1)

	client := &Client{
		Client: mockClient,
	}

	err := client.Register("testuser", "password123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "registration failed")
}
