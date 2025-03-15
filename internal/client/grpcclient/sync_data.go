package grpcclient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Sofja96/GophKeeper.git/internal/client/localstorage"
	"github.com/Sofja96/GophKeeper.git/internal/models"
	"github.com/Sofja96/GophKeeper.git/proto"
)

// SyncData синхронизирует данные между сервером и клиентом
func (c *Client) SyncData() error {
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", c.GetToken())

	// 1. Получаем данные с сервера
	serverData, err := c.GetAllDataFromServer(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			fmt.Println("Нет данных на сервере, продолжаем синхронизацию...")
			serverData = make(map[int64]models.Data)
		} else {
			return fmt.Errorf("ошибка получения данных с сервера: %w", err)
		}
	}

	// 2. Получаем данные из локального хранилища
	localData, err := localstorage.GetAllData(c.UserID)
	if err != nil {
		return fmt.Errorf("ошибка получения данных из локального хранилища: %w", err)
	}

	// 3. Синхронизация данных с сервера в локальное хранилище
	for serverID, serverItem := range serverData {
		localItem, exists := localData[serverID]

		if !exists {
			fmt.Println("Синхронизация данных с сервера в локальное хранилище")
			// Данных нет локально — сохраняем их в локальное хранилище
			if err := localstorage.SaveData(c.UserID, serverItem); err != nil {
				return fmt.Errorf("ошибка сохранения данных в локальное хранилище: %w", err)
			}
		} else if serverItem.UpdatedAt.After(localItem.UpdatedAt) {
			// Серверные данные новее — обновляем локальные
			if err := localstorage.SaveData(c.UserID, serverItem); err != nil {
				return fmt.Errorf("ошибка обновления локальных данных: %w", err)
			}
		}
	}

	// 4. Синхронизация данных из локального хранилища на сервер
	for localID, localItem := range localData {

		serverItem, exists := serverData[localID]

		if !exists {
			// Данных нет на сервере — отправляем их на сервер
			newID, err := c.SendDataToServer(ctx, localItem)
			if err != nil {
				return fmt.Errorf("ошибка отправки данных на сервер: %w", err)
			}

			// Обновляем локальный ID на новый, который вернул сервер
			if err := localstorage.UpdateID(c.UserID, localID, newID); err != nil {
				return fmt.Errorf("ошибка обновления локального ID: %w", err)
			}
		} else if localItem.UpdatedAt.After(serverItem.UpdatedAt) {
			// Локальные данные новее — обновляем их на сервере
			if err := c.UpdateDataOnServer(ctx, localItem, localID); err != nil {
				return fmt.Errorf("ошибка обновления данных на сервере: %w", err)
			}
		}
	}

	fmt.Println("Данные успешно синхронизированы.")
	return nil
}

// GetAllDataFromServer получает данные с сервера
func (c *Client) GetAllDataFromServer(ctx context.Context) (map[int64]models.Data, error) {
	resp, err := c.Client.GetAllData(ctx, &proto.GetAllDataRequest{})
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных с сервера: %w", err)
	}

	serverData := make(map[int64]models.Data)
	for _, item := range resp.Data {
		dataType, err := models.GetModelType(item.DataType)
		if err != nil {
			return nil, fmt.Errorf("ошибка конвертации типа данных: %w", err)
		}

		location, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			return nil, fmt.Errorf("ошибка загрузки часового пояса: %w", err)
		}

		updatedAt, err := time.ParseInLocation(time.RFC3339, item.UpdatedAt, location)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования времени UpdatedAt: %w", err)
		}

		serverData[item.DataId] = models.Data{
			ID:          item.DataId,
			DataType:    dataType,
			DataContent: item.DataContent,
			Metadata:    item.Metadata.AsMap(),
			UpdatedAt:   updatedAt,
		}
	}

	return serverData, nil
}

// SendDataToServer отправляет данные на сервер
func (c *Client) SendDataToServer(ctx context.Context, data models.Data) (int64, error) {
	structMetadata, err := models.ConvertJSONBToStruct(data.Metadata)
	if err != nil {
		return 0, fmt.Errorf("ошибка преобразования метаданных: %w", err)
	}

	req := &proto.CreateDataRequest{
		DataType:    proto.DataType(proto.DataType_value[data.DataType.String()]),
		DataContent: data.DataContent,
		Metadata:    structMetadata,
		FileName:    data.FileName,
	}

	resp, err := c.Client.CreateData(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("ошибка отправки данных на сервер: %w", err)
	}

	return resp.DataId, nil
}

// UpdateDataOnServer обновляет данные на сервере
func (c *Client) UpdateDataOnServer(ctx context.Context, data models.Data, dataId int64) error {

	structMetadata, err := models.ConvertJSONBToStruct(data.Metadata)
	if err != nil {
		return fmt.Errorf("ошибка преобразования метаданных: %w", err)
	}

	req := &proto.UpdateDataRequest{
		DataId:      dataId,
		DataContent: data.DataContent,
		Metadata:    structMetadata,
		FileName:    data.FileName,
	}

	_, err = c.Client.UpdateData(ctx, req)
	if err != nil {
		return fmt.Errorf("ошибка обновления данных на сервере: %w", err)
	}

	return nil
}
