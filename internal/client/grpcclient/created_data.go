package grpcclient

import (
	"fmt"
	"time"

	"github.com/Sofja96/GophKeeper.git/internal/client/encryption"
	"github.com/Sofja96/GophKeeper.git/internal/client/localstorage"
	"github.com/Sofja96/GophKeeper.git/internal/client/models"
	mdata "github.com/Sofja96/GophKeeper.git/internal/models"
	"github.com/Sofja96/GophKeeper.git/proto"
)

// CreateData создает новые данные и сохраняет их в локальное хранилище.
// Функция выполняет валидацию входных данных, преобразует их в JSON, шифрует в зависимости от типа данных,
// а затем сохраняет данные в локальное хранилище с уникальным идентификатором.
func (c *Client) CreateData(reqData models.CreateData) (int64, error) {
	var encryptedData string
	var fileName string

	if err := reqData.Data.Validate(); err != nil {
		return 0, fmt.Errorf("ошибка валидации данных: %w", err)
	}

	rawData, err := reqData.Data.ToJSON()
	if err != nil {
		return 0, fmt.Errorf("ошибка преобразования в JSON: %w", err)
	}

	switch reqData.DataType {
	case proto.DataType_BINARY_DATA:
		if binaryData, ok := reqData.Data.(*models.BinaryDataType); ok {
			encryptedData = encryption.EncodeData(binaryData.Content)
			fileName = binaryData.Filename
		} else {
			return 0, fmt.Errorf("неверный тип данных для бинарного контента")
		}
	default:
		var err error
		if encryptedData, err = encryption.EncryptData(rawData, reqData.EncryptionKey); err != nil {
			return 0, fmt.Errorf("ошибка шифрования: %w", err)
		}
	}

	tempID := time.Now().UnixNano()

	dataType, err := mdata.GetModelType(reqData.DataType)
	if err != nil {
		return 0, fmt.Errorf("ошибка конвертации типа данных: %w", err)
	}

	data := mdata.Data{
		ID:          tempID,
		UserID:      c.UserID,
		DataType:    dataType,
		DataContent: []byte(encryptedData),
		Metadata:    reqData.Metadata.AsMap(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		FileName:    fileName,
	}

	if err := localstorage.SaveData(c.UserID, data); err != nil {
		return 0, fmt.Errorf("ошибка сохранения данных в локальное хранилище: %w", err)
	}

	fmt.Println("Данные успешно сохранены в локальное хранилище.")
	return tempID, nil
}
