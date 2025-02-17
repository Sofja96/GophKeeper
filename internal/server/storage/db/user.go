package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Sofja96/GophKeeper.git/internal/models"
)

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
