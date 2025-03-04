package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Sofja96/GophKeeper.git/internal/client/grpcclient"
	"github.com/Sofja96/GophKeeper.git/shared/buildinfo"
)

// StartCLI инициализирует CLI, принимая gRPC-клиент
func StartCLI(client *grpcclient.Client) error {
	rootCmd := &cobra.Command{
		Use:   "gophkeeper",
		Short: "CLI client for GophKeeper",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := InteractiveMode(client)
			if err != nil {
				return err
			}
			return nil
		},
	}

	rootCmd.AddCommand(LoginCmd(client), RegisterCmd(client),
		VersionCmd(), CreateDataCmd(client), GetDataCmd(client), DeleteDataCmd(client), UpdateDataCmd(client))

	return rootCmd.Execute()
}

// InteractiveMode запускает интерактивный режим для работы с клиентом.
// В этом режиме пользователь может выбрать одну из команд для выполнения различных операций,
// таких как логин, регистрация, создание, получение, удаление и обновление данных,
func InteractiveMode(client *grpcclient.Client) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nВыберите команду:")
		fmt.Println("1. Логин")
		fmt.Println("2. Регистрация")
		fmt.Println("3. Создать данные")
		fmt.Println("4. Получить данные")
		fmt.Println("5. Удалить данные")
		fmt.Println("6. Обновить данные")
		fmt.Println("7. Получить информацию о версии и дате сборке клиента")
		fmt.Println("8. Выйти")

		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		dummyCmd := &cobra.Command{}

		switch input {
		case "1":
			LoginCmd(client).Run(dummyCmd, nil)
		case "2":
			RegisterCmd(client).Run(dummyCmd, nil)
		case "3":
			CreateDataCmd(client).RunE(dummyCmd, nil)
		case "4":
			GetDataCmd(client).Run(dummyCmd, nil)
		case "5":
			DeleteDataCmd(client).RunE(dummyCmd, nil)
		case "6":
			UpdateDataCmd(client).RunE(dummyCmd, nil)
		case "7":
			VersionCmd().Run(dummyCmd, nil)
		case "8":
			fmt.Println("Выход из программы.")
			return nil
		default:
			fmt.Println("Неизвестная команда. Пожалуйста, выберите число от 1 до 4.")
		}
	}
}

// VersionCmd возвращает команду для отображения версии
func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show build version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Version: %s\nBuild Date: %s\n", buildinfo.Version, buildinfo.BuildDate)
		},
	}
}
