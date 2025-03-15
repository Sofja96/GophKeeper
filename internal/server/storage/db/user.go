package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Sofja96/GophKeeper.git/internal/models"
)

// CreateUser создает нового пользователя в базе данных.
//
// Если пользователь с таким именем уже существует, его пароль будет обновлен.
// Возвращает созданного пользователя или ошибку, если операция не удалась.
func (db *dbAdapter) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `insert into users (
                   username, password) values ($1, $2) on conflict(username) do update set
                   password = EXCLUDED.password;`

	_, err := db.conn.ExecContext(ctx, query, user.Username, user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, err

}

// GetUserIDByName проверяет, существует ли пользователь с указанным именем.
//
// Возвращает true, если пользователь найден, и false, если не найден,
// а также ошибку, если произошла ошибка при выполнении запроса.
func (db *dbAdapter) GetUserIDByName(ctx context.Context, username string) (bool, error) {
	var id string
	query := `SELECT id FROM users WHERE username = $1`

	err := db.conn.GetContext(ctx, &id, query, username)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("error getting id users: %w", err)
	}

	return true, nil
}

// GetUserHashPassword получает хеш пароля пользователя по его имени.
//
// Возвращает хеш пароля пользователя или ошибку, если он не найден.
func (db *dbAdapter) GetUserHashPassword(ctx context.Context, username string) (string, error) {
	var password string

	query := `SELECT password FROM users WHERE username = $1 `

	err := db.conn.GetContext(ctx, &password, query, username)
	if errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	if err != nil {
		return "", fmt.Errorf("error getting password on user: %w", err)
	}

	return password, nil
}

// GetUserID возвращает ID пользователя по его имени.
//
// Возвращает ID пользователя или ошибку, если имя не найдено в базе данных.
func (db *dbAdapter) GetUserID(ctx context.Context, username string) (int64, error) {
	var id int64
	row := db.conn.QueryRowContext(ctx, "SELECT id FROM users WHERE username = $1", username)
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("unable select id: %w", err)
	}

	return id, nil
}
