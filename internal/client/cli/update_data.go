package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Sofja96/GophKeeper.git/internal/client/grpcclient"
	"github.com/Sofja96/GophKeeper.git/internal/client/models"
	mdata "github.com/Sofja96/GophKeeper.git/internal/models"
	"github.com/Sofja96/GophKeeper.git/proto"
)

// UpdateDataCmd создает команду Cobra для обновления данных пользователя.
func UpdateDataCmd(client *grpcclient.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-data",
		Short: "Обновить данные",
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Println("\nРежим обновления данных. Введите '8' или 'exit' для выхода.")

			reader := bufio.NewReader(os.Stdin)

			data, err := client.GetData()
			if err != nil {
				cmd.Println("Ошибка получения данных:", err)
				return err
			}

			if len(data) == 0 {
				cmd.Println("У вас нет данных для обновления.")
				return fmt.Errorf("у вас нет данных для обновления")
			}

			cmd.Println("Ваши данные:")
			for _, item := range data {
				cmd.Printf("ID: %d, Тип данных: %s, Метаданные: %v\n", item.ID, item.DataType, item.Metadata)
			}

			cmd.Print("Введите ID документа, который хотите обновить: ")
			var dataID string
			_, err = fmt.Scanln(&dataID)
			if err != nil {
				cmd.Println("Ошибка ввода:", err)
				return fmt.Errorf("ошибка ввода: %w", err)
			}

			id, err := strconv.ParseInt(dataID, 10, 64)
			if err != nil {
				cmd.Println("Ошибка: неверный формат ID")
				return fmt.Errorf("ошибка: неверный формат ID")
			}

			var selectedData *mdata.Data
			for _, item := range data {
				if item.ID == id {
					selectedData = &item
					break
				}
			}

			if selectedData == nil {
				cmd.Println("Ошибка: данные с таким ID не найдены")
				return fmt.Errorf("ошибка: данные с таким ID не найдены")
			}

			var newData models.DataType

			protoDataType, err := mdata.ConvertModelDataTypeToProto(selectedData.DataType)
			if err != nil {
				return err
			}

			switch protoDataType {
			case proto.DataType_LOGIN_PASSWORD:
				newData = inputLoginPassword(reader)
			case proto.DataType_TEXT_DATA:
				newData = inputText(reader)
			case proto.DataType_BANK_CARD:
				newData = inputBankCard(reader)
			case proto.DataType_BINARY_DATA:
				newData = inputBinaryData(reader)
			default:
				cmd.Println("Ошибка: неизвестный тип данных")
				return fmt.Errorf("неверный выбор")
			}

			metadata, err := inputMetadata(reader)
			if err != nil {
				cmd.Println("Ошибка при вводе метаданных:", err)
				return fmt.Errorf("ошибка при вводе метаданных: %v", err)
			}

			req := models.CreateData{
				Data:          newData,
				DataType:      protoDataType,
				EncryptionKey: client.GetMasterKey(),
				Metadata:      metadata,
			}

			err = client.UpdateData(req, id)
			if err != nil {
				cmd.Println("Ошибка при обновлении данных:", err)
				return err
			}

			cmd.Println("Данные успешно обновлены!")
			return nil
		},
	}
}
