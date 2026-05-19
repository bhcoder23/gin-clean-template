package client

import (
	"time"

	"github.com/bhcoder23/gin-clean-template/pkg/logger"
)

// Option -.
type Option func(*Client)

// Timeout -.
func Timeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// ConnWaitTime -.
func ConnWaitTime(timeout time.Duration) Option {
	return func(c *Client) {
		c.conn.WaitTime = timeout
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *Client) {
		c.conn.Attempts = attempts
	}
}

// Logger configures optional structured logs for connection retries and acknowledgements.
func Logger(l logger.Interface) Option {
	return func(c *Client) {
		c.logger = l
		c.conn.Logger = l
	}
}
