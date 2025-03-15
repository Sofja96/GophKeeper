//nolint:errcheck
package grpcclient

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Sofja96/GophKeeper.git/internal/client/encryption"
	"github.com/Sofja96/GophKeeper.git/internal/client/localstorage"
	"github.com/Sofja96/GophKeeper.git/internal/client/models"
	mdata "github.com/Sofja96/GophKeeper.git/internal/models"
	"github.com/Sofja96/GophKeeper.git/proto"
	mproto "github.com/Sofja96/GophKeeper.git/proto/mocks"
)

func TestCreateData(t *testing.T) {
	t.Cleanup(func() {
		os.RemoveAll("user_data")
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)

	masterKey := make([]byte, 32)
	copy(masterKey, "16-byte-master-key")

	grpcClient := &Client{
		Client:        mockClient,
		UserID:        12345,
		EncryptionKey: masterKey,
	}

	t.Run("Успешное создание текстовых данных", func(t *testing.T) {
		testData := &models.TextDataType{
			Text: "Sample text data",
		}

		reqData := models.CreateData{
			Data:          testData,
			DataType:      proto.DataType_TEXT_DATA,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		id, err := grpcClient.CreateData(reqData)

		assert.NoError(t, err)
		assert.True(t, id > 0)

		data, err := localstorage.GetAllData(grpcClient.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, data[id])

		storedData := data[id]
		assert.Equal(t, "TEXT_DATA", storedData.DataType.String())
		assert.NotEmpty(t, storedData.DataContent)
		assert.Equal(t, int64(12345), storedData.UserID)
		assert.False(t, storedData.UpdatedAt.IsZero())
	})
	t.Run("Успешное создание бинарных данных", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "testfile.txt")
		if err != nil {
			t.Fatalf("Не удалось создать временный файл: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		testData := []byte("test data")
		if _, err := tmpFile.Write(testData); err != nil {
			t.Fatalf("Не удалось записать данные во временный файл: %v", err)
		}
		tmpFile.Close()

		binaryData := &models.BinaryDataType{
			FilePath: tmpFile.Name(),
			Content:  testData,
			Filename: "testfile.txt",
		}

		reqData := models.CreateData{
			Data:          binaryData,
			DataType:      proto.DataType_BINARY_DATA,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		id, err := grpcClient.CreateData(reqData)

		assert.NoError(t, err)
		assert.True(t, id > 0)

		data, err := localstorage.GetAllData(grpcClient.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, data[id])

		storedData := data[id]
		assert.Equal(t, "BINARY_DATA", storedData.DataType.String())
		assert.NotEmpty(t, storedData.DataContent)
		assert.Equal(t, int64(12345), storedData.UserID)
		assert.False(t, storedData.UpdatedAt.IsZero())
	})
	t.Run("Ошибка валидации текстовых данных", func(t *testing.T) {
		testData := &models.TextDataType{
			Text: "",
		}

		reqData := models.CreateData{
			Data:          testData,
			DataType:      proto.DataType_TEXT_DATA,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		_, err := grpcClient.CreateData(reqData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "поле не может быть пустым")
	})
	t.Run("Ошибка валидации бинарных данных (неверный путь к файлу)", func(t *testing.T) {
		binaryData := &models.BinaryDataType{
			Content:  []byte("test data"),
			Filename: "",
		}

		reqData := models.CreateData{
			Data:          binaryData,
			DataType:      proto.DataType_BINARY_DATA,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		_, err := grpcClient.CreateData(reqData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "поле не может быть пустым")
	})
	t.Run("Ошибка преобразования данных в JSON", func(t *testing.T) {
		invalidData := &InvalidDataType{}

		reqData := models.CreateData{
			Data:          invalidData,
			DataType:      proto.DataType_TEXT_DATA,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		_, err := grpcClient.CreateData(reqData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка преобразования в JSON")
	})
	t.Run("Ошибка при неверном типе данных для бинарного контента", func(t *testing.T) {
		invalidData := &models.TextDataType{
			Text: "Sample text data",
		}

		reqData := models.CreateData{
			Data:          invalidData,
			DataType:      proto.DataType_BINARY_DATA,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		_, err := grpcClient.CreateData(reqData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "неверный тип данных для бинарного контента")
	})
	t.Run("Ошибка шифрования данных", func(t *testing.T) {
		invalidKey := make([]byte, 15)

		testData := &models.TextDataType{
			Text: "Sample text data",
		}

		reqData := models.CreateData{
			Data:          testData,
			DataType:      proto.DataType_TEXT_DATA,
			EncryptionKey: invalidKey,
			Metadata:      nil,
		}

		_, err := grpcClient.CreateData(reqData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка шифрования")
	})
	t.Run("Ошибка конвертации типа данных", func(t *testing.T) {
		invalidDataType := proto.DataType(999)

		testData := &models.TextDataType{
			Text: "Sample text data",
		}

		reqData := models.CreateData{
			Data:          testData,
			DataType:      invalidDataType,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		_, err := grpcClient.CreateData(reqData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка конвертации типа данных")
	})
	t.Run("Ошибка сохранения данных в локальное хранилище", func(t *testing.T) {
		userDir := filepath.Join("user_data", fmt.Sprintf("%d", grpcClient.UserID))
		dataFilePath := filepath.Join(userDir, "data.json")
		err := os.Chmod(dataFilePath, 0444)
		if err != nil {
			t.Fatalf("Не удалось установить права на файл: %v", err)
		}

		testData := &models.TextDataType{
			Text: "Sample text data",
		}

		reqData := models.CreateData{
			Data:          testData,
			DataType:      proto.DataType_TEXT_DATA,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		_, err = grpcClient.CreateData(reqData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка сохранения данных в локальное хранилище")
	})

}

func TestUpdateData(t *testing.T) {
	t.Cleanup(func() {
		os.RemoveAll("user_data")
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)

	masterKey := make([]byte, 32)
	copy(masterKey, "16-byte-master-key")

	grpcClient := &Client{
		Client:        mockClient,
		UserID:        12345,
		EncryptionKey: masterKey,
	}

	t.Run("Успешное обновление текстовых данных", func(t *testing.T) {
		testData := &models.TextDataType{
			Text: "Updated text data",
		}

		reqData := models.CreateData{
			Data:          testData,
			DataType:      proto.DataType_TEXT_DATA,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		dataId := int64(1)

		err := grpcClient.UpdateData(reqData, dataId)

		assert.NoError(t, err)

		data, err := localstorage.GetAllData(grpcClient.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, data[dataId])

		storedData := data[dataId]
		assert.Equal(t, "TEXT_DATA", storedData.DataType.String())
		assert.NotEmpty(t, storedData.DataContent)
		assert.Equal(t, int64(12345), storedData.UserID)
		assert.False(t, storedData.UpdatedAt.IsZero())
	})
	t.Run("Ошибка валидации текстовых данных", func(t *testing.T) {
		testData := &models.TextDataType{
			Text: "",
		}

		reqData := models.CreateData{
			Data:          testData,
			DataType:      proto.DataType_TEXT_DATA,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		dataId := int64(1)

		err := grpcClient.UpdateData(reqData, dataId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка валидации данных")
	})
	t.Run("Ошибка преобразования данных в JSON", func(t *testing.T) {
		invalidData := &InvalidDataType{}

		reqData := models.CreateData{
			Data:          invalidData,
			DataType:      proto.DataType_TEXT_DATA,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		dataId := int64(1)

		err := grpcClient.UpdateData(reqData, dataId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка преобразования в JSON")
	})
	t.Run("Ошибка шифрования данных", func(t *testing.T) {
		invalidKey := make([]byte, 15)

		testData := &models.TextDataType{
			Text: "Sample text data",
		}

		reqData := models.CreateData{
			Data:          testData,
			DataType:      proto.DataType_TEXT_DATA,
			EncryptionKey: invalidKey,
			Metadata:      nil,
		}

		dataId := int64(1)

		err := grpcClient.UpdateData(reqData, dataId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка шифрования")
	})
	t.Run("Ошибка конвертации типа данных", func(t *testing.T) {
		invalidDataType := proto.DataType(999)

		testData := &models.TextDataType{
			Text: "Sample text data",
		}

		reqData := models.CreateData{
			Data:          testData,
			DataType:      invalidDataType,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		dataId := int64(1)

		err := grpcClient.UpdateData(reqData, dataId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка конвертации типа данных")
	})
	t.Run("Успешное обновление бинарных данных", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "testfile.txt")
		if err != nil {
			t.Fatalf("Не удалось создать временный файл: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		testData := []byte("test data")
		if _, err := tmpFile.Write(testData); err != nil {
			t.Fatalf("Не удалось записать данные во временный файл: %v", err)
		}
		tmpFile.Close()

		binaryData := &models.BinaryDataType{
			Content:  testData,
			Filename: "testfile.txt",
			FilePath: tmpFile.Name(),
		}

		reqData := models.CreateData{
			Data:          binaryData,
			DataType:      proto.DataType_BINARY_DATA,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		dataId := int64(1)

		err = grpcClient.UpdateData(reqData, dataId)

		assert.NoError(t, err)

		data, err := localstorage.GetAllData(grpcClient.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, data[dataId])

		storedData := data[dataId]
		assert.Equal(t, "BINARY_DATA", storedData.DataType.String())
		assert.NotEmpty(t, storedData.DataContent)
		assert.Equal(t, int64(12345), storedData.UserID)
		assert.False(t, storedData.UpdatedAt.IsZero())
	})
	t.Run("Ошибка при неверном типе данных для бинарного контента", func(t *testing.T) {
		invalidData := &models.TextDataType{
			Text: "Sample text data",
		}

		reqData := models.CreateData{
			Data:          invalidData,
			DataType:      proto.DataType_BINARY_DATA,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		dataId := int64(1)

		err := grpcClient.UpdateData(reqData, dataId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "неверный тип данных для бинарного контента")
	})
	t.Run("Ошибка сохранения данных в локальное хранилище", func(t *testing.T) {
		userDir := filepath.Join("user_data", fmt.Sprintf("%d", grpcClient.UserID))
		dataFilePath := filepath.Join(userDir, "data.json")
		err := os.Chmod(dataFilePath, 0444)
		if err != nil {
			t.Fatalf("Не удалось установить права на файл: %v", err)
		}
		testData := &models.TextDataType{
			Text: "Sample text data",
		}

		reqData := models.CreateData{
			Data:          testData,
			DataType:      proto.DataType_TEXT_DATA,
			EncryptionKey: masterKey,
			Metadata:      nil,
		}

		dataId := int64(1)

		err = grpcClient.UpdateData(reqData, dataId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка обновления данных в локальном хранилище")
	})
}

func TestDeleteData(t *testing.T) {
	t.Cleanup(func() {
		os.RemoveAll("user_data")
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)

	masterKey := make([]byte, 32)
	copy(masterKey, "16-byte-master-key")

	grpcClient := &Client{
		Client:        mockClient,
		UserID:        12345,
		EncryptionKey: masterKey,
	}

	createTestData := func(dataId int64) {
		encryptedData, err := encryption.EncryptData([]byte(`{"username":"unutest","password":"test"}`), masterKey)
		if err != nil {
			t.Fatalf("Не удалось зашифровать тестовые данные: %v", err)
		}

		data := mdata.Data{
			UserID:      grpcClient.UserID,
			ID:          dataId,
			DataType:    mdata.TextData,
			DataContent: []byte(encryptedData),
			Metadata:    nil,
			UpdatedAt:   time.Now(),
		}

		err = localstorage.SaveData(grpcClient.UserID, data)
		if err != nil {
			t.Fatalf("Не удалось создать тестовые данные: %v", err)
		}
	}

	t.Run("Успешное удаление данных", func(t *testing.T) {
		dataId := int64(1)
		createTestData(dataId)

		mockClient.EXPECT().DeleteData(gomock.Any(), &proto.DeleteDataRequest{DataId: dataId}).
			Return(&proto.DeleteDataResponse{Message: "Данные с ID 1 успешно удалены"}, nil)

		err := grpcClient.DeleteData(dataId)
		assert.NoError(t, err)

		data, err := localstorage.GetAllData(grpcClient.UserID)
		assert.NoError(t, err)
		assert.NotContains(t, data, dataId)
	})
	t.Run("Ошибка удаления данных из локального хранилища", func(t *testing.T) {
		dataId := int64(1)
		createTestData(dataId)

		userDir := filepath.Join("user_data", fmt.Sprintf("%d", grpcClient.UserID))
		dataFilePath := filepath.Join(userDir, "data.json")
		err := os.Chmod(dataFilePath, 0444)
		if err != nil {
			t.Fatalf("Не удалось установить права на файл: %v", err)
		}

		err = grpcClient.DeleteData(dataId)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка удаления данных из локального хранилища")

		err = os.Chmod(dataFilePath, 0644)
		if err != nil {
			t.Fatalf("Не удалось восстановить права на файл: %v", err)
		}
	})
	t.Run("Ошибка удаления данных через gRPC", func(t *testing.T) {
		dataId := int64(1)
		createTestData(dataId)

		mockClient.EXPECT().DeleteData(gomock.Any(), &proto.DeleteDataRequest{DataId: dataId}).
			Return(nil, fmt.Errorf("ошибка gRPC"))

		err := grpcClient.DeleteData(dataId)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка удаления данных")
	})
}

func TestGetData(t *testing.T) {
	t.Cleanup(func() {
		os.RemoveAll("user_data")
	})

	masterKey := make([]byte, 32)
	copy(masterKey, "16-byte-master-key")

	grpcClient := &Client{
		UserID:        12345,
		EncryptionKey: masterKey,
	}

	createTestData := func(dataId int64) {
		encryptedData, err := encryption.EncryptData([]byte(`{"username":"unutest","password":"test"}`), masterKey)
		if err != nil {
			t.Fatalf("Не удалось зашифровать тестовые данные: %v", err)
		}

		data := mdata.Data{
			UserID:      grpcClient.UserID,
			ID:          dataId,
			DataType:    mdata.TextData,
			DataContent: []byte(encryptedData),
			Metadata:    nil,
			UpdatedAt:   time.Now(),
		}

		err = localstorage.SaveData(grpcClient.UserID, data)
		if err != nil {
			t.Fatalf("Не удалось создать тестовые данные: %v", err)
		}
	}

	t.Run("Успешное получение данных", func(t *testing.T) {
		dataId := int64(1)
		createTestData(dataId)

		data, err := grpcClient.GetData()

		assert.NoError(t, err)

		assert.NotEmpty(t, data)
		assert.Equal(t, dataId, data[0].ID)
		assert.Equal(t, mdata.TextData, data[0].DataType)

		expectedData := []byte(`{"username":"unutest","password":"test"}`)
		assert.Equal(t, expectedData, data[0].DataContent)
	})
	t.Run("Успешное получение бинарных данных", func(t *testing.T) {
		dataId := int64(2)
		binaryContent := []byte("binary data")

		encodedData := encryption.EncodeData(binaryContent)

		data := mdata.Data{
			UserID:      grpcClient.UserID,
			ID:          dataId,
			DataType:    mdata.BinaryData,
			DataContent: []byte(encodedData),
			Metadata:    nil,
			UpdatedAt:   time.Now(),
		}

		err := localstorage.SaveData(grpcClient.UserID, data)
		if err != nil {
			t.Fatalf("Не удалось создать тестовые данные: %v", err)
		}

		retrievedData, err := grpcClient.GetData()

		assert.NoError(t, err)
		assert.NotEmpty(t, retrievedData)

		var foundData *mdata.Data
		for _, item := range retrievedData {
			if item.ID == dataId {
				foundData = &item
				break
			}
		}
		assert.NotNil(t, foundData, "Данные с dataId должны быть найдены")
		assert.Equal(t, mdata.BinaryData, foundData.DataType)
		assert.Equal(t, binaryContent, foundData.DataContent)
	})
	t.Run("Ошибка декодирования бинарных данных", func(t *testing.T) {
		dataId := int64(3)
		invalidBinaryContent := []byte("invalid binary data")

		data := mdata.Data{
			UserID:      grpcClient.UserID,
			ID:          dataId,
			DataType:    mdata.BinaryData,
			DataContent: invalidBinaryContent,
			Metadata:    nil,
			UpdatedAt:   time.Now(),
		}

		err := localstorage.SaveData(grpcClient.UserID, data)
		if err != nil {
			t.Fatalf("Не удалось создать тестовые данные: %v", err)
		}

		retrievedData, err := grpcClient.GetData()

		assert.NoError(t, err)

		assert.NotEmpty(t, retrievedData)

		var foundData *mdata.Data
		for _, item := range retrievedData {
			if item.ID == dataId {
				foundData = &item
				break
			}
		}
		assert.NotNil(t, foundData, "Данные с dataId должны быть найдены")

		assert.Equal(t, mdata.BinaryData, foundData.DataType)

		assert.Equal(t, invalidBinaryContent, foundData.DataContent)
	})
	t.Run("Ошибка расшифровки данных", func(t *testing.T) {
		dataId := int64(4)
		invalidEncryptedData := []byte("invalid encrypted data")

		data := mdata.Data{
			UserID:      grpcClient.UserID,
			ID:          dataId,
			DataType:    mdata.TextData,
			DataContent: invalidEncryptedData,
			Metadata:    nil,
			UpdatedAt:   time.Now(),
		}

		err := localstorage.SaveData(grpcClient.UserID, data)
		if err != nil {
			t.Fatalf("Не удалось создать тестовые данные: %v", err)
		}
		_, err = grpcClient.GetData()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка расшифровки данных")
	})
	t.Run("Получение данных из пустого хранилища", func(t *testing.T) {
		err := os.RemoveAll("user_data")
		assert.NoError(t, err)
		data, err := grpcClient.GetData()
		assert.NoError(t, err)
		assert.Empty(t, data)
	})
	t.Run("Ошибка получения данных из локального хранилища", func(t *testing.T) {
		userDir := filepath.Join("user_data", fmt.Sprintf("%d", grpcClient.UserID))
		err := os.MkdirAll(userDir, 0700)
		if err != nil {
			t.Fatalf("Не удалось создать директорию пользователя: %v", err)
		}

		dataFilePath := filepath.Join(userDir, "data.json")
		err = os.WriteFile(dataFilePath, []byte(""), 0644)
		assert.NoError(t, err)
		_, err = grpcClient.GetData()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка получения данных из локального хранилища")
	})
}

func TestSyncData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)

	masterKey := make([]byte, 32)
	copy(masterKey, "16-byte-master-key")

	grpcClient := &Client{
		Client:        mockClient,
		UserID:        12345,
		EncryptionKey: masterKey,
	}

	createTestData := func(dataId int64, dataType mdata.DataType, content []byte, updatedAt time.Time) {
		data := mdata.Data{
			UserID:      grpcClient.UserID,
			ID:          dataId,
			DataType:    dataType,
			DataContent: content,
			Metadata:    nil,
			UpdatedAt:   updatedAt,
		}

		err := localstorage.SaveData(grpcClient.UserID, data)
		if err != nil {
			t.Fatalf("Не удалось создать тестовые данные: %v", err)
		}
	}

	t.Run("Успешная синхронизация данных", func(t *testing.T) {
		t.Cleanup(func() {
			os.RemoveAll("user_data")
		})

		os.RemoveAll("user_data")

		dataId := int64(1)
		updatedAt := time.Now()
		createTestData(dataId, mdata.TextData, []byte("test data"), updatedAt)

		mockClient.EXPECT().GetAllData(gomock.Any(), &proto.GetAllDataRequest{}).Return(&proto.GetAllDataResponse{
			Data: []*proto.DataItem{
				{
					DataId:      dataId,
					DataType:    proto.DataType_TEXT_DATA,
					DataContent: []byte("test data"),
					Metadata:    nil,
					UpdatedAt:   updatedAt.Format(time.RFC3339),
				},
			},
		}, nil)

		mockClient.EXPECT().UpdateData(gomock.Any(), gomock.Any()).
			Return(&proto.UpdateDataResponse{}, nil)

		err := grpcClient.SyncData()
		assert.NoError(t, err)
	})

	t.Run("Обновление локальных данных, если серверные данные новее", func(t *testing.T) {
		t.Cleanup(func() {
			os.RemoveAll("user_data")
		})

		os.RemoveAll("user_data")

		dataId := int64(1)
		localUpdatedAt := time.Now().Add(-1 * time.Hour) // Локальные данные старше
		serverUpdatedAt := time.Now()

		createTestData(dataId, mdata.TextData, []byte("old data"), localUpdatedAt)

		mockClient.EXPECT().GetAllData(gomock.Any(), &proto.GetAllDataRequest{}).
			Return(&proto.GetAllDataResponse{
				Data: []*proto.DataItem{
					{
						DataId:      dataId,
						DataType:    proto.DataType_TEXT_DATA,
						DataContent: []byte("new data"),
						Metadata:    nil,
						UpdatedAt:   serverUpdatedAt.Format(time.RFC3339),
					},
				},
			}, nil)

		err := grpcClient.SyncData()
		assert.NoError(t, err)

		data, err := localstorage.GetAllData(grpcClient.UserID)
		assert.NoError(t, err)
		assert.Equal(t, []byte("new data"), data[dataId].DataContent)
	})

	t.Run("Синхронизация данных с сервера в локальное хранилище", func(t *testing.T) {
		t.Cleanup(func() {
			os.RemoveAll("user_data")
		})

		os.RemoveAll("user_data")

		dataId := int64(1)
		updatedAt := time.Now()

		mockClient.EXPECT().GetAllData(gomock.Any(), &proto.GetAllDataRequest{}).
			Return(&proto.GetAllDataResponse{
				Data: []*proto.DataItem{
					{
						DataId:      dataId,
						DataType:    proto.DataType_TEXT_DATA,
						DataContent: []byte("test data"),
						Metadata:    nil,
						UpdatedAt:   updatedAt.Format(time.RFC3339),
					},
				},
			}, nil)

		err := grpcClient.SyncData()
		assert.NoError(t, err)

		data, err := localstorage.GetAllData(grpcClient.UserID)
		assert.NoError(t, err)
		assert.Contains(t, data, dataId)
		assert.Equal(t, []byte("test data"), data[dataId].DataContent)
	})

	t.Run("Обновление локального ID после отправки данных на сервер", func(t *testing.T) {
		t.Cleanup(func() {
			os.RemoveAll("user_data")
		})

		os.RemoveAll("user_data")

		localID := int64(1)
		newID := int64(2) // Новый ID, который вернет сервер
		updatedAt := time.Now()

		createTestData(localID, mdata.TextData, []byte("test data"), updatedAt)

		mockClient.EXPECT().GetAllData(gomock.Any(), &proto.GetAllDataRequest{}).
			Return(&proto.GetAllDataResponse{
				Data: []*proto.DataItem{}, // Пустой ответ от сервера
			}, nil)

		mockClient.EXPECT().CreateData(gomock.Any(), gomock.Any()).
			Return(&proto.CreateDataResponse{
				DataId: newID,
			}, nil)

		err := grpcClient.SyncData()
		assert.NoError(t, err)

		data, err := localstorage.GetAllData(grpcClient.UserID)
		assert.NoError(t, err)
		assert.NotContains(t, data, localID)
		assert.Contains(t, data, newID)
		assert.Equal(t, []byte("test data"), data[newID].DataContent)
	})
}

func TestSyncData_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)

	masterKey := make([]byte, 32)
	copy(masterKey, "16-byte-master-key")

	grpcClient := &Client{
		Client:        mockClient,
		UserID:        12345,
		EncryptionKey: masterKey,
	}

	createTestData := func(dataId int64, dataType mdata.DataType, content []byte, updatedAt time.Time) {
		data := mdata.Data{
			UserID:      grpcClient.UserID,
			ID:          dataId,
			DataType:    dataType,
			DataContent: content,
			Metadata:    nil,
			UpdatedAt:   updatedAt,
		}

		err := localstorage.SaveData(grpcClient.UserID, data)
		if err != nil {
			t.Fatalf("Не удалось создать тестовые данные: %v", err)
		}
	}

	t.Run("Ошибка получения данных с сервера", func(t *testing.T) {
		t.Cleanup(func() {
			os.RemoveAll("user_data")
		})

		mockClient.EXPECT().GetAllData(gomock.Any(), &proto.GetAllDataRequest{}).
			Return(nil, fmt.Errorf("ошибка получения данных с сервера"))

		err := grpcClient.SyncData()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка получения данных с сервера")
	})
	t.Run("Ошибка получения данных из локального хранилища", func(t *testing.T) {
		t.Cleanup(func() {
			os.RemoveAll("user_data")
		})

		userDir := filepath.Join("user_data", fmt.Sprintf("%d", grpcClient.UserID))
		err := os.MkdirAll(userDir, 0700)
		if err != nil {
			t.Fatalf("Не удалось создать директорию пользователя: %v", err)
		}

		dataFilePath := filepath.Join(userDir, "data.json")
		err = os.WriteFile(dataFilePath, []byte(""), 0644)

		mockClient.EXPECT().GetAllData(gomock.Any(), &proto.GetAllDataRequest{}).
			Return(&proto.GetAllDataResponse{
				Data: []*proto.DataItem{},
			}, nil)

		err = grpcClient.SyncData()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка получения данных из локального хранилища")

		err = os.Chmod(dataFilePath, 0644)
		if err != nil {
			t.Fatalf("Не удалось восстановить права на файл: %v", err)
		}
	})
	t.Run("Ошибка отправки данных на сервер", func(t *testing.T) {
		t.Cleanup(func() {
			os.RemoveAll("user_data")
		})

		dataId := int64(1)
		updatedAt := time.Now()

		createTestData(dataId, mdata.TextData, []byte("test data"), updatedAt)

		mockClient.EXPECT().GetAllData(gomock.Any(), &proto.GetAllDataRequest{}).
			Return(&proto.GetAllDataResponse{
				Data: []*proto.DataItem{},
			}, nil)

		mockClient.EXPECT().CreateData(gomock.Any(), gomock.Any()).
			Return(nil, fmt.Errorf("ошибка отправки данных на сервер"))

		err := grpcClient.SyncData()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка отправки данных на сервер")
	})
	t.Run("Ошибка обновления данных на сервере", func(t *testing.T) {
		t.Cleanup(func() {
			os.RemoveAll("user_data")
		})

		os.RemoveAll("user_data")

		dataId := int64(1)
		updatedAt := time.Now()

		createTestData(dataId, mdata.TextData, []byte("test data"), updatedAt)

		mockClient.EXPECT().GetAllData(gomock.Any(), &proto.GetAllDataRequest{}).
			Return(&proto.GetAllDataResponse{
				Data: []*proto.DataItem{
					{
						DataId:      dataId,
						DataType:    proto.DataType_TEXT_DATA,
						DataContent: []byte("test data"),
						Metadata:    nil,
						UpdatedAt:   updatedAt.Format(time.RFC3339),
					},
				},
			}, nil)

		mockClient.EXPECT().UpdateData(gomock.Any(), gomock.Any()).
			Return(nil, fmt.Errorf("ошибка обновления данных на сервере"))

		err := grpcClient.SyncData()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка обновления данных на сервере")
	})
	t.Run("Ошибка сохранения данных в локальное хранилище", func(t *testing.T) {
		t.Cleanup(func() {
			os.RemoveAll("user_data")
		})

		os.RemoveAll("user_data")

		dataId := int64(1)
		updatedAt := time.Now()

		mockClient.EXPECT().GetAllData(gomock.Any(), &proto.GetAllDataRequest{}).
			Return(&proto.GetAllDataResponse{
				Data: []*proto.DataItem{
					{
						DataId:      dataId,
						DataType:    proto.DataType_TEXT_DATA,
						DataContent: []byte("test data"),
						Metadata:    nil,
						UpdatedAt:   updatedAt.Format(time.RFC3339),
					},
				},
			}, nil)

		userDir := filepath.Join("user_data", fmt.Sprintf("%d", grpcClient.UserID))
		err := os.MkdirAll(userDir, 0700)
		if err != nil {
			t.Fatalf("Не удалось создать директорию пользователя: %v", err)
		}

		dataFilePath := filepath.Join(userDir, "data.json")
		err = os.WriteFile(dataFilePath, []byte("{}"), 0444)
		if err != nil {
			t.Fatalf("Не удалось установить права на файл: %v", err)
		}

		err = grpcClient.SyncData()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка сохранения данных в локальное хранилище")
	})
	t.Run("Ошибка обновления локальных данных", func(t *testing.T) {
		t.Cleanup(func() {
			os.RemoveAll("user_data")
		})

		os.RemoveAll("user_data")

		dataId := int64(1)
		localUpdatedAt := time.Now().Add(-1 * time.Hour) // Локальные данные старше
		serverUpdatedAt := time.Now()

		createTestData(dataId, mdata.TextData, []byte("old data"), localUpdatedAt)

		mockClient.EXPECT().GetAllData(gomock.Any(), &proto.GetAllDataRequest{}).Return(&proto.GetAllDataResponse{
			Data: []*proto.DataItem{
				{
					DataId:      dataId,
					DataType:    proto.DataType_TEXT_DATA,
					DataContent: []byte("new data"),
					Metadata:    nil,
					UpdatedAt:   serverUpdatedAt.Format(time.RFC3339),
				},
			},
		}, nil)

		userDir := filepath.Join("user_data", fmt.Sprintf("%d", grpcClient.UserID))
		dataFilePath := filepath.Join(userDir, "data.json")
		err := os.Chmod(dataFilePath, 0444)
		if err != nil {
			t.Fatalf("Не удалось восстановить права на файл: %v", err)
		}
		err = grpcClient.SyncData()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка обновления локальных данных")
	})
	t.Run("Ошибка обновления локального ID", func(t *testing.T) {
		t.Cleanup(func() {
			os.RemoveAll("user_data")
		})

		os.RemoveAll("user_data")

		localID := int64(1)
		newID := int64(2) // Новый ID, который вернет сервер
		updatedAt := time.Now()

		createTestData(localID, mdata.TextData, []byte("test data"), updatedAt)

		mockClient.EXPECT().GetAllData(gomock.Any(), &proto.GetAllDataRequest{}).
			Return(&proto.GetAllDataResponse{
				Data: []*proto.DataItem{},
			}, nil)

		mockClient.EXPECT().CreateData(gomock.Any(), gomock.Any()).
			Return(&proto.CreateDataResponse{
				DataId: newID,
			}, nil)

		userDir := filepath.Join("user_data", fmt.Sprintf("%d", grpcClient.UserID))
		dataFilePath := filepath.Join(userDir, "data.json")
		err := os.Chmod(dataFilePath, 0444)
		if err != nil {
			t.Fatalf("Не удалось восстановить права на файл: %v", err)
		}

		err = grpcClient.SyncData()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка обновления локального ID")
	})
}

// InvalidDataType Вспомогательная структура для тестирования ошибки преобразования в JSON
type InvalidDataType struct{}

func (i *InvalidDataType) Validate() error {
	return nil
}

func (i *InvalidDataType) ToJSON() ([]byte, error) {
	return nil, errors.New("ошибка преобразования в JSON")
}
