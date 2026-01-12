package grpc_server

import (
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"boilerplate/internal/model"
)

type Server interface {
	Start() error
	Stop()
}

type server struct {
	addr   string
	server *grpc.Server
}

func NewServer(
	host, port string,
	middleware []grpc.UnaryServerInterceptor,
	handlers []model.GRPCHandler,
) Server {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware...),
	)

	reflection.Register(s)

	for _, handler := range handlers {
		handler.RegisterGRPCServer(s)
	}

	return &server{
		addr:   net.JoinHostPort(host, port),
		server: s,
	}
}

func (s *server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("create listener grpc server: %w", err)
	}

	err = s.server.Serve(listener)
	if err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	return nil
}

func (s *server) Stop() {
	go func() {
		time.Sleep(time.Second * 5)
		s.server.Stop()
	}()
	s.server.GracefulStop()
}
