package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Sofja96/GophKeeper.git/internal/models"
	mockdb "github.com/Sofja96/GophKeeper.git/internal/server/storage/db/mocks"
	"github.com/Sofja96/GophKeeper.git/internal/server/utils"
)

func TestService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockdb.NewMockAdapter(ctrl)
	service := New(mockDB)

	ctx := context.Background()
	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}

	t.Run("successful registration", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).Return(false, nil)
		mockDB.EXPECT().CreateUser(ctx, gomock.Any()).Return(user, nil)

		result, err := service.RegisterUser(ctx, user)
		assert.NoError(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).Return(true, nil)

		_, err := service.RegisterUser(ctx, user)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, utils.ErrUserExists))
	})
	t.Run("error create user", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).Return(false, nil)
		mockDB.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil, errors.New("failed to create user"))

		_, err := service.RegisterUser(ctx, user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create user")
	})
}

func TestService_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockdb.NewMockAdapter(ctrl)
	service := New(mockDB)

	ctx := context.Background()
	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}

	t.Run("successful login", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).Return(true, nil)
		mockDB.EXPECT().GetUserHashPassword(ctx, user.Username).Return("$2a$10$k8sLGTcrvuI36ZsTddy7EOgarUqltq2nlu5qv2ZG1IiZbqzvYAqjG", nil)

		token, err := service.LoginUser(ctx, user)
		assert.NoError(t, err)
		assert.Contains(t, token, "Bearer ")
	})

	t.Run("user not found", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).Return(false, nil)

		_, err := service.LoginUser(ctx, user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "users not found")
	})
	t.Run("error checking user", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).Return(false, fmt.Errorf("error checking existing user"))

		_, err := service.LoginUser(ctx, user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error checking existing user")
	})

	t.Run("error getting password", func(t *testing.T) {
		mockDB.EXPECT().GetUserIDByName(ctx, user.Username).Return(true, nil)
		mockDB.EXPECT().GetUserHashPassword(ctx, user.Username).Return("", fmt.Errorf("error getting password on user"))

		_, err := service.LoginUser(ctx, user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error getting password on user")
	})
}
