package s3

import logger_pkg "boilerplate/internal/pkg/logger"

type option func(*client)

func WithLogger(logger logger_pkg.Logger) option {
	return func(c *client) {
		c.logger = logger
	}
}
