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

// UpdateData обновляет данные в локальном хранилище на основе переданных данных.
// Функция выполняет валидацию входных данных, преобразует их в JSON, шифрует в зависимости от типа данных,
// и затем сохраняет обновленные данные в локальное хранилище.
func (c *Client) UpdateData(reqData models.CreateData, dataId int64) error {
	var encryptedData string
	var fileName string

	if err := reqData.Data.Validate(); err != nil {
		return fmt.Errorf("ошибка валидации данных: %w", err)
	}

	rawData, err := reqData.Data.ToJSON()
	if err != nil {
		return fmt.Errorf("ошибка преобразования в JSON: %w", err)
	}

	switch reqData.DataType {
	case proto.DataType_BINARY_DATA:
		if binaryData, ok := reqData.Data.(*models.BinaryDataType); ok {
			encryptedData = encryption.EncodeData(binaryData.Content)
			fileName = binaryData.Filename
		} else {
			return fmt.Errorf("неверный тип данных для бинарного контента")
		}
	default:
		var err error
		if encryptedData, err = encryption.EncryptData(rawData, reqData.EncryptionKey); err != nil {
			return fmt.Errorf("ошибка шифрования: %w", err)
		}
	}

	dataType, err := mdata.GetModelType(reqData.DataType)
	if err != nil {
		return fmt.Errorf("ошибка конвертации типа данных: %w", err)
	}

	data := mdata.Data{
		UserID:      c.UserID,
		ID:          dataId,
		DataType:    dataType,
		DataContent: []byte(encryptedData),
		Metadata:    reqData.Metadata.AsMap(),
		UpdatedAt:   time.Now(),
		FileName:    fileName,
	}

	if err := localstorage.SaveData(c.UserID, data); err != nil {
		return fmt.Errorf("ошибка обновления данных в локальном хранилище: %w", err)
	}

	fmt.Println("Данные успешно обновлены в локальном хранилище.")
	return nil
}
