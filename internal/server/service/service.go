package service

import (
	"context"

	"github.com/Sofja96/GophKeeper.git/internal/models"
	"github.com/Sofja96/GophKeeper.git/internal/server/storage/db"
)

// todo покрыть тестами и вызывать в хендлере
type Service interface {
	RegisterUser(ctx context.Context, user *models.User) (*models.User, error)
	LoginUser(ctx context.Context, user *models.User) (string, error)
}

type service struct {
	dbAdapter db.Adapter
}

func New(dbAdapter db.Adapter) Service {
	return &service{
		dbAdapter: dbAdapter,
	}
}
