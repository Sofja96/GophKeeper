package grpcserver

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/Sofja96/GophKeeper.git/internal/models"
	amock "github.com/Sofja96/GophKeeper.git/internal/server/app/mocks"
	smock "github.com/Sofja96/GophKeeper.git/internal/server/service/mocks"
	"github.com/Sofja96/GophKeeper.git/internal/server/utils"
	"github.com/Sofja96/GophKeeper.git/proto"
)

func TestCreateData(t *testing.T) {
	type (
		args struct {
			req *proto.CreateDataRequest
		}
		mockBehavior func(m *mocks, args args)
	)

	tests := []struct {
		name            string
		args            args
		mockBehavior    mockBehavior
		expectedError   error
		expectedMessage string
		expectedDataID  int64
	}{
		{
			name: "TestCreateDataSuccess",
			args: args{
				req: &proto.CreateDataRequest{
					DataContent: []byte("test content"),
					DataType:    proto.DataType_BINARY_DATA,
					FileName:    "testfile",
					Metadata:    &structpb.Struct{},
				},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service).Times(2)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(1), nil)
				m.service.EXPECT().CreateData(gomock.Any(), gomock.Any()).Return(int64(1), nil)
			},
			expectedError:   nil,
			expectedMessage: "Data successfully created",
			expectedDataID:  1,
		},
		{
			name: "TestCreateDataServiceError",
			args: args{
				req: &proto.CreateDataRequest{
					DataContent: []byte("test content"),
					DataType:    proto.DataType_BINARY_DATA,
					FileName:    "testfile",
					Metadata:    &structpb.Struct{},
				},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service).Times(2)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(1), nil)

				m.service.EXPECT().CreateData(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("service error"))
			},
			expectedError:   status.Errorf(codes.Internal, "failed to create data: service error"),
			expectedMessage: "",
			expectedDataID:  0,
		},
		{
			name: "TestCreateDataUnauthenticated",
			args: args{
				req: &proto.CreateDataRequest{
					DataContent: []byte("test content"),
					DataType:    proto.DataType_BINARY_DATA,
					FileName:    "testfile",
					Metadata:    &structpb.Struct{},
				},
			},
			mockBehavior: func(m *mocks, args args) {
			},
			expectedError:   status.Errorf(codes.Unauthenticated, "invalid user authentication"),
			expectedMessage: "",
			expectedDataID:  0,
		},
		{
			name: "TestCreateDataGetUserIDError",
			args: args{
				req: &proto.CreateDataRequest{
					DataContent: []byte("test content"),
					DataType:    proto.DataType_BINARY_DATA,
					FileName:    "testfile",
					Metadata:    &structpb.Struct{},
				},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(0), errors.New("failed to get user ID"))
			},
			expectedError:   status.Errorf(codes.Internal, "failed to get user ID: failed to get user ID"),
			expectedMessage: "",
			expectedDataID:  0,
		},
		{
			name: "TestCreateDataInvalidDataType",
			args: args{
				req: &proto.CreateDataRequest{
					DataContent: []byte("test content"),
					DataType:    proto.DataType_UNKNOWN,
					FileName:    "testfile",
					Metadata:    &structpb.Struct{},
				},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(1), nil)
			},
			expectedError:   status.Errorf(codes.InvalidArgument, "invalid data type :unknown proto.DataType"),
			expectedMessage: "",
			expectedDataID:  0,
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

			var ctx context.Context
			if tt.name != "TestCreateDataUnauthenticated" {
				ctx = context.WithValue(context.Background(), models.ContextKeyUser, "testuser")
			} else {
				ctx = context.Background()
			}

			resp, err := server.CreateData(ctx, tt.args.req)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMessage, resp.Message)
				assert.Equal(t, tt.expectedDataID, resp.DataId)
			}
		})
	}
}

func TestGetAllData(t *testing.T) {
	type (
		args struct {
			req *proto.GetAllDataRequest
		}
		mockBehavior func(m *mocks, args args)
	)

	tests := []struct {
		name            string
		args            args
		mockBehavior    mockBehavior
		expectedError   error
		expectedMessage string
		expectedData    []*proto.DataItem
	}{
		{
			name: "TestGetAllDataSuccess",
			args: args{
				req: &proto.GetAllDataRequest{},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service).Times(2)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(1), nil)
				m.service.EXPECT().GetData(gomock.Any(), int64(1)).Return([]models.Data{
					{
						ID:          1,
						DataType:    models.TextData,
						DataContent: []byte("test content"),
						Metadata:    map[string]interface{}{"key": "value"},
						UpdatedAt:   time.Now(),
					},
				}, nil)
			},
			expectedError:   nil,
			expectedMessage: "",
			expectedData: []*proto.DataItem{
				{
					DataId:      1,
					DataType:    proto.DataType_TEXT_DATA,
					DataContent: []byte("test content"),
					Metadata:    &structpb.Struct{Fields: map[string]*structpb.Value{"key": {Kind: &structpb.Value_StringValue{StringValue: "value"}}}},
					UpdatedAt:   time.Now().Format(time.RFC3339),
				},
			},
		},
		{
			name: "TestGetAllDataUnauthenticated",
			args: args{
				req: &proto.GetAllDataRequest{},
			},
			mockBehavior: func(m *mocks, args args) {
			},
			expectedError:   status.Errorf(codes.Unauthenticated, "invalid user authentication"),
			expectedMessage: "",
			expectedData:    nil,
		},
		{
			name: "TestGetAllDataGetUserIDError",
			args: args{
				req: &proto.GetAllDataRequest{},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(0), errors.New("failed to get user ID"))
			},
			expectedError:   status.Errorf(codes.Internal, "failed to get user ID: failed to get user ID"),
			expectedMessage: "",
			expectedData:    nil,
		},
		{
			name: "TestGetAllDataNotFound",
			args: args{
				req: &proto.GetAllDataRequest{},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service).Times(2)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(1), nil)
				m.service.EXPECT().GetData(gomock.Any(), int64(1)).Return(nil, utils.ErrUserDataNotFound)
			},
			expectedError:   status.Errorf(codes.NotFound, "not found data for user testuser"),
			expectedMessage: "",
			expectedData:    nil,
		},
		{
			name: "TestGetAllDataError",
			args: args{
				req: &proto.GetAllDataRequest{},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service).Times(2)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(1), nil)
				m.service.EXPECT().GetData(gomock.Any(), int64(1)).Return(nil, errors.New("error service"))
			},
			expectedError:   status.Errorf(codes.Internal, "failed to get data: error service"),
			expectedMessage: "",
			expectedData:    nil,
		},
		{
			name: "TestGetAllDataConvertError",
			args: args{
				req: &proto.GetAllDataRequest{},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service).Times(2)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(1), nil)
				m.service.EXPECT().GetData(gomock.Any(), int64(1)).Return([]models.Data{
					{
						ID:          1,
						DataType:    models.DataType(rune(999)),
						DataContent: []byte("test content"),
						Metadata:    map[string]interface{}{"key": "value"},
						UpdatedAt:   time.Now(),
					},
				}, nil)
			},
			expectedError:   status.Errorf(codes.Internal, "failed to convert data type: unknown DataType: ϧ"),
			expectedMessage: "",
			expectedData:    nil,
		},
		{
			name: "TestGetAllDataConvertMetadataError",
			args: args{
				req: &proto.GetAllDataRequest{},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service).Times(2)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(1), nil)
				m.service.EXPECT().GetData(gomock.Any(), int64(1)).Return([]models.Data{
					{
						ID:          1,
						DataType:    models.TextData,
						DataContent: []byte("test content"),
						Metadata:    map[string]interface{}{"key": func() {}},
						UpdatedAt:   time.Now(),
					},
				}, nil)
			},
			expectedError: status.Errorf(codes.Internal,
				"failed to convert metadata: failed to convert JSONB to structpb.Struct: proto: invalid type: func()"),
			expectedMessage: "",
			expectedData:    nil,
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

			var ctx context.Context
			if tt.name != "TestGetAllDataUnauthenticated" {
				ctx = context.WithValue(context.Background(), models.ContextKeyUser, "testuser")
			} else {
				ctx = context.Background()
			}

			resp, err := server.GetAllData(ctx, tt.args.req)
			if tt.expectedError != nil {
				assert.Error(t, err)
				actual := strings.ReplaceAll(err.Error(), "\u00a0", " ")
				assert.Equal(t, tt.expectedError.Error(), actual)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedData, resp.Data)
			}
		})
	}
}

func TestDeleteData(t *testing.T) {
	type (
		args struct {
			req    *proto.DeleteDataRequest
			userId int64
		}
		mockBehavior func(m *mocks, args args)
	)

	tests := []struct {
		name            string
		args            args
		mockBehavior    mockBehavior
		expectedError   error
		expectedMessage string
	}{
		{
			name: "TestDeleteDataSuccess",
			args: args{
				req: &proto.DeleteDataRequest{
					DataId: 1,
				},
				userId: int64(3),
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service).Times(2)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(3), nil)
				m.service.EXPECT().DeleteData(gomock.Any(), args.req.DataId, args.userId).
					Return(true, nil)
			},
			expectedError:   nil,
			expectedMessage: "Данные с ID 1 успешно удалены",
		},
		{
			name: "TestDeleteDataNotFound",
			args: args{
				req: &proto.DeleteDataRequest{
					DataId: 1,
				},
				userId: int64(2),
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service).Times(2)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(2), nil)
				m.service.EXPECT().DeleteData(gomock.Any(), args.req.DataId, args.userId).
					Return(false, utils.ErrUserDataNotFound)
			},
			expectedError:   status.Errorf(codes.NotFound, "данные с ID 1 не найдены"),
			expectedMessage: "",
		},
		{
			name: "TestDeleteDataInternalError",
			args: args{
				req: &proto.DeleteDataRequest{
					DataId: 1,
				},
				userId: int64(1),
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service).Times(2)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(1), nil)
				m.service.EXPECT().DeleteData(gomock.Any(), args.req.DataId, args.userId).
					Return(false, errors.New("internal error"))
			},
			expectedError:   status.Errorf(codes.Internal, "failed delete data with ID 1"),
			expectedMessage: "",
		},
		{
			name: "TestDeleteDataUnauthenticated",
			args: args{
				req: &proto.DeleteDataRequest{
					DataId: 1,
				},
			},
			mockBehavior: func(m *mocks, args args) {
			},
			expectedError:   status.Errorf(codes.Unauthenticated, "invalid user authentication"),
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

			var ctx context.Context
			if tt.name != "TestDeleteDataUnauthenticated" {
				ctx = context.WithValue(context.Background(), models.ContextKeyUser, "testuser")
			} else {
				ctx = context.Background()
			}

			resp, err := server.DeleteData(ctx, tt.args.req)
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

func TestUpdateData(t *testing.T) {
	type (
		args struct {
			req *proto.UpdateDataRequest
		}
		mockBehavior func(m *mocks, args args)
	)

	tests := []struct {
		name            string
		args            args
		mockBehavior    mockBehavior
		expectedError   error
		expectedMessage string
	}{
		{
			name: "TestUpdateDataSuccess",
			args: args{
				req: &proto.UpdateDataRequest{
					DataId:      1,
					DataContent: []byte("updated content"),
					FileName:    "updatedfile",
					Metadata:    &structpb.Struct{},
				},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service).Times(2)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(1), nil)
				m.service.EXPECT().UpdateData(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError:   nil,
			expectedMessage: "Data successfully updated",
		},
		{
			name: "TestUpdateDataUnauthenticated",
			args: args{
				req: &proto.UpdateDataRequest{
					DataId:      1,
					DataContent: []byte("updated content"),
					FileName:    "updatedfile",
					Metadata:    &structpb.Struct{},
				},
			},
			mockBehavior: func(m *mocks, args args) {
			},
			expectedError:   status.Errorf(codes.Unauthenticated, "invalid user authentication"),
			expectedMessage: "",
		},
		{
			name: "TestUpdateDataGetUserIDError",
			args: args{
				req: &proto.UpdateDataRequest{
					DataId:      1,
					DataContent: []byte("updated content"),
					FileName:    "updatedfile",
					Metadata:    &structpb.Struct{},
				},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(0), errors.New("failed to get user ID"))
			},
			expectedError:   status.Errorf(codes.Internal, "failed to get user ID: failed to get user ID"),
			expectedMessage: "",
		},
		{
			name: "TestUpdateDataServiceError",
			args: args{
				req: &proto.UpdateDataRequest{
					DataId:      1,
					DataContent: []byte("updated content"),
					FileName:    "updatedfile",
					Metadata:    &structpb.Struct{},
				},
			},
			mockBehavior: func(m *mocks, args args) {
				m.app.EXPECT().GetService().Return(m.service).Times(2)
				m.service.EXPECT().GetUserIDByUsername(gomock.Any(), "testuser").Return(int64(1), nil)
				m.service.EXPECT().UpdateData(gomock.Any(), gomock.Any()).Return(errors.New("service error"))
			},
			expectedError:   status.Errorf(codes.Internal, "failed to update data: service error"),
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

			var ctx context.Context
			if tt.name != "TestUpdateDataUnauthenticated" {
				ctx = context.WithValue(context.Background(), models.ContextKeyUser, "testuser")
			} else {
				ctx = context.Background()
			}

			resp, err := server.UpdateData(ctx, tt.args.req)
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
