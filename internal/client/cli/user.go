package cli

import (
	"bufio"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Sofja96/GophKeeper.git/internal/client/encryption"
	"github.com/Sofja96/GophKeeper.git/internal/client/grpcclient"
)

// LoginCmd возвращает команду CLI для входа
func LoginCmd(client *grpcclient.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate and receive a token",
		Run: func(cmd *cobra.Command, _ []string) {

			reader := bufio.NewReader(os.Stdin)

			cmd.Print("Enter username: ")
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)

			cmd.Print("Enter password: ")
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)

			token, err := client.Login(username, password)
			if err != nil {
				client.Logger.Error("Login failed: %v", err)
			}

			salt := []byte(username)
			encryptionKey := encryption.GenerateEncryptionKey(password, salt)

			client.SetMasterKey(encryptionKey)
			client.SetToken(token)

			if err := client.SyncData(); err != nil {
				cmd.Println("Ошибка синхронизации данных:", err)
				return
			}

		},
	}
}

// RegisterCmd возвращает команду CLI для регистрации
func RegisterCmd(client *grpcclient.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "register",
		Short: "Register a new user",
		Run: func(cmd *cobra.Command, _ []string) {

			reader := bufio.NewReader(os.Stdin)

			cmd.Print("Enter username: ")
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)

			cmd.Print("Enter password: ")
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)

			err := client.Register(username, password)
			if err != nil {
				client.Logger.Error("Registration failed: %v", err)
			}

			cmd.Println("Registration successful!")
		},
	}
}
