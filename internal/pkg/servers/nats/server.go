package nats

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	nats_server "github.com/nats-io/nats-server/v2/server"

	"boilerplate/internal/model"
)

type server struct {
	name            string
	withJetStream   bool
	jetStreamDomain string
	server          *nats_server.Server
	debug           bool
}

func NewServer(config *model.ConfigNats, opts ...Option) (model.BrokerServer, error) {
	s := &server{
		name: "nats-server",
	}
	for _, opt := range opts {
		opt(s)
	}

	port, err := strconv.Atoi(config.Port)
	if err != nil {
		return nil, fmt.Errorf("nats parse conf port: %w", err)
	}

	httpPort, err := strconv.Atoi(config.HTTPPort)
	if err != nil {
		return nil, fmt.Errorf("nats parse conf http port: %w", err)
	}

	serverOpts := &nats_server.Options{
		Host:            config.Host,
		Port:            port,
		ServerName:      s.name,
		DontListen:      false,
		JetStream:       s.withJetStream,
		JetStreamDomain: s.jetStreamDomain,
		HTTPHost:        config.Host,
		HTTPPort:        httpPort,
		StoreDir:        config.DataDir,
		NoSigs:          true,
		Debug:           s.debug,
	}

	s.server, err = nats_server.NewServer(serverOpts)
	if err != nil {
		return nil, err
	}

	s.server.ConfigureLogger()

	return s, nil
}

func (s *server) Start() error {
	go s.server.Start()
	if !s.server.ReadyForConnections(5 * time.Second) {
		return errors.New("nats server start timeout")
	}

	return nil
}

func (s *server) Stop() error {
	if !s.server.Running() {
		return nil
	}
	s.server.Shutdown()
	return nil
}

func (s *server) GetConn() (net.Conn, error) {
	return s.server.InProcessConn()
}
