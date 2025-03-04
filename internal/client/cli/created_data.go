package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/Sofja96/GophKeeper.git/internal/client/grpcclient"
	"github.com/Sofja96/GophKeeper.git/internal/client/models"
	"github.com/Sofja96/GophKeeper.git/proto"
)

// CreateDataCmd создаёт CLI команду для создания данных
func CreateDataCmd(client *grpcclient.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "create-data",
		Short: "Создать новые данные",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := bufio.NewReader(os.Stdin)

			cmd.Println("\nРежим создания данных. Введите '8' или 'exit' для выхода.")
			cmd.Println("Выберите тип данных:")
			cmd.Println("1. Логин/Пароль")
			cmd.Println("2. Текстовые данные")
			cmd.Println("3. Путь к файлу")
			cmd.Println("4. Банковская карта")

			choiceStr, _ := reader.ReadString('\n')
			choiceStr = strings.TrimSpace(choiceStr)

			if choiceStr == "" {
				cmd.Println("Ошибка ввода. Введите число от 1 до 4.")
				return fmt.Errorf("ошибка ввода. Введите число от 1 до 4")
			}

			choice, err := strconv.Atoi(choiceStr)
			if err != nil {
				cmd.Println("Ошибка ввода. Введите число от 1 до 4.")
				return fmt.Errorf("ошибка ввода. Введите число от 1 до 4")
			}

			var data models.DataType
			var dataType proto.DataType

			switch choice {
			case 1:
				data = inputLoginPassword(reader)
				dataType = proto.DataType_LOGIN_PASSWORD
			case 2:
				data = inputText(reader)
				dataType = proto.DataType_TEXT_DATA
			case 3:
				data = inputBinaryData(reader)
				dataType = proto.DataType_BINARY_DATA
			case 4:
				data = inputBankCard(reader)
				dataType = proto.DataType_BANK_CARD
			default:
				cmd.Println("Неверный выбор")
				return fmt.Errorf("неверный выбор")
			}

			metadata, err := inputMetadata(reader)
			if err != nil {
				cmd.Println("Ошибка при вводе метаданных:", err)
				return fmt.Errorf("ошибка при вводе метаданных: %v", err)
			}

			req := models.CreateData{
				Data:          data,
				DataType:      dataType,
				EncryptionKey: client.GetMasterKey(),
				Metadata:      metadata,
			}

			dataID, err := client.CreateData(req)
			if err != nil {
				cmd.Println("Ошибка:", err)
				return fmt.Errorf("ошибка: %v", err)
			}

			cmd.Printf("Данные успешно сохранены с ID: %d\n", dataID)
			if err := client.SyncData(); err != nil {
				cmd.Println("Ошибка синхронизации данных:", err)
				return fmt.Errorf("ошибка синхронизации данных: %v", err)
			}
			return nil
		},
	}
}

// inputLoginPassword запрашивает логин и пароль
func inputLoginPassword(reader *bufio.Reader) *models.LoginPasswordType {
	fmt.Print("Введите логин: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Введите пароль: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	return &models.LoginPasswordType{
		Username: username,
		Password: password,
	}
}

// inputBankCard запрашивает данные банковской карты
func inputBankCard(reader *bufio.Reader) *models.BankCardType {
	fmt.Print("Введите номер карты: ")
	cardNumber, _ := reader.ReadString('\n')
	cardNumber = strings.TrimSpace(cardNumber)

	fmt.Print("Введите срок действия (MM/YY): ")
	expiryDate, _ := reader.ReadString('\n')
	expiryDate = strings.TrimSpace(expiryDate)

	fmt.Print("Введите CVV: ")
	cvv, _ := reader.ReadString('\n')
	cvv = strings.TrimSpace(cvv)

	fmt.Print("Введите имя владельца: ")
	holderName, _ := reader.ReadString('\n')
	holderName = strings.TrimSpace(holderName)

	return &models.BankCardType{
		CardNumber: cardNumber,
		ExpiryDate: expiryDate,
		CVV:        cvv,
		HolderName: holderName,
	}
}

// inputText запрашивает произвольные текстовые данные
func inputText(reader *bufio.Reader) *models.TextDataType {
	fmt.Print("Введите текст: ")
	textData, _ := reader.ReadString('\n')
	textData = strings.TrimSpace(textData)

	return &models.TextDataType{
		Text: textData,
	}
}

// inputBinaryData запрашивает бинарные данные
func inputBinaryData(reader *bufio.Reader) *models.BinaryDataType {
	fmt.Print("Введите путь к файлу: ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	binaryData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Ошибка чтения файла", err)
		return nil
	}

	fmt.Print("Введите название файла: ")
	fileName, _ := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)

	return &models.BinaryDataType{
		FilePath: filePath,
		Content:  binaryData,
		Filename: fileName,
	}
}

// inputMetadata запрашивает у пользователя метаданные как список ключ-значение
func inputMetadata(reader *bufio.Reader) (*structpb.Struct, error) {
	metadata := make(map[string]interface{})

	fmt.Println("Введите метаданные (в формате ключ=значение). Оставьте пустую строку для завершения:")

	for {
		fmt.Print("Метаданные: ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if line == "" {
			break
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			fmt.Println("Неверный формат. Используйте ключ=значение.")
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		metadata[key] = value
	}

	structMetadata, err := structpb.NewStruct(metadata)
	if err != nil {
		return nil, fmt.Errorf("ошибка преобразования метаданных: %w", err)
	}

	return structMetadata, nil
}
