package app

import (
	"testing"

	"github.com/bhcoder23/gin-clean-template/config"
	"github.com/stretchr/testify/require"
)

func TestEnabledTransports(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		HTTP: config.HTTP{Enabled: true},
		GRPC: config.GRPC{Enabled: true},
		RMQ:  config.RMQ{Enabled: false},
		NATS: config.NATS{Enabled: true},
	}

	enabled := enabledTransports(cfg)

	require.True(t, enabled.http)
	require.True(t, enabled.grpc)
	require.False(t, enabled.rmq)
	require.True(t, enabled.nats)
}
