package db

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/Sofja96/GophKeeper.git/internal/models"
)

func TestCreateData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	defer db.Close()

	type (
		args struct {
			data *models.Data
		}
		mockBehavior func(m *mocks, args args)
	)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expectedID   int64
		wantErr      bool
		err          error
	}{
		{
			name: "CreateDataSuccessfully",
			args: args{
				data: &models.Data{
					UserID:      10,
					DataType:    "LOGIN_PASSWORD",
					DataContent: []byte("$1$212345"),
					Metadata:    nil,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			mockBehavior: func(m *mocks, args args) {
				mock.ExpectBegin()
				expectedQuery := `insert into data(user_id, data_type, data_content, metadata)
			values ($1, $2, $3, $4) RETURNING id`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(
						args.data.UserID,
						args.data.DataType,
						args.data.DataContent,
						args.data.Metadata).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectedID: 1,
			wantErr:    false,
			err:        nil,
		},
		{
			name: "InvalidTransactionStart",
			args: args{
				data: &models.Data{
					UserID:      10,
					DataType:    "LOGIN_PASSWORD",
					DataContent: []byte("$1$212345"),
					Metadata:    nil,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			mockBehavior: func(m *mocks, args args) {
				mock.ExpectBegin().WillReturnError(
					fmt.Errorf("failed to begin transaction"))
			},
			expectedID: 0,
			wantErr:    true,
			err:        fmt.Errorf("failed to begin transaction"),
		},
		{
			name: "InvalidCommitTransaction",
			args: args{
				data: &models.Data{
					UserID:      10,
					DataType:    "LOGIN_PASSWORD",
					DataContent: []byte("$1$212345"),
					Metadata:    nil,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			mockBehavior: func(m *mocks, args args) {
				mock.ExpectBegin()
				expectedQuery := `insert into data(user_id, data_type, data_content, metadata)
			values ($1, $2, $3, $4) RETURNING id`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(
						args.data.UserID,
						args.data.DataType,
						args.data.DataContent,
						args.data.Metadata).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit().WillReturnError(fmt.Errorf("failed to commit transaction"))
			},
			expectedID: 0,
			wantErr:    true,
			err:        fmt.Errorf("failed to commit transaction"),
		},
		{
			name: "CreateDataError",
			args: args{
				data: &models.Data{
					DataType:    "LOGIN_PASSWORD",
					DataContent: []byte("$1$212345"),
					Metadata:    nil,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			mockBehavior: func(m *mocks, args args) {
				mock.ExpectBegin()
				expectedQuery := `insert into data(user_id, data_type, data_content, metadata)
			values ($1, $2, $3, $4) RETURNING id`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(
						args.data.UserID,
						args.data.DataType,
						args.data.DataContent,
						args.data.Metadata).
					WillReturnError(fmt.Errorf("failed to insert data"))

				mock.ExpectCommit()
			},
			expectedID: 0,
			wantErr:    true,
			err:        fmt.Errorf("failed to insert data"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			m := &mocks{}

			m.db = sqlx.NewDb(db, "sqlmock")
			pg := dbAdapter{conn: m.db}

			tt.mockBehavior(m, tt.args)
			returnedID, err := pg.CreateData(context.Background(), tt.args.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, returnedID, "The returned ID %d does not match the expected ID %d", tt.expectedID, returnedID)
			}

		})
	}
}

func TestGetData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	defer db.Close()

	type (
		args struct {
			userId int64
		}
		mockBehavior func(m *mocks, args args)
	)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expectedData []models.Data
		wantErr      bool
		err          error
	}{
		{
			name: "GetDataSuccessfully",
			args: args{
				userId: 1,
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `select id, data_type, 
       								data_content, metadata, updated_at
			 						from data where user_id = $1 order by created_at`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.userId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "data_type"}).AddRow(1, "LOGIN_PASSWORD"))
			},
			expectedData: []models.Data{
				{
					ID:          1,
					DataType:    "LOGIN_PASSWORD",
					DataContent: []byte(nil),
					Metadata:    nil,
					UpdatedAt:   time.Time{},
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "DataNotFound",
			args: args{
				userId: 11,
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `select id, data_type, 
       								data_content, metadata, updated_at
			 						from data where user_id = $1 order by created_at`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.userId).
					WillReturnError(sql.ErrNoRows)
			},
			expectedData: nil,
			wantErr:      false,
			err:          nil,
		},
		{
			name: "GetDataError",
			args: args{
				userId: 11,
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `select id, data_type, 
       								data_content, metadata, updated_at
			 						from data where user_id = $1 order by created_at`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.userId).
					WillReturnError(fmt.Errorf("error getting data info"))
			},
			expectedData: nil,
			wantErr:      true,
			err:          fmt.Errorf("error getting data info"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			m := &mocks{}

			m.db = sqlx.NewDb(db, "sqlmock")
			pg := dbAdapter{conn: m.db}

			tt.mockBehavior(m, tt.args)
			returnedData, err := pg.GetData(context.Background(), tt.args.userId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedData, returnedData, "The returned Data does not match the expected Data")
			}

		})
	}
}

func TestGetDataByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	defer db.Close()

	type (
		args struct {
			dataId int64
		}
		mockBehavior func(m *mocks, args args)
	)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expectedData *models.Data
		wantErr      bool
		err          error
	}{
		{
			name: "GetDataSuccessfully",
			args: args{
				dataId: 1,
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `SELECT id, user_id, data_type, 
               				data_content, metadata, updated_at 
	          				FROM data WHERE id = $1`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.dataId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "data_type"}).AddRow(1, "LOGIN_PASSWORD"))
			},
			expectedData: &models.Data{
				ID:          1,
				DataType:    "LOGIN_PASSWORD",
				DataContent: []byte(nil),
				Metadata:    nil,
				UpdatedAt:   time.Time{},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "GetDataError",
			args: args{
				dataId: 11,
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `SELECT id, user_id, data_type, 
               				data_content, metadata, updated_at 
	          				FROM data WHERE id = $1`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.dataId).
					WillReturnError(fmt.Errorf("error getting data by ID"))
			},
			expectedData: nil,
			wantErr:      true,
			err:          fmt.Errorf("error getting data by ID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			m := &mocks{}

			m.db = sqlx.NewDb(db, "sqlmock")
			pg := dbAdapter{conn: m.db}

			tt.mockBehavior(m, tt.args)
			returnedData, err := pg.GetDataByID(context.Background(), tt.args.dataId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedData, returnedData, "The returned Data does not match the expected Data")
			}

		})
	}
}

func TestDeleteData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	defer db.Close()

	type (
		args struct {
			dataId int64
			userId int64
		}
		mockBehavior func(m *mocks, args args)
	)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expectedDel  bool
		wantErr      bool
		err          error
	}{
		{
			name: "DeleteDataSuccessfully",
			args: args{
				dataId: 1,
				userId: 3,
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `delete from data where id = $1 and user_id= $2`
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.dataId, args.userId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedDel: true,
			wantErr:     false,
			err:         nil,
		},
		{
			name: "DeleteDataError",
			args: args{
				dataId: 1,
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `delete from data where id = $1 and user_id= $2`
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.dataId, args.userId).
					WillReturnError(fmt.Errorf("error deleting data"))
			},
			expectedDel: false,
			wantErr:     true,
			err:         fmt.Errorf("error deleting data"),
		},
		{
			name: "DeleteDataErrorRowsAffected",
			args: args{
				dataId: 1,
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `delete from data where id = $1 and user_id= $2`
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.dataId, args.userId).
					WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("error getting rows affected")))
			},
			expectedDel: false,
			wantErr:     true,
			err:         fmt.Errorf("error get count rows: error getting rows affected"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			m := &mocks{}

			m.db = sqlx.NewDb(db, "sqlmock")
			pg := dbAdapter{conn: m.db}

			tt.mockBehavior(m, tt.args)
			returned, err := pg.DeleteData(context.Background(), tt.args.dataId, tt.args.userId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedDel, returned, "The returned Data does not match the expected Data")
			}

		})
	}
}

func TestUpdateData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	defer db.Close()

	type (
		args struct {
			data *models.Data
		}
		mockBehavior func(m *mocks, args args)
	)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
		err          error
	}{
		{
			name: "DeleteDataSuccessfully",
			args: args{
				data: &models.Data{
					DataContent: []byte("$1$212345"),
					Metadata:    nil,
					UpdatedAt:   time.Time{},
					ID:          1,
					UserID:      2,
				},
			},
			mockBehavior: func(m *mocks, args args) {
				mock.ExpectBegin()
				expectedQuery := `update data
					 set data_content = $1, metadata = $2, updated_at = now()
					 where id = $3 and user_id = $4`
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.data.DataContent,
						args.data.Metadata, args.data.ID, args.data.UserID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "InvalidTransactionStart",
			args: args{
				data: &models.Data{
					DataContent: []byte("$1$212345"),
					Metadata:    nil,
					UpdatedAt:   time.Time{},
					ID:          1,
					UserID:      2,
				},
			},
			mockBehavior: func(m *mocks, args args) {
				mock.ExpectBegin().WillReturnError(
					fmt.Errorf("failed to begin transaction"))
			},
			wantErr: true,
			err:     fmt.Errorf("failed to begin transaction"),
		},
		{
			name: "InvalidCommitTransaction",
			args: args{
				data: &models.Data{
					DataContent: []byte("$1$212345"),
					Metadata:    nil,
					UpdatedAt:   time.Time{},
					ID:          1,
					UserID:      2,
				},
			},
			mockBehavior: func(m *mocks, args args) {
				mock.ExpectBegin()
				expectedQuery := `update data
					 set data_content = $1, metadata = $2, updated_at = now()
					 where id = $3 and user_id = $4`
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.data.DataContent,
						args.data.Metadata, args.data.ID, args.data.UserID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit().WillReturnError(fmt.Errorf("failed to commit transaction"))
			},
			wantErr: true,
			err:     fmt.Errorf("failed to commit transaction"),
		},
		{
			name: "DeleteDataError",
			args: args{
				data: &models.Data{
					DataContent: []byte("$1$212345"),
					Metadata:    nil,
					UpdatedAt:   time.Time{},
					ID:          1,
					UserID:      2,
				},
			},
			mockBehavior: func(m *mocks, args args) {
				mock.ExpectBegin()
				expectedQuery := `update data
					 set data_content = $1, metadata = $2, updated_at = now()
					 where id = $3 and user_id = $4`
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.data.DataContent,
						args.data.Metadata, args.data.ID, args.data.UserID).
					WillReturnError(fmt.Errorf("error update update data"))
				mock.ExpectCommit()

			},
			wantErr: true,
			err:     fmt.Errorf("error update update data"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			m := &mocks{}

			m.db = sqlx.NewDb(db, "sqlmock")
			pg := dbAdapter{conn: m.db}

			tt.mockBehavior(m, tt.args)
			err := pg.UpdateData(context.Background(), tt.args.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
