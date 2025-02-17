package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/Sofja96/GophKeeper.git/internal/server/app"
	"github.com/Sofja96/GophKeeper.git/internal/server/grpcserver"
)

func main() {
	errorCh := make(chan error)
	defer close(errorCh)

	// Инициализация основного сервера
	srv, err := app.Run()
	if err != nil {
		log.Fatalf("cannot start application: %v", err)
	}
	defer srv.GetDbAdapter().Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	// Запуск gRPC сервера
	go func() {
		errorCh <- grpcserver.Run(ctx, srv)
	}()

	// Ожидание ошибок
	err = <-errorCh
	if err != nil {
		log.Fatalf("application was aborted: %v", err)
	}
}
