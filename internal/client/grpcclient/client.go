package grpcclient

import (
	"crypto/tls"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/Sofja96/GophKeeper.git/internal/server/settings"
	"github.com/Sofja96/GophKeeper.git/proto"
)

type Client struct {
	conn   *grpc.ClientConn
	Client proto.GophKeeperClient
}

func NewGRPCClient(settings *settings.Settings) (*Client, error) {
	//Загружаем сертификат
	cert, err := tls.LoadX509KeyPair(settings.PathCert, settings.PathKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read cert file: %w", err)
	}

	//Создаем TLS credentials
	cred := credentials.NewTLS(&tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	})

	// Устанавливаем соединение
	conn, err := grpc.NewClient(settings.Host+":"+settings.Port,
		grpc.WithTransportCredentials(cred))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	client := proto.NewGophKeeperClient(conn)
	return &Client{conn: conn, Client: client}, nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		log.Fatalf("error closing connection: %v", err)
	}
}
