package app

import (
	"fmt"

	logging "github.com/Sofja96/GophKeeper.git/internal/server/logger"
	"github.com/Sofja96/GophKeeper.git/internal/server/service"
	"github.com/Sofja96/GophKeeper.git/internal/server/settings"
	"github.com/Sofja96/GophKeeper.git/internal/server/storage/db"
	"github.com/Sofja96/GophKeeper.git/internal/server/storage/minio"
)

// Server - интерфейс, предоставляющий доступ к различным компонентам сервера,
type Server interface {
	GetSettings() settings.Settings
	GetDbAdapter() db.Adapter
	GetService() service.Service
	GetLogger() logging.ILogger
	GetMinioClient() minio.Client
}

// server - структура, которая реализует интерфейс Server.
type server struct {
	settings    settings.Settings
	dbAdapter   db.Adapter
	service     service.Service
	logger      logging.ILogger
	minioClient minio.Client
}

// GetSettings возвращает настройки сервера.
func (s *server) GetSettings() settings.Settings {
	return s.settings
}

// GetDbAdapter возвращает адаптер для работы с базой данных.
func (s *server) GetDbAdapter() db.Adapter {
	return s.dbAdapter
}

// GetService возвращает экземпляр сервиса для основной логики.
func (s *server) GetService() service.Service {
	return s.service
}

// GetLogger возвращает логгер для записи логов.
func (s *server) GetLogger() logging.ILogger {
	return s.logger
}

// GetMinioClient возвращает клиент для работы с MinIO.
func (s *server) GetMinioClient() minio.Client {
	return s.minioClient
}

// Run инициализирует все компоненты сервера, включая конфигурацию, базу данных,
// логгер, клиент MinIO и сам сервис. Возвращает экземпляр сервера.
func Run() (Server, error) {
	conf, err := settings.GetSettings()
	if err != nil {
		return nil, fmt.Errorf("error load configuration: %w", err)
	}

	logger := logging.New(conf)

	dbAdapter, err := db.NewAdapter(conf)
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.NewMinioClient(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	return &server{
		settings:    *conf,
		dbAdapter:   dbAdapter,
		logger:      logger,
		minioClient: minioClient,
		service:     service.New(dbAdapter, minioClient, logger),
	}, nil
}
