package grpcclient

import (
	"fmt"
	"sort"

	"github.com/Sofja96/GophKeeper.git/internal/client/encryption"
	"github.com/Sofja96/GophKeeper.git/internal/client/localstorage"
	"github.com/Sofja96/GophKeeper.git/internal/models"
)

// GetData получает все данные пользователя из локального хранилища,
// расшифровывая их при необходимости, и возвращает их в виде среза данных.
// В случае ошибки при получении или расшифровке данных возвращается ошибка.
func (c *Client) GetData() ([]models.Data, error) {
	dataMap, err := localstorage.GetAllData(c.UserID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных из локального хранилища: %w", err)
	}

	key := c.GetMasterKey()
	var data []models.Data

	for _, item := range dataMap {
		var decryptedData []byte

		if item.DataType == models.BinaryData {
			decryptedData, err = encryption.DecodeData(string(item.DataContent))
			if err != nil {
				decryptedData = item.DataContent
			}
		} else {
			decryptedData, err = encryption.DecryptData(string(item.DataContent), key)
			if err != nil {
				return nil, fmt.Errorf("ошибка расшифровки данных: %w", err)
			}

		}

		data = append(data, models.Data{
			ID:          item.ID,
			DataType:    item.DataType,
			DataContent: decryptedData,
			Metadata:    item.Metadata,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].ID < data[j].ID
	})

	return data, nil
}
