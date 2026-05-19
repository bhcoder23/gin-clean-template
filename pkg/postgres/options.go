package postgres

import (
	"time"

	"github.com/bhcoder23/gin-clean-template/pkg/logger"
)

// Option -.
type Option func(*Postgres)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *Postgres) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.connTimeout = timeout
	}
}

// Logger configures optional structured logs for connection retries.
func Logger(l logger.Interface) Option {
	return func(c *Postgres) {
		c.logger = l
	}
}
