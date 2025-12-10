package model

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Handler interface {
	Mount(router *gin.RouterGroup)
}

type HandlerError struct {
	Error string `json:"error"`
}

type HandlerMessage struct {
	Message string `json:"message"`
}

type GRPCHandler interface {
	RegisterGRPCServer(server *grpc.Server)
	RegisterHTTPHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
}
