package grpcclient

import (
	"context"
	"fmt"

	"github.com/Sofja96/GophKeeper.git/proto"
)

// Login выполняет аутентификацию пользователя с использованием логина и пароля.
// При успешной аутентификации возвращает токен пользователя, который можно использовать для дальнейших запросов.
// Если аутентификация не удалась, возвращает ошибку.
func (c *Client) Login(username, password string) (string, error) {
	req := &proto.LoginRequest{
		Username: username,
		Password: password,
	}
	resp, err := c.Client.Login(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	c.UserID = resp.UserId
	return resp.Token, nil
}

// Register регистрирует нового пользователя с заданным логином и паролем.
// Если регистрация прошла успешно, функция возвращает nil. В случае ошибки возвращается ошибка с описанием причины.
func (c *Client) Register(username, password string) error {
	req := &proto.RegisterRequest{
		Username: username,
		Password: password,
	}
	_, err := c.Client.Register(context.Background(), req)
	if err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}
	return nil
}
