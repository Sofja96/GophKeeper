package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Sofja96/GophKeeper.git/internal/models"
	mlogger "github.com/Sofja96/GophKeeper.git/internal/server/logger/mocks"
	mockdb "github.com/Sofja96/GophKeeper.git/internal/server/storage/db/mocks"
	mockminio "github.com/Sofja96/GophKeeper.git/internal/server/storage/minio/mocks"
	"github.com/Sofja96/GophKeeper.git/internal/server/utils"
)

func TestService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockdb.NewMockAdapter(ctrl)
	service := New(mockDB, nil, nil)

	ctx := context.Background()
	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}

	t.Run("successful registration", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).Return(false, nil)
		mockDB.EXPECT().CreateUser(ctx, gomock.Any()).Return(user, nil)

		result, err := service.RegisterUser(ctx, user)
		assert.NoError(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).Return(true, nil)

		_, err := service.RegisterUser(ctx, user)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, utils.ErrUserExists))
	})
	t.Run("error create user", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).Return(false, nil)
		mockDB.EXPECT().CreateUser(ctx, gomock.Any()).
			Return(nil, errors.New("failed to create user"))

		_, err := service.RegisterUser(ctx, user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create user")
	})
	t.Run("error get userID", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).
			Return(false, fmt.Errorf("error checking existing user"))

		_, err := service.RegisterUser(ctx, user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error checking existing user")
	})
}

func TestService_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockdb.NewMockAdapter(ctrl)
	service := New(mockDB, nil, nil)

	ctx := context.Background()
	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}

	t.Run("successful login", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).Return(true, nil)
		mockDB.EXPECT().GetUserHashPassword(ctx, user.Username).
			Return("$2a$10$k8sLGTcrvuI36ZsTddy7EOgarUqltq2nlu5qv2ZG1IiZbqzvYAqjG", nil)

		token, err := service.LoginUser(ctx, user)
		assert.NoError(t, err)
		assert.Contains(t, token, "Bearer ")
	})

	t.Run("user not found", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).Return(false, nil)

		_, err := service.LoginUser(ctx, user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "users not found")
	})
	t.Run("error checking user", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).
			Return(false, fmt.Errorf("error checking existing user"))

		_, err := service.LoginUser(ctx, user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error checking existing user")
	})

	t.Run("error getting password", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).Return(true, nil)
		mockDB.EXPECT().GetUserHashPassword(ctx, user.Username).
			Return("", fmt.Errorf("error getting password on user"))

		_, err := service.LoginUser(ctx, user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error getting password on user")
	})
}

func TestGetUserIDByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockdb.NewMockAdapter(ctrl)
	service := New(mockDB, nil, nil)

	ctx := context.Background()
	user := &models.User{
		Username: "testuser",
	}

	t.Run("successful get user", func(t *testing.T) {
		mockDB.EXPECT().GetUserID(ctx, user.Username).Return(int64(1), nil)

		userID, err := service.GetUserIDByUsername(ctx, user.Username)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), userID)
	})
}

func TestCreateData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockdb.NewMockAdapter(ctrl)
	mockMinio := mockminio.NewMockClient(ctrl)
	mockLogger := mlogger.NewMockILogger(ctrl)

	s := New(mockDB, mockMinio, mockLogger)

	t.Run("successful creation of binary data", func(t *testing.T) {
		data := &models.Data{
			DataType:    models.BinaryData,
			FileName:    "testfile",
			DataContent: []byte("test content"),
		}

		mockMinio.EXPECT().UploadFile(gomock.Any(), data.FileName, data.DataContent).
			Return("file_url", nil)
		mockDB.EXPECT().CreateData(gomock.Any(), gomock.Any()).Return(int64(1), nil)

		id, err := s.CreateData(context.Background(), data)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})

	t.Run("failed to upload binary data to MinIO", func(t *testing.T) {
		data := &models.Data{
			DataType:    models.BinaryData,
			FileName:    "testfile",
			DataContent: []byte("test content"),
		}

		mockMinio.EXPECT().UploadFile(gomock.Any(), data.FileName, data.DataContent).
			Return("", errors.New("upload failed"))
		mockLogger.EXPECT().Error("ошибка загрузки в Minio: %v", gomock.Any()).Times(1)

		_, err := s.CreateData(context.Background(), data)
		assert.Error(t, err)
	})

	t.Run("successful creation of non-binary data", func(t *testing.T) {
		data := &models.Data{
			DataType:    models.TextData,
			DataContent: []byte("test content"),
		}

		mockDB.EXPECT().CreateData(gomock.Any(), gomock.Any()).Return(int64(1), nil)

		id, err := s.CreateData(context.Background(), data)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})
}

func TestGetData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockdb.NewMockAdapter(ctrl)
	mockMinio := mockminio.NewMockClient(ctrl)
	mockLogger := mlogger.NewMockILogger(ctrl)

	s := New(mockDB, mockMinio, mockLogger)

	t.Run("successful retrieval of data", func(t *testing.T) {
		data := []models.Data{
			{
				ID:          1,
				DataType:    models.TextData,
				DataContent: []byte("test content"),
			},
		}

		mockDB.EXPECT().GetData(gomock.Any(), int64(1)).Return(data, nil)

		result, err := s.GetData(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, data, result)
	})

	t.Run("successful retrieval of binary data", func(t *testing.T) {
		data := []models.Data{
			{
				ID:       1,
				DataType: models.BinaryData,
				Metadata: map[string]interface{}{"file_url": "file_url"},
			},
		}

		mockDB.EXPECT().GetData(gomock.Any(), int64(1)).Return(data, nil)
		mockMinio.EXPECT().GetFile(gomock.Any(), "file_url").Return([]byte("test content"), nil)

		result, err := s.GetData(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, []byte("test content"), result[0].DataContent)
	})

	t.Run("failed to retrieve binary data from MinIO", func(t *testing.T) {
		data := []models.Data{
			{
				ID:       1,
				DataType: models.BinaryData,
				Metadata: map[string]interface{}{"file_url": "file_url"},
			},
		}

		mockDB.EXPECT().GetData(gomock.Any(), int64(1)).Return(data, nil)
		mockMinio.EXPECT().GetFile(gomock.Any(), "file_url").
			Return(nil, errors.New("failed to get file"))
		mockLogger.EXPECT().Error("failed to load file from MinIO: %v", gomock.Any()).Times(1)

		result, err := s.GetData(context.Background(), 1)
		assert.NoError(t, err)
		assert.Nil(t, result[0].DataContent)
	})

	t.Run("no data found", func(t *testing.T) {
		mockDB.EXPECT().GetData(gomock.Any(), int64(1)).Return(nil, nil)

		_, err := s.GetData(context.Background(), 1)
		assert.Error(t, err)
		assert.Equal(t, utils.ErrUserDataNotFound, err)
	})

	t.Run("file_url not found in metadata for binary data", func(t *testing.T) {
		data := []models.Data{
			{
				ID:       1,
				DataType: models.BinaryData,
				Metadata: map[string]interface{}{},
			},
		}

		mockDB.EXPECT().GetData(gomock.Any(), int64(1)).Return(data, nil)
		mockLogger.EXPECT().Error("file_url not found in metadata for binary data").Times(1)

		result, err := s.GetData(context.Background(), 1)
		assert.NoError(t, err)
		assert.Nil(t, result[0].DataContent)
	})
	t.Run("file_url is not a valid string", func(t *testing.T) {
		data := []models.Data{
			{
				ID:       1,
				DataType: models.BinaryData,
				Metadata: map[string]interface{}{"file_url": 123},
			},
		}

		mockDB.EXPECT().GetData(gomock.Any(), int64(1)).Return(data, nil)
		mockLogger.EXPECT().Error("file_url is not a valid string").Times(1)

		result, err := s.GetData(context.Background(), 1)
		assert.NoError(t, err)
		assert.Nil(t, result[0].DataContent)
	})
}

func TestDeleteData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockdb.NewMockAdapter(ctrl)
	mockMinio := mockminio.NewMockClient(ctrl)
	mockLogger := mlogger.NewMockILogger(ctrl)

	s := New(mockDB, mockMinio, mockLogger)

	t.Run("successful deletion of binary data", func(t *testing.T) {
		data := &models.Data{
			ID:       1,
			DataType: models.BinaryData,
			Metadata: map[string]interface{}{"file_url": "file_url"},
		}

		mockDB.EXPECT().GetDataByID(gomock.Any(), int64(1)).Return(data, nil)
		mockMinio.EXPECT().DeleteFile(gomock.Any(), "file_url").Return(nil)
		mockDB.EXPECT().DeleteData(gomock.Any(), int64(1), int64(3)).Return(true, nil)

		success, err := s.DeleteData(context.Background(), 1, 3)
		assert.NoError(t, err)
		assert.True(t, success)
	})

	t.Run("failed to delete binary data from MinIO", func(t *testing.T) {
		data := &models.Data{
			ID:       1,
			DataType: models.BinaryData,
			Metadata: map[string]interface{}{"file_url": "file_url"},
		}

		mockDB.EXPECT().GetDataByID(gomock.Any(), int64(1)).Return(data, nil)
		mockMinio.EXPECT().DeleteFile(gomock.Any(), "file_url").
			Return(errors.New("failed to delete file"))

		_, err := s.DeleteData(context.Background(), 1, 3)
		assert.Error(t, err)
	})

	t.Run("successful deletion of non-binary data", func(t *testing.T) {
		data := &models.Data{
			ID:       1,
			DataType: models.TextData,
		}

		mockDB.EXPECT().GetDataByID(gomock.Any(), int64(1)).Return(data, nil)
		mockDB.EXPECT().DeleteData(gomock.Any(), int64(1), int64(3)).Return(true, nil)

		success, err := s.DeleteData(context.Background(), 1, 3)
		assert.NoError(t, err)
		assert.True(t, success)
	})

	t.Run("data not found", func(t *testing.T) {
		mockDB.EXPECT().GetDataByID(gomock.Any(), int64(1)).
			Return(nil, errors.New("data not found"))

		_, err := s.DeleteData(context.Background(), 1, 3)
		assert.Error(t, err)
	})

	t.Run("file_url not found in metadata for binary data", func(t *testing.T) {
		data := &models.Data{
			ID:       1,
			DataType: models.BinaryData,
			Metadata: map[string]interface{}{},
		}

		mockDB.EXPECT().GetDataByID(gomock.Any(), int64(1)).Return(data, nil)

		_, err := s.DeleteData(context.Background(), 1, 3)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file_url не найден в метаданных")
	})

	t.Run("failed to delete data from database", func(t *testing.T) {
		data := &models.Data{
			ID:       1,
			DataType: models.TextData,
		}

		mockDB.EXPECT().GetDataByID(gomock.Any(), int64(1)).Return(data, nil)
		mockDB.EXPECT().DeleteData(gomock.Any(), int64(1), int64(3)).
			Return(false, fmt.Errorf("ошибка удаления данных из базы данных"))

		_, err := s.DeleteData(context.Background(), 1, 3)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка удаления данных из базы данных")
	})

	t.Run("data not found in database", func(t *testing.T) {
		data := &models.Data{
			ID:       1,
			DataType: models.TextData,
		}

		mockDB.EXPECT().GetDataByID(gomock.Any(), int64(1)).Return(data, nil)
		mockDB.EXPECT().DeleteData(gomock.Any(), int64(1), int64(3)).Return(false, nil)

		_, err := s.DeleteData(context.Background(), 1, 3)
		assert.Error(t, err)
		assert.Equal(t, utils.ErrUserDataNotFound, err)
	})
}

func TestUpdateData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockdb.NewMockAdapter(ctrl)
	mockMinio := mockminio.NewMockClient(ctrl)
	mockLogger := mlogger.NewMockILogger(ctrl)

	s := New(mockDB, mockMinio, mockLogger)

	t.Run("successful update of binary data", func(t *testing.T) {
		oldData := &models.Data{
			ID:       1,
			DataType: models.BinaryData,
			Metadata: map[string]interface{}{"file_url": "old_file_url"},
		}

		newData := &models.Data{
			ID:          1,
			DataType:    models.BinaryData,
			FileName:    "new_file",
			DataContent: []byte("new content"),
		}

		mockDB.EXPECT().GetDataByID(gomock.Any(), int64(1)).Return(oldData, nil)
		mockMinio.EXPECT().UpdateFile(gomock.Any(),
			"old_file_url", newData.FileName, newData.DataContent).
			Return("new_file_url", nil)
		mockDB.EXPECT().UpdateData(gomock.Any(), newData).Return(nil)

		err := s.UpdateData(context.Background(), newData)
		assert.NoError(t, err)
	})

	t.Run("failed to update binary data in MinIO", func(t *testing.T) {
		oldData := &models.Data{
			ID:       1,
			DataType: models.BinaryData,
			Metadata: map[string]interface{}{"file_url": "old_file_url"},
		}

		newData := &models.Data{
			ID:          1,
			DataType:    models.BinaryData,
			FileName:    "new_file",
			DataContent: []byte("new content"),
		}

		mockDB.EXPECT().GetDataByID(gomock.Any(), int64(1)).Return(oldData, nil)
		mockMinio.EXPECT().UpdateFile(gomock.Any(),
			"old_file_url", newData.FileName, newData.DataContent).
			Return("", errors.New("failed to update file"))

		err := s.UpdateData(context.Background(), newData)
		assert.Error(t, err)
	})

	t.Run("successful update of non-binary data", func(t *testing.T) {
		oldData := &models.Data{
			ID:       1,
			DataType: models.TextData,
		}

		newData := &models.Data{
			ID:          1,
			DataType:    models.TextData,
			DataContent: []byte("new content"),
		}

		mockDB.EXPECT().GetDataByID(gomock.Any(), int64(1)).Return(oldData, nil)
		mockDB.EXPECT().UpdateData(gomock.Any(), newData).Return(nil)

		err := s.UpdateData(context.Background(), newData)
		assert.NoError(t, err)
	})

	t.Run("data not found", func(t *testing.T) {
		mockDB.EXPECT().GetDataByID(gomock.Any(), int64(1)).
			Return(nil, errors.New("data not found"))

		err := s.UpdateData(context.Background(), &models.Data{ID: 1})
		assert.Error(t, err)
	})
	t.Run("failed to update data from database", func(t *testing.T) {
		oldData := &models.Data{
			ID:       1,
			DataType: models.TextData,
		}

		newData := &models.Data{
			ID:          1,
			DataType:    models.TextData,
			DataContent: []byte("new content"),
		}

		mockDB.EXPECT().GetDataByID(gomock.Any(), int64(1)).Return(oldData, nil)
		mockDB.EXPECT().UpdateData(gomock.Any(), newData).
			Return(errors.New("error update update data"))

		err := s.UpdateData(context.Background(), newData)
		assert.Error(t, err)
	})

	t.Run("file_url not found in metadata for binary data", func(t *testing.T) {
		oldData := &models.Data{
			ID:       1,
			DataType: models.BinaryData,
			Metadata: map[string]interface{}{},
		}

		newData := &models.Data{
			ID:          1,
			DataType:    models.BinaryData,
			FileName:    "new_file",
			DataContent: []byte("new content"),
		}

		mockDB.EXPECT().GetDataByID(gomock.Any(), int64(1)).Return(oldData, nil)

		err := s.UpdateData(context.Background(), newData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file_url не найден в метаданных")
	})
}
