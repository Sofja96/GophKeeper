package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Sofja96/GophKeeper.git/internal/models"
)

// CreateData создает новую запись данных в базе данных.
//
// Эта функция использует транзакцию для создания записи данных. Она вставляет запись в таблицу данных,
// включая данные пользователя, тип данных, содержимое и метаданные. В случае успеха возвращает ID созданной записи.
//
// Если при создании данных происходит ошибка, транзакция будет отменена, и функция вернет ошибку.
func (db *dbAdapter) CreateData(ctx context.Context, data *models.Data) (int64, error) {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() { _ = tx.Rollback() }()

	query := `insert into data(user_id, data_type, data_content, metadata)
			values ($1, $2, $3, $4) RETURNING id`

	var id int64
	err = tx.QueryRowContext(ctx, query, data.UserID, data.DataType, data.DataContent, data.Metadata).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert data: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetData получает все данные пользователя из базы данных по его ID.
//
// Функция извлекает список данных пользователя по его ID, сортируя записи по дате создания.
// Если пользователь не имеет данных, возвращает пустой срез.
func (db *dbAdapter) GetData(ctx context.Context, userId int64) ([]models.Data, error) {
	dataList := make([]models.Data, 0)

	query := `select id, data_type, data_content, metadata, updated_at
			 from data where user_id = $1 order by created_at`

	err := db.conn.SelectContext(ctx, &dataList, query, userId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error getting data info: %w", err)
	}

	return dataList, err
}

// GetDataByID получает данные по их ID.
//
// Функция извлекает конкретную запись данных по указанному ID из базы данных.
// Возвращает ошибку, если запись не найдена или произошла другая ошибка при извлечении.
func (db *dbAdapter) GetDataByID(ctx context.Context, dataID int64) (*models.Data, error) {
	query := `SELECT id, user_id, data_type, data_content, metadata, updated_at 
	          FROM data WHERE id = $1`

	var data models.Data
	err := db.conn.GetContext(ctx, &data, query, dataID)
	if err != nil {
		return nil, fmt.Errorf("error getting data by ID: %w", err)
	}

	return &data, nil
}

// DeleteData удаляет запись данных из базы данных по ID и ID пользователя.
//
// Функция удаляет данные, если они принадлежат указанному пользователю (проверка по ID).
// Возвращает true, если запись была успешно удалена, и false, если запись не найдена или произошла ошибка.
func (db *dbAdapter) DeleteData(ctx context.Context, dataId int64, userId int64) (bool, error) {
	query := `delete from data where id = $1 and user_id= $2`

	result, err := db.conn.ExecContext(ctx, query, dataId, userId)
	if err != nil {
		return false, fmt.Errorf("error deleting data: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error get count rows: %w", err)
	}

	return rowsAffected > 0, nil
}

// UpdateData обновляет существующую запись данных в базе данных.
//
// Эта функция использует транзакцию для обновления записи данных. Обновляются поля содержимого данных и метаданных.
// Если транзакция успешна, изменения сохраняются в базе данных, если произошла ошибка — транзакция откатывается.
func (db *dbAdapter) UpdateData(ctx context.Context, data *models.Data) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	query := `update data
			 set data_content = $1, metadata = $2, updated_at = now()
             where id = $3 and user_id = $4`

	_, err = tx.ExecContext(ctx, query, data.DataContent, data.Metadata, data.ID, data.UserID)
	if err != nil {
		return fmt.Errorf("error update update data: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
