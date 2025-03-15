package grpcclient

import (
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	logging "github.com/Sofja96/GophKeeper.git/internal/server/logger"
	"github.com/Sofja96/GophKeeper.git/internal/server/settings"
	"github.com/Sofja96/GophKeeper.git/proto"
)

// Client представляет клиентское подключение к серверу GophKeeper.
type Client struct {
	conn          *grpc.ClientConn
	Client        proto.GophKeeperClient
	Logger        logging.ILogger
	EncryptionKey []byte
	Token         string
	UserID        int64
}

// NewGRPCClient создает новый клиент для подключения к серверу GophKeeper.
// Настройки подключения передаются через объект settings.
// Возвращает объект Client и ошибку, если подключение не удалось.
func NewGRPCClient(settings *settings.Settings) (*Client, error) {
	logger := logging.New(settings)

	cert, err := tls.LoadX509KeyPair(settings.PathCert, settings.PathKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read cert file: %w", err)
	}

	cred := credentials.NewTLS(&tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	})

	conn, err := grpc.NewClient(settings.Host+":"+settings.Port,
		grpc.WithTransportCredentials(cred))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	client := proto.NewGophKeeperClient(conn)
	return &Client{
		conn:   conn,
		Client: client,
		Logger: logger,
	}, nil
}

// Close закрывает соединение.
func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		c.Logger.Fatal("error closing connection: %v", err)
	}
}

// SetMasterKey Устанавливает мастер-пароль
func (c *Client) SetMasterKey(masterKey []byte) {
	c.EncryptionKey = masterKey
}

// GetMasterKey Получает мастер-пароль
func (c *Client) GetMasterKey() []byte {
	return c.EncryptionKey
}

// SetToken устанавливает токен для аунтентификации
func (c *Client) SetToken(token string) {
	c.Token = token
}

// GetToken получает токен для аутентификации
func (c *Client) GetToken() string {
	return c.Token
}
