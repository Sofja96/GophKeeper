package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Sofja96/GophKeeper.git/internal/client/grpcclient"
)

// LoginCmd возвращает команду CLI для входа
func LoginCmd(client *grpcclient.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate and receive a token",
		Run: func(cmd *cobra.Command, args []string) {
			//username := "testuser"
			//password := "password123"

			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter username: ")
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)

			fmt.Print("Enter password: ")
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)

			token, err := client.Login(username, password)
			if err != nil {
				log.Fatalf("Login failed: %v", err)
			}

			fmt.Println("Token:", token)
		},
	}
}

// RegisterCmd возвращает команду CLI для регистрации
func RegisterCmd(client *grpcclient.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "register",
		Short: "Register a new user",
		Run: func(cmd *cobra.Command, args []string) {

			//username := "testuser2"
			//password := "password123"

			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter username: ")
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)

			fmt.Print("Enter password: ")
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)

			err := client.Register(username, password)
			if err != nil {
				log.Fatalf("Registration failed: %v", err)
			}

			fmt.Println("Registration successful!")
		},
	}
}
