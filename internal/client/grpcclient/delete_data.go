package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/Sofja96/GophKeeper.git/internal/client/localstorage"
	"github.com/Sofja96/GophKeeper.git/proto"
)

// DeleteData удаляет данные с указанным ID из локального хранилища и сервера.
// Использует токен для аутентификации при взаимодействии с сервером.
// Возвращает ошибку в случае неудачи.
func (c *Client) DeleteData(dataId int64) error {
	ctx := metadata.AppendToOutgoingContext(context.Background(),
		"authorization", c.GetToken())

	if err := localstorage.DeleteData(c.UserID, dataId); err != nil {
		return fmt.Errorf("ошибка удаления данных из локального хранилища: %w", err)
	}

	req := &proto.DeleteDataRequest{DataId: dataId}

	_, err := c.Client.DeleteData(ctx, req)
	if err != nil {
		return fmt.Errorf("ошибка удаления данных: %w", err)
	}

	fmt.Println("Данные успешно удалены из локального хранилища.")

	return nil
}
