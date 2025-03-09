package localstorage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sofja96/GophKeeper.git/internal/models"
)

// Storage представляет собой структуру для хранения данных пользователя.
type Storage struct {
	Data map[int64]models.Data `json:"data"`
}

// getUserDir возвращает путь к папке пользователя.
func getUserDir(userID int64) string {
	return filepath.Join("user_data", fmt.Sprintf("%d", userID))
}

// getUserDataPath возвращает путь к файлу данных пользователя.
func getUserDataPath(userID int64) string {
	return filepath.Join(getUserDir(userID), "data.json")
}

// SaveData сохраняет данные в локальное хранилище.
func SaveData(userID int64, data models.Data) error {
	userDir := getUserDir(userID)
	if err := os.MkdirAll(userDir, 0700); err != nil {
		return fmt.Errorf("ошибка создания папки пользователя: %w", err)
	}

	storage, err := readUserData(userID)
	if err != nil {
		return fmt.Errorf("ошибка чтения данных пользователя: %w", err)
	}

	storage.Data[data.ID] = data

	return writeUserData(userID, storage)
}

// GetAllData возвращает все данные из локального хранилища для указанного пользователя.
func GetAllData(userID int64) (map[int64]models.Data, error) {
	storage, err := readUserData(userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения данных пользователя: %w", err)
	}
	return storage.Data, nil
}

// DeleteData удаляет данные из локального хранилища по userID и dataID.
func DeleteData(userID, dataID int64) error {
	storage, err := readUserData(userID)
	if err != nil {
		return fmt.Errorf("ошибка чтения данных пользователя: %w", err)
	}

	if _, exists := storage.Data[dataID]; !exists {
		return fmt.Errorf("данные с ID %d не найдены", dataID)
	}

	delete(storage.Data, dataID)
	return writeUserData(userID, storage)
}

// UpdateID обновляет ID записи в локальном хранилище.
func UpdateID(userID, oldID, newID int64) error {
	storage, err := readUserData(userID)
	if err != nil {
		return fmt.Errorf("ошибка чтения данных пользователя: %w", err)
	}

	data, exists := storage.Data[oldID]
	if !exists {
		return fmt.Errorf("данные с ID %d не найдены", oldID)
	}

	delete(storage.Data, oldID)
	data.ID = newID
	storage.Data[newID] = data

	return writeUserData(userID, storage)
}

// readUserData читает данные пользователя из файла.
func readUserData(userID int64) (*Storage, error) {
	filePath := getUserDataPath(userID)
	storage := &Storage{Data: make(map[int64]models.Data)}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return storage, nil
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	if err := json.Unmarshal(file, storage); err != nil {
		return nil, fmt.Errorf("ошибка десериализации данных: %w", err)
	}

	return storage, nil
}

// writeUserData записывает данные пользователя в файл.
func writeUserData(userID int64, storage *Storage) error {
	filePath := getUserDataPath(userID)
	file, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %w", err)
	}

	if err := os.WriteFile(filePath, file, 0644); err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}

	return nil
}
