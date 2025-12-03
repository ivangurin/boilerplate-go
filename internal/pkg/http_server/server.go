package http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const Timeout = 3 * time.Second

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

type server struct {
	server *http.Server
}

func NewServer(port string, handler http.Handler) Server {
	server := &server{
		server: &http.Server{
			Addr:              fmt.Sprintf("127.0.0.1:%s", port),
			ReadHeaderTimeout: Timeout,
			ReadTimeout:       Timeout,
			Handler:           handler,
		},
	}

	return server
}

func (s *server) Start() error {
	err := s.server.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("start http server: %w", err)
		}
	}

	return nil
}

func (s *server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("stop http server: %w", err)
	}

	return nil
}
