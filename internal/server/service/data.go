package service

import (
	"context"
	"fmt"

	"github.com/Sofja96/GophKeeper.git/internal/models"
	"github.com/Sofja96/GophKeeper.git/internal/server/utils"
)

// CreateData создает новые данные в базе данных и (если необходимо) загружает бинарные данные в MinIO.
// Если данные являются бинарными, файл загружается в MinIO, и его URL сохраняется в метаданных.
func (s *service) CreateData(ctx context.Context, data *models.Data) (int64, error) {
	putData := *data

	if putData.DataType == models.BinaryData {
		fileURL, err := s.minioClient.UploadFile(ctx, data.FileName, data.DataContent)
		if err != nil {
			s.logger.Error("ошибка загрузки в Minio: %v", err)
			return 0, err
		}

		putData.SetMetadata("file_url", fileURL)

		putData.DataContent = nil
	} else {
		putData.DataContent = data.DataContent
	}

	return s.dbAdapter.CreateData(ctx, &putData)
}

// GetData получает все данные для указанного пользователя. Если данные являются бинарными,
// они загружаются из MinIO с использованием URL, сохраненного в метаданных.
func (s *service) GetData(ctx context.Context, userId int64) ([]models.Data, error) {
	data, err := s.dbAdapter.GetData(ctx, userId)
	if len(data) == 0 && err == nil {
		return nil, utils.ErrUserDataNotFound
	}

	for i := range data {
		if data[i].DataType == models.BinaryData {
			fileURLValue, ok := data[i].GetMetadata("file_url")
			if !ok {
				s.logger.Error("file_url not found in metadata for binary data")
				continue
			}

			fileURL, ok := fileURLValue.(string)
			if !ok || fileURL == "" {
				s.logger.Error("file_url is not a valid string")
				continue
			}

			fileContent, err := s.minioClient.GetFile(ctx, fileURL)
			if err != nil {
				s.logger.Error("failed to load file from MinIO: %v", err)
				continue
			}

			data[i].DataContent = fileContent
		}
	}

	return data, err
}

// DeleteData удаляет данные с заданным идентификатором (dataId) для указанного пользователя (userId).
// Если данные бинарные, соответствующий файл также удаляется из MinIO.
func (s *service) DeleteData(ctx context.Context, dataId int64, userId int64) (bool, error) {
	data, err := s.dbAdapter.GetDataByID(ctx, dataId)
	if err != nil {
		return false, err
	}

	if data.DataType == models.BinaryData {
		fileURL, ok := data.Metadata["file_url"].(string)
		if !ok || fileURL == "" {
			return false, fmt.Errorf("file_url не найден в метаданных")
		}

		err := s.minioClient.DeleteFile(ctx, fileURL)
		if err != nil {
			return false, fmt.Errorf("ошибка удаления файла из MinIO: %w", err)
		}
	}

	success, err := s.dbAdapter.DeleteData(ctx, dataId, userId)
	if err != nil {
		return false, fmt.Errorf("ошибка удаления данных из базы данных: %w", err)
	}

	if !success {
		return false, utils.ErrUserDataNotFound
	}

	return success, nil
}

// UpdateData обновляет данные с заданным идентификатором (dataId) для указанного пользователя.
// Если данные бинарные, файл обновляется в MinIO.
func (s *service) UpdateData(ctx context.Context, data *models.Data) error {
	oldData, err := s.dbAdapter.GetDataByID(ctx, data.ID)
	if err != nil {
		return err
	}

	if oldData.DataType == models.BinaryData {
		OldFileURL, ok := oldData.Metadata["file_url"].(string)
		if !ok || OldFileURL == "" {
			return fmt.Errorf("file_url не найден в метаданных")
		}

		minioUrl, err := s.minioClient.UpdateFile(ctx, OldFileURL, data.FileName, data.DataContent)
		if err != nil {
			return err
		}

		data.SetMetadata("file_url", minioUrl)
		data.DataContent = nil
	}

	err = s.dbAdapter.UpdateData(ctx, data)
	if err != nil {
		return err
	}

	return nil

}
