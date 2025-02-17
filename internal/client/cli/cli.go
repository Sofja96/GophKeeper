package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Sofja96/GophKeeper.git/internal/client/grpcclient"
	"github.com/Sofja96/GophKeeper.git/shared/buildinfo"
)

// StartCLI инициализирует CLI, принимая gRPC-клиент
func StartCLI(client *grpcclient.Client) error {
	rootCmd := &cobra.Command{
		Use:   "gophkeeper",
		Short: "CLI client for GophKeeper",
	}

	// Добавляем команды, передавая клиент
	rootCmd.AddCommand(LoginCmd(client), RegisterCmd(client), VersionCmd())

	return rootCmd.Execute()
}

// VersionCmd возвращает команду для отображения версии
func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show build version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\nBuild Date: %s\n", buildinfo.Version, buildinfo.BuildDate)
		},
	}
}
