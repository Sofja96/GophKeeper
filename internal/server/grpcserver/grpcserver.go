package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/Sofja96/GophKeeper.git/internal/server/app"
	"github.com/Sofja96/GophKeeper.git/internal/server/grpcserver/interceptors"
	logging "github.com/Sofja96/GophKeeper.git/internal/server/logger"
	"github.com/Sofja96/GophKeeper.git/proto"
)

// GRPCServer управляет gRPC сервером.
type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
	logger   logging.ILogger
}

// NewGRPCServer создает новый экземпляр GRPCServer.
func NewGRPCServer(srv app.Server) (*GRPCServer, error) {
	lis, err := net.Listen("tcp", srv.GetSettings().Host+":"+srv.GetSettings().Port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	cred, err := credentials.NewServerTLSFromFile(srv.GetSettings().PathCert, srv.GetSettings().PathKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create credentials: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(cred),
		grpc.ChainUnaryInterceptor(
			interceptors.LoggingInterceptor(srv.GetLogger()),
			interceptors.AuthInterceptor(),
		),
	)

	proto.RegisterGophKeeperServer(grpcServer, NewGophKeeperServer(srv))
	reflection.Register(grpcServer)

	return &GRPCServer{
		server:   grpcServer,
		listener: lis,
		logger:   srv.GetLogger(),
	}, nil
}

// Run запускает gRPC сервер и обрабатывает graceful shutdown.
func Run(ctx context.Context, srv app.Server) error {
	grpcSrv, err := NewGRPCServer(srv)
	if err != nil {
		return fmt.Errorf("failed to create gRPC server: %w", err)
	}
	errorCh := make(chan error, 1)

	go func() {
		srv.GetLogger().Info("gRPC server listening at %v", grpcSrv.listener.Addr())
		if err := grpcSrv.server.Serve(grpcSrv.listener); err != nil {
			errorCh <- fmt.Errorf("gRPC server failed: %w", err)
		}
	}()

	defer func() {
		srv.GetLogger().Info("Initiating graceful shutdown...")
		serverCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		grpcSrv.Stop(serverCtx)
		srv.GetLogger().Info("Gracefully stopped")
	}()

	select {
	case err := <-errorCh:
		srv.GetLogger().Info("Error starting server: %v\n", err)
		return err
	case <-ctx.Done():
		srv.GetLogger().Info("Shutdown signal received, shutting down gracefully...")
		return ctx.Err()
	}
}

// Stop зарквает соединение gRPC сервера gracefully.
func (s *GRPCServer) Stop(ctx context.Context) {
	stopped := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.server.Stop()
		s.logger.Info("gRPC server was forcefully stopped")
	case <-stopped:
		s.logger.Info("gRPC server stopped gracefully")
	}
}
