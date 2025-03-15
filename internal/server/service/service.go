package service

import (
	"context"

	"github.com/Sofja96/GophKeeper.git/internal/models"
	logging "github.com/Sofja96/GophKeeper.git/internal/server/logger"
	"github.com/Sofja96/GophKeeper.git/internal/server/storage/db"
	"github.com/Sofja96/GophKeeper.git/internal/server/storage/minio"
)

// Service интерфейс предоставляет методы для работы с пользователями и данными.
// Он включает операции для регистрации, авторизации, создания, получения, удаления и обновления данных.
type Service interface {
	RegisterUser(ctx context.Context, user *models.User) (*models.User, error)
	LoginUser(ctx context.Context, user *models.User) (string, error)
	CreateData(ctx context.Context, data *models.Data) (int64, error)
	GetUserIDByUsername(ctx context.Context, username string) (int64, error)
	GetData(ctx context.Context, userId int64) ([]models.Data, error)
	DeleteData(ctx context.Context, dataId int64, userId int64) (bool, error)
	UpdateData(ctx context.Context, data *models.Data) error
}

type service struct {
	dbAdapter   db.Adapter
	minioClient minio.Client
	logger      logging.ILogger
}

// New создаёт новый экземпляр service с переданными зависимостями.
// Возвращает интерфейс Service, который можно использовать для работы с данными и пользователями.
func New(dbAdapter db.Adapter, minioClient minio.Client, logger logging.ILogger) Service {
	return &service{
		dbAdapter:   dbAdapter,
		minioClient: minioClient,
		logger:      logger,
	}
}
