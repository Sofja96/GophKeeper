package app

import (
	"context"
	"fmt"

	logging "github.com/Sofja96/GophKeeper.git/internal/server/logger"
	"github.com/Sofja96/GophKeeper.git/internal/server/service"
	"github.com/Sofja96/GophKeeper.git/internal/server/settings"
	"github.com/Sofja96/GophKeeper.git/internal/server/storage/db"
)

type Server interface {
	GetContext() context.Context
	GetSettings() settings.Settings
	GetDbAdapter() db.Adapter    //добавить интерфейс ДБ
	GetService() service.Service //добавить сервис интерфейс
	GetLogger() logging.ILogger
	//GetMinio() // если будет структруры или интерфейс
}

type server struct {
	ctx       context.Context
	settings  settings.Settings
	dbAdapter db.Adapter
	service   service.Service
	logger    logging.ILogger
	//todo возможно логгер или его оставить только как мидлваре
}

func (s *server) GetContext() context.Context {
	return s.ctx
}

func (s *server) GetSettings() settings.Settings {
	return s.settings
}

func (s *server) GetDbAdapter() db.Adapter {
	return s.dbAdapter
}

func (s *server) GetService() service.Service {
	return s.service
}

func (s *server) GetLogger() logging.ILogger {
	return s.logger
}

// todo moжно переименовать в New и ниже добавить Start(в которой будет инициализироваться сам grpc)

func Run() (Server, error) {
	ctx := context.Background()

	conf, err := settings.GetSettings()
	if err != nil {
		return nil, fmt.Errorf("error load configuration: %w", err)
	}

	logger := logging.New(conf)

	dbAdapter, err := db.NewAdapter(conf)
	if err != nil {
		return nil, err
	}

	return &server{
		ctx:       ctx,
		settings:  *conf,
		dbAdapter: dbAdapter,
		service:   service.New(dbAdapter),
		logger:    logger,
	}, nil
}
