package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Sofja96/GophKeeper.git/internal/models"
	"github.com/Sofja96/GophKeeper.git/internal/server/settings"
)

//todo покрыть тестами

type dbAdapter struct {
	conn *sqlx.DB
}

func (db *dbAdapter) Close() {
	_ = db.conn.Close()
}

type Adapter interface {
	Close()
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserIDByName(ctx context.Context, username string) (bool, error)
	GetUserHashPassword(ctx context.Context, username string) (string, error)
}

func NewAdapter(settings *settings.Settings) (Adapter, error) {
	db, err := sqlx.Connect("postgres", settings.DbDsn)
	if err != nil {
		return nil, err
	}

	dbClient := dbAdapter{conn: db}

	if settings.DbAutoMigration {
		err = dbClient.migration(settings)
		if err != nil {
			return nil, err
		}
	}
	return &dbClient, err
}

func (db *dbAdapter) migration(settings *settings.Settings) error {
	dbInstance, err := sql.Open("postgres", settings.DbDsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	driver, err := postgres.WithInstance(dbInstance, &postgres.Config{MigrationsTable: "migration"})

	if err != nil {
		return err
	}

	dbDsnParse := strings.Split(settings.DbDsn, "/")
	lastPart := strings.Split(dbDsnParse[len(dbDsnParse)-1], "?")
	m, err := migrate.NewWithDatabaseInstance("file:./internal/server/storage/db/migrations", lastPart[0], driver)
	if err != nil {
		return fmt.Errorf("error creaate migrate: %w", err)
	}

	defer func() { _, _ = m.Close() }()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Database migration completed successfully.")

	return nil
}
