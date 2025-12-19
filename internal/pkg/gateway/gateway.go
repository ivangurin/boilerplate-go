package gateway

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"boilerplate/internal/model"
	"boilerplate/internal/pkg/swagger"
)

func Setup(ctx context.Context, configAPI model.ConfigAPI, grpcHandlers []model.GRPCHandler) (*http.ServeMux, error) {
	// Create gRPC client connection
	grpcConn, err := grpc.NewClient(
		net.JoinHostPort(configAPI.Host, configAPI.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("create grpc client: %w", err)
	}

	// Setup main router
	router := http.NewServeMux()

	// Register swagger UI
	swagger.Register(router)

	// Setup gRPC gateway
	gwRouter := runtime.NewServeMux(
		// Добавляем аннотатор metadata для прокидывания cookies из HTTP в gRPC
		runtime.WithMetadata(WithMetadata),
		// Добавляем forwarder для прокидывания cookies из gRPC в HTTP
		runtime.WithForwardResponseOption(WithForwardResponseOption),
	)
	// Register gRPC handlers
	for _, handler := range grpcHandlers {
		if err := handler.RegisterHTTPHandler(ctx, gwRouter, grpcConn); err != nil {
			return nil, fmt.Errorf("register grpc handler: %w", err)
		}
	}

	// Mount gateway to main router
	router.Handle("/", gwRouter)

	return router, nil
}
