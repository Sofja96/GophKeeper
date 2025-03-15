package grpcserver

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	amock "github.com/Sofja96/GophKeeper.git/internal/server/app/mocks"
	mlogging "github.com/Sofja96/GophKeeper.git/internal/server/logger/mocks"
	"github.com/Sofja96/GophKeeper.git/internal/server/settings"
	"github.com/Sofja96/GophKeeper.git/pkg"
	"github.com/Sofja96/GophKeeper.git/proto"
	mproto "github.com/Sofja96/GophKeeper.git/proto/mocks"
)

func TestGophKeeperServer_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mproto.NewMockGophKeeperServer(ctrl)

	mockServer.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(&proto.LoginResponse{Token: "mock-token", Message: "Login successful"}, nil).
		Times(1)

	resp, err := mockServer.Login(context.Background(), &proto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.Equal(t, "mock-token", resp.Token)
	assert.Equal(t, "Login successful", resp.Message)
}

func TestGophKeeperServer_Login_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mproto.NewMockGophKeeperServer(ctrl)

	mockServer.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(nil, status.Errorf(codes.Unauthenticated, "invalid credentials")).
		Times(1)

	resp, err := mockServer.Login(context.Background(), &proto.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestGophKeeperServer_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mproto.NewMockGophKeeperServer(ctrl)

	mockServer.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(&proto.RegisterResponse{Message: "User registered successfully"}, nil).
		Times(1)

	resp, err := mockServer.Register(context.Background(), &proto.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", resp.Message)
}

func TestGophKeeperServer_Register_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mproto.NewMockGophKeeperServer(ctrl)

	mockServer.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(nil, status.Errorf(codes.AlreadyExists, "user already exists")).
		Times(1)

	resp, err := mockServer.Register(context.Background(), &proto.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))
	assert.Contains(t, err.Error(), "user already exists")
}

func TestGophKeeperServer_CreateData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mproto.NewMockGophKeeperServer(ctrl)

	mockServer.EXPECT().
		CreateData(gomock.Any(), gomock.Any()).
		Return(&proto.CreateDataResponse{
			Message: "Data successfully created",
			DataId:  int64(2)}, nil).
		Times(1)

	resp, err := mockServer.CreateData(context.Background(), &proto.CreateDataRequest{
		DataContent: []byte("test content"),
		DataType:    3,
		Metadata:    nil,
		FileName:    "testfile",
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(2), resp.DataId)
	assert.Equal(t, "Data successfully created", resp.Message)
}

func TestGophKeeperServer_CreateData_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mproto.NewMockGophKeeperServer(ctrl)

	mockServer.EXPECT().
		CreateData(gomock.Any(), gomock.Any()).
		Return(nil, status.Errorf(codes.Internal, "failed to create data")).
		Times(1)

	resp, err := mockServer.CreateData(context.Background(), &proto.CreateDataRequest{
		DataContent: []byte("test content"),
		DataType:    3,
		Metadata:    nil,
		FileName:    "testfile",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Equal(t, err.Error(), "rpc error: code = Internal desc = failed to create data")
}

func TestGophKeeperServer_GetAllData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mproto.NewMockGophKeeperServer(ctrl)

	mockServer.EXPECT().
		GetAllData(gomock.Any(), gomock.Any()).
		Return(&proto.GetAllDataResponse{
			Data: []*proto.DataItem{
				{
					DataId:      1,
					DataType:    proto.DataType_TEXT_DATA,
					DataContent: []byte("test content"),
					Metadata:    &structpb.Struct{},
					UpdatedAt:   "2023-10-01T12:00:00Z",
				},
			},
		}, nil).
		Times(1)

	resp, err := mockServer.GetAllData(context.Background(), &proto.GetAllDataRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, int64(1), resp.Data[0].DataId)
	assert.Equal(t, proto.DataType_TEXT_DATA, resp.Data[0].DataType)
	assert.Equal(t, []byte("test content"), resp.Data[0].DataContent)
}

func TestGophKeeperServer_GetAllData_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mproto.NewMockGophKeeperServer(ctrl)

	mockServer.EXPECT().
		GetAllData(gomock.Any(), gomock.Any()).
		Return(nil, status.Errorf(codes.NotFound, "not found data for user testuser")).
		Times(1)

	resp, err := mockServer.GetAllData(context.Background(), &proto.GetAllDataRequest{})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.NotFound, status.Code(err))
	assert.Equal(t, "rpc error: code = NotFound desc = not found data for user testuser", err.Error())
}

func TestGophKeeperServer_DeleteData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mproto.NewMockGophKeeperServer(ctrl)

	mockServer.EXPECT().
		DeleteData(gomock.Any(), gomock.Any()).
		Return(&proto.DeleteDataResponse{
			Message: "Данные с ID 1 успешно удалены",
		}, nil).
		Times(1)

	resp, err := mockServer.DeleteData(context.Background(), &proto.DeleteDataRequest{
		DataId: 1,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Данные с ID 1 успешно удалены", resp.Message)
}

func TestGophKeeperServer_DeleteData_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mproto.NewMockGophKeeperServer(ctrl)

	mockServer.EXPECT().
		DeleteData(gomock.Any(), gomock.Any()).
		Return(nil, status.Errorf(codes.Internal, "failed delete data with ID 1")).
		Times(1)

	resp, err := mockServer.DeleteData(context.Background(), &proto.DeleteDataRequest{
		DataId: 1,
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Equal(t, "rpc error: code = Internal desc = failed delete data with ID 1", err.Error())
}

func TestGophKeeperServer_UpdateData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mproto.NewMockGophKeeperServer(ctrl)

	mockServer.EXPECT().
		UpdateData(gomock.Any(), gomock.Any()).
		Return(&proto.UpdateDataResponse{
			Message: "Data successfully updated",
		}, nil).
		Times(1)

	resp, err := mockServer.UpdateData(context.Background(), &proto.UpdateDataRequest{
		DataId:      1,
		DataContent: []byte("updated content"),
		FileName:    "updatedfile",
		Metadata:    &structpb.Struct{},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Data successfully updated", resp.Message)
}

func TestGophKeeperServer_UpdateData_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mproto.NewMockGophKeeperServer(ctrl)

	mockServer.EXPECT().
		UpdateData(gomock.Any(), gomock.Any()).
		Return(nil, status.Errorf(codes.Internal, "failed to update data")).
		Times(1)

	resp, err := mockServer.UpdateData(context.Background(), &proto.UpdateDataRequest{
		DataId:      1,
		DataContent: []byte("updated content"),
		FileName:    "updatedfile",
		Metadata:    &structpb.Struct{},
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Equal(t, "rpc error: code = Internal desc = failed to update data", err.Error())
}

func TestNewGRPCServer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := &mocks{
		app: amock.NewMockServer(ctrl),
	}

	certPath := "rootCACert.pem"
	keyPath := "rootCAKey.pem"

	defer os.Remove(certPath)
	defer os.Remove(keyPath)

	err := pkg.GenerateCertificate(certPath, keyPath)
	assert.NoError(t, err)

	m.app.EXPECT().GetSettings().Return(settings.Settings{
		Host:     "localhost",
		Port:     "50052",
		PathCert: certPath,
		PathKey:  keyPath,
	}).Times(4)

	m.app.EXPECT().GetLogger().Times(2)

	server, err := NewGRPCServer(m.app)
	assert.NoError(t, err)
	assert.NotNil(t, server)
}

func TestNewGRPCServer_Error_Cred(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := &mocks{
		app: amock.NewMockServer(ctrl),
	}

	m.app.EXPECT().GetSettings().Times(4)

	server, err := NewGRPCServer(m.app)
	assert.Error(t, err)
	assert.Nil(t, server)
	assert.Contains(t, err.Error(), "failed to create credentials")
}

func TestRun_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApp := amock.NewMockServer(ctrl)
	mockLogger := mlogging.NewMockILogger(ctrl)

	certPath := "rootCACert.pem"
	keyPath := "rootCAKey.pem"
	defer os.Remove(certPath)
	defer os.Remove(keyPath)

	err := pkg.GenerateCertificate(certPath, keyPath)
	assert.NoError(t, err)

	mockApp.EXPECT().GetSettings().Return(settings.Settings{
		Host:     "localhost",
		Port:     "50053",
		PathCert: certPath,
		PathKey:  keyPath,
	}).Times(4)

	mockApp.EXPECT().GetLogger().Return(mockLogger).Times(3)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).Times(1)

	go func() {
		err = Run(context.Background(), mockApp)
		assert.NoError(t, err)
	}()

	time.Sleep(100 * time.Millisecond)

}
