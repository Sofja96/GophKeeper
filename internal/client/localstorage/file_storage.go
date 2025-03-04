package localstorage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Sofja96/GophKeeper.git/internal/models"
)

const storageFile = "local_storage.json"

// Storage представляет собой структуру для хранения данных.
type Storage struct {
	Data map[int64]map[int64]models.Data `json:"data"`
}

// SaveData сохраняет данные в локальное хранилище
func SaveData(userID int64, data models.Data) error {
	storage, err := readData()
	if err != nil {
		return fmt.Errorf("ошибка чтения данных: %w", err)
	}

	if storage.Data[userID] == nil {
		storage.Data[userID] = make(map[int64]models.Data)
	}

	storage.Data[userID][data.ID] = data
	return writeData(storage)
}

// GetAllData возвращает все данные из локального хранилища для указанного пользователя
func GetAllData(userID int64) (map[int64]models.Data, error) {
	storage, err := readData()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения данных: %w", err)
	}
	return storage.Data[userID], nil
}

// readData читает данные из файла
func readData() (*Storage, error) {
	storage := &Storage{Data: make(map[int64]map[int64]models.Data)}

	if _, err := os.Stat(storageFile); os.IsNotExist(err) {
		return storage, nil
	}

	file, err := os.ReadFile(storageFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	if err := json.Unmarshal(file, storage); err != nil {
		return nil, fmt.Errorf("ошибка десериализации данных: %w", err)
	}

	return storage, nil
}

// writeData записывает данные в файл
func writeData(storage *Storage) error {
	file, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %w", err)
	}

	if err := os.WriteFile(storageFile, file, 0644); err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}

	return nil
}

// DeleteData удаляет данные из локального хранилища по userID и dataID
func DeleteData(userID, dataID int64) error {
	storage, err := readData()
	if err != nil {
		return fmt.Errorf("ошибка чтения данных: %w", err)
	}

	if _, exists := storage.Data[userID]; !exists {
		return fmt.Errorf("пользователь с ID %d не найден", userID)
	}

	delete(storage.Data[userID], dataID)
	return writeData(storage)
}

// UpdateID обновляет ID записи в локальном хранилище
func UpdateID(userID, oldID, newID int64) error {
	storage, err := readData()
	if err != nil {
		return fmt.Errorf("ошибка чтения данных: %w", err)
	}

	userData, exists := storage.Data[userID]
	if !exists {
		return fmt.Errorf("пользователь с ID %d не найден", userID)
	}

	data, exists := userData[oldID]
	if !exists {
		return fmt.Errorf("данные с ID %d не найдены", oldID)
	}

	delete(userData, oldID)
	data.ID = newID
	userData[newID] = data
	storage.Data[userID] = userData

	return writeData(storage)
}
