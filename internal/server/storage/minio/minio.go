package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/Sofja96/GophKeeper.git/internal/server/settings"
)

// Client представляет интерфейс для работы с MinIO: загрузка, получение, обновление и удаление файлов.
type Client interface {
	UploadFile(ctx context.Context, fileName string, fileContent []byte) (string, error)
	DeleteFile(ctx context.Context, fileURL string) error
	GetFile(ctx context.Context, fileURL string) ([]byte, error)
	UpdateFile(ctx context.Context, oldFileName, fileName string, content []byte) (string, error)
}

type client struct {
	Client *minio.Client
	Bucket string
}

// NewMinioClient Инициализация клиента MinIO
func NewMinioClient(settings *settings.Settings) (Client, error) {
	minioClient, err := minio.New(settings.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(settings.MinioUser, settings.MinioPassword, ""),
		Secure: settings.MinioUseSsl,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	ctx := context.Background()
	location := "us-east-1"
	exists, errBucketExists := minioClient.BucketExists(ctx, settings.MinioBucketName)
	if errBucketExists != nil {
		return nil, errBucketExists
	}
	if !exists {
		if err = minioClient.MakeBucket(ctx, settings.MinioBucketName, minio.MakeBucketOptions{Region: location}); err != nil {
			return nil, err
		}
		log.Printf("Создан новый bucket: %s\n", settings.MinioBucketName)
	}

	return &client{
		Client: minioClient,
		Bucket: settings.MinioBucketName,
	}, nil
}

// UploadFile загружает файл в MinIO.
//
// Принимает имя файла и его содержимое в виде байтов. Возвращает URL загруженного файла или ошибку,
// если загрузка не удалась.
func (m *client) UploadFile(ctx context.Context, fileName string, fileContent []byte) (string, error) {
	objectName := fmt.Sprintf("uploads/%s", fileName)
	_, err := m.Client.PutObject(ctx, m.Bucket, objectName, bytes.NewReader(fileContent), int64(len(fileContent)), minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	fileURL := fmt.Sprintf("%s/%s/%s", m.Client.EndpointURL(), m.Bucket, objectName)
	return fileURL, nil
}

// DeleteFile удаляет файл из MinIO по его URL.
//
// Принимает URL файла и удаляет его из хранилища MinIO. Возвращает ошибку, если файл не удается удалить.
func (m *client) DeleteFile(ctx context.Context, fileURL string) error {
	objectName := strings.TrimPrefix(fileURL, fmt.Sprintf("%s/%s/", m.Client.EndpointURL(), m.Bucket))
	err := m.Client.RemoveObject(ctx, m.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetFile загружает файл из MinIO по его URL.
//
// Принимает URL файла, извлекает содержимое и возвращает его в виде байтов. Возвращает ошибку,
// если файл не удается получить или URL пуст.
func (m *client) GetFile(ctx context.Context, fileURL string) ([]byte, error) {
	if fileURL == "" {
		return nil, fmt.Errorf("fileURL is empty")
	}
	objectName := strings.TrimPrefix(fileURL, fmt.Sprintf("%s/%s/", m.Client.EndpointURL(), m.Bucket))

	object, err := m.Client.GetObject(ctx, m.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file from MinIO: %w", err)
	}
	defer object.Close()

	fileContent, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return fileContent, nil
}

// UpdateFile обновляет файл в MinIO, если имя файла изменилось.
//
// Если имя старого файла совпадает с новым, то файл просто загружается заново, иначе старый файл удаляется,
// а новый загружается в хранилище. Возвращает URL нового файла или ошибку.
func (m *client) UpdateFile(ctx context.Context, oldFileName, fileName string, content []byte) (string, error) {
	oldFileName = strings.TrimPrefix(oldFileName, fmt.Sprintf("%s/%s/", m.Client.EndpointURL(), m.Bucket))
	if oldFileName == fileName {
		FileUrl, err := m.UploadFile(ctx, fileName, content)
		if err != nil {
			return "", err
		}
		return FileUrl, nil
	}

	err := m.DeleteFile(ctx, oldFileName)
	if err != nil {
		return "", err
	}

	fileUrl, err := m.UploadFile(ctx, fileName, content)
	if err != nil {
		return "", err
	}
	return fileUrl, nil
}
