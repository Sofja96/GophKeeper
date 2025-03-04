package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Sofja96/GophKeeper.git/internal/client/grpcclient"
)

// GetDataCmd создает команду для получения всех данных пользователя.
func GetDataCmd(client *grpcclient.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "get-data",
		Short: "Получить все данные пользователя",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("\nРежим получения данных. Введите '8' или 'exit' для выхода.")

			data, err := client.GetData()
			if err != nil {
				cmd.Println("Ошибка получения данных:", err)
				return
			}

			for _, item := range data {
				cmd.Println("ID:", item.ID)
				cmd.Println("Тип данных:", item.DataType)
				cmd.Println("Содержимое:", string(item.DataContent))
				cmd.Println("Метаданные:", item.Metadata)
				cmd.Println("Обновлено:", item.UpdatedAt)
				cmd.Println("---")
			}
		},
	}
}
