package db

import (
	"context"
	"database/sql"
	_ "database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/Sofja96/GophKeeper.git/internal/models"
	dmock "github.com/Sofja96/GophKeeper.git/internal/server/storage/db/mocks"
)

type mocks struct {
	db      *sqlx.DB
	storage *dmock.MockAdapter
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	defer db.Close()

	type (
		args struct {
			user *models.User
		}
		mockBehavior func(m *mocks, args args)
	)
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expectedUser *models.User
		wantErr      bool
		err          error
	}{
		{
			name: "SuccessfulCreateUser",
			args: args{
				user: &models.User{
					Username: "testuser",
					Password: "password123",
				},
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `insert into users (
                   username, password) values ($1, $2) on conflict(username) do update set
                   password = EXCLUDED.password;`
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).WithArgs(args.user.Username, args.user.Password).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedUser: &models.User{
				Username: "testuser",
				Password: "password123",
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "ErrorCreateUser",
			args: args{
				user: &models.User{
					Username: "testuser",
					Password: "password123",
				},
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `insert into users (
                   username, password) values ($1, $2) on conflict(username) do update set
                   password = EXCLUDED.password;`
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).WithArgs(args.user.Username, args.user.Password).
					WillReturnError(fmt.Errorf("failed to create user"))
			},
			expectedUser: nil,
			wantErr:      true,
			err:          fmt.Errorf("failed to create user"),
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
			returnedUser, err := pg.CreateUser(context.Background(), tt.args.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, returnedUser, "The returned user does not match the expected user")
			}

		})

	}
}

func TestGetUserIDByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	defer db.Close()

	type (
		args struct {
			username string
		}
		mockBehavior func(m *mocks, args args)
	)
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expectedUser bool
		wantErr      bool
		err          error
	}{
		{
			name: "UserExists",
			args: args{
				username: "testuser",
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `SELECT id FROM users WHERE username = $1`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.username).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("12345"))
			},
			expectedUser: true,
			wantErr:      false,
			err:          nil,
		},
		{
			name: "UserDoesNotExist",
			args: args{
				username: "unknownuser",
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `SELECT id FROM users WHERE username = $1`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.username).
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser: false,
			wantErr:      false,
			err:          nil,
		},
		{
			name: "GetUserIdError",
			args: args{
				username: "erroruser",
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `SELECT id FROM users WHERE username = $1`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.username).
					WillReturnError(fmt.Errorf("error getting id users"))
			},
			expectedUser: false,
			wantErr:      true,
			err:          fmt.Errorf("error getting id users"),
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
			returnedUser, err := pg.GetUserIDByName(context.Background(), tt.args.username)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, returnedUser, "The returned user does not match the expected user")
			}

		})

	}
}

func TestUserHashPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	defer db.Close()

	type (
		args struct {
			username string
		}
		mockBehavior func(m *mocks, args args)
	)
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expectedPass string
		wantErr      bool
		err          error
	}{
		{
			name: "GetPasswordSuccess",
			args: args{
				username: "testuser",
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `SELECT password FROM users WHERE username = $1`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.username).
					WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow("$1$212345"))
			},
			expectedPass: "$1$212345",
			wantErr:      false,
			err:          nil,
		},
		{
			name: "PassNotFound",
			args: args{
				username: "unknownuser",
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `SELECT password FROM users WHERE username = $1`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.username).
					WillReturnError(sql.ErrNoRows)
			},
			expectedPass: "",
			wantErr:      true,
			err:          sql.ErrNoRows,
		},
		{
			name: "GetPasswordError",
			args: args{
				username: "erroruser",
			},
			mockBehavior: func(m *mocks, args args) {
				expectedQuery := `SELECT password FROM users WHERE username = $1`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.username).
					WillReturnError(fmt.Errorf("error getting password on user"))
			},
			expectedPass: "",
			wantErr:      true,
			err:          fmt.Errorf("error getting password on user"),
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
			returnedPass, err := pg.GetUserHashPassword(context.Background(), tt.args.username)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPass, returnedPass, "The returned pass does not match the expected pass")
			}

		})

	}
}
