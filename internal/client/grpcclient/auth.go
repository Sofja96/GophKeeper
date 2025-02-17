package grpcclient

import (
	"context"
	"fmt"

	"github.com/Sofja96/GophKeeper.git/proto"
)

func (c *Client) Login(username, password string) (string, error) {
	req := &proto.LoginRequest{
		Username: username,
		Password: password,
	}
	resp, err := c.Client.Login(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}
	return resp.Token, nil
}

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
