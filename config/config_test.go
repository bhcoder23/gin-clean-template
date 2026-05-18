package config_test

import (
	"testing"

	"github.com/bhcoder23/gin-clean-template/config"
	"github.com/stretchr/testify/require"
)

func TestNewConfigTransportDefaults(t *testing.T) {
	t.Setenv("APP_NAME", "gin-clean-template")
	t.Setenv("APP_VERSION", "1.0.0")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("PG_POOL_MAX", "2")
	t.Setenv("PG_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("JWT_SECRET", "secret")

	cfg, err := config.NewConfig()

	require.NoError(t, err)
	require.True(t, cfg.HTTP.Enabled)
	require.True(t, cfg.GRPC.Enabled)
	require.True(t, cfg.RMQ.Enabled)
	require.True(t, cfg.NATS.Enabled)
}
