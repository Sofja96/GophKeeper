package interceptors

import (
	"context"
	"encoding/json"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	logging "github.com/Sofja96/GophKeeper.git/internal/server/logger"
)

// LoggingInterceptor - интерцептор для логирования запросов и ответов.
func LoggingInterceptor(log logging.ILogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		log.Info("gRPC method %s called", info.FullMethod)

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			log.Info("Metadata: %v", md)
		}

		resp, err := handler(ctx, req)

		st, _ := status.FromError(err)
		statusCode := st.Code().String()

		duration := time.Since(start)

		var size int
		if resp != nil {
			jsonResp, _ := json.Marshal(resp)
			size = len(jsonResp)
		}

		if err != nil {
			log.Info("gRPC method %s failed with code %s: %s, "+
				"duration: %s", info.FullMethod, st.Code(), st.Message(), duration)
			return resp, err
		}

		log.Info("gRPC method %s completed successfully, duration: %s, "+
			"status: %s, size: %d", info.FullMethod, duration, statusCode, size)

		return resp, err
	}
}
