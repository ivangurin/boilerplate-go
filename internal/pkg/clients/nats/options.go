package nats

import (
	"net"
)

type Option func(*client)

func WithName(name string) Option {
	return func(c *client) {
		c.name = name
	}
}

func WithUrl(url string) Option {
	return func(c *client) {
		c.url = url
	}
}

func WithConn(conn net.Conn) Option {
	return func(c *client) {
		c.conn = conn
	}
}
