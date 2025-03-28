package service

import (
	"context"
	"fmt"

	"github.com/Sofja96/GophKeeper.git/internal/models"
	"github.com/Sofja96/GophKeeper.git/internal/server/grpcserver/interceptors"
	"github.com/Sofja96/GophKeeper.git/internal/server/utils"
)

// RegisterUser регистрирует нового пользователя.
// Проверяет, существует ли уже пользователь с данным именем.
// Если существует, возвращает ошибку ErrUserExists.
// Хеширует пароль пользователя перед сохранением в базе данных.
// Возвращает зарегистрированного пользователя или ошибку.
func (s *service) RegisterUser(ctx context.Context, user *models.User) (*models.User, error) {
	existingUser, err := s.dbAdapter.GetUserIDByName(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}
	if existingUser {
		return nil, utils.ErrUserExists
	}

	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user = &models.User{
		Username: user.Username,
		Password: hash,
	}

	newUser, err := s.dbAdapter.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// LoginUser выполняет аутентификацию пользователя.
// Проверяет, существует ли пользователь с данным именем.
// Сравнивает введенный пароль с хешированным паролем в базе данных.
// Если пароль верный, генерирует JWT токен для аутентифицированного пользователя.
// Возвращает JWT токен в виде строки или ошибку.
func (s *service) LoginUser(ctx context.Context, user *models.User) (string, error) {
	existingUser, err := s.dbAdapter.GetUserIDByName(ctx, user.Username)
	if err != nil {
		return "", fmt.Errorf("error checking existing user: %w", err)
	}
	if !existingUser {
		return "", fmt.Errorf("users not found, please to registration")
	}

	hash, err := s.dbAdapter.GetUserHashPassword(ctx, user.Username)
	if err != nil {
		return "", err
	}

	err = utils.CheckPassword(user.Password, hash)
	if err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}
	token, err := interceptors.CreateToken(user.Username)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %w", err)
	}
	var bearer = "Bearer " + token

	return bearer, nil

}

// GetUserIDByUsername возвращает ID пользователя по его имени.
// Используется для получения уникального идентификатора пользователя из базы данных.
func (s *service) GetUserIDByUsername(ctx context.Context, username string) (int64, error) {
	return s.dbAdapter.GetUserID(ctx, username)
}
