package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Sofja96/GophKeeper.git/internal/client/grpcclient"
)

// DeleteDataCmd  создаёт CLI команду для удаления данных
func DeleteDataCmd(client *grpcclient.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-data",
		Short: "Удалить данные по ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println("\nРежим удаления данных. Введите '8' или 'exit' для выхода.")

			data, err := client.GetData()
			if err != nil {
				cmd.Println("Ошибка получения данных:", err)
				return fmt.Errorf("ошибка получения данных: %w", err)
			}

			if len(data) == 0 {
				cmd.Println("У вас нет данных для удаления.")
				return fmt.Errorf("у вас нет данных для удаления")
			}

			cmd.Println("Ваши данные:")
			for _, item := range data {
				cmd.Printf("ID: %d, Тип данных: %s, Метаданные: %v\n", item.ID, item.DataType, item.Metadata)
			}

			cmd.Print("Введите ID документа, который хотите удалить: ")
			var dataID string
			_, err = fmt.Scanln(&dataID)
			if err != nil {
				cmd.Println("Ошибка ввода:", err)
				return fmt.Errorf("ошибка ввода %w", err)
			}

			id, err := strconv.ParseInt(dataID, 10, 64)
			if err != nil {
				cmd.Println("Ошибка: неверный формат ID")
				return fmt.Errorf("ошибка: неверный формат ID %w", err)
			}

			err = client.DeleteData(id)
			if err != nil {
				cmd.Println("Ошибка удаления данных:", err)
				return fmt.Errorf("ошибка удаления данных: %w", err)
			}

			cmd.Println("Данные успешно удалены")
			return nil
		},
	}
}
