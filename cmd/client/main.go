package main

import (
	"log"

	"github.com/Sofja96/GophKeeper.git/internal/client/cli"
	"github.com/Sofja96/GophKeeper.git/internal/client/grpcclient"
	"github.com/Sofja96/GophKeeper.git/internal/server/settings"
)

func main() {
	conf, err := settings.GetSettings()
	if err != nil {
		log.Fatalf("error load configuration: %v", err)
	}

	client, err := grpcclient.NewGRPCClient(conf)
	if err != nil {
		log.Fatalf("failed to create gRPC client: %v", err)
	}
	defer client.Close()

	err = cli.StartCLI(client)
	if err != nil {
		log.Fatalf("cannot start CLI applictaion :%v", err)
	}
}
