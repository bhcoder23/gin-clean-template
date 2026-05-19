package config_test

import (
	"testing"
	"time"

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
	require.Equal(t, "release", cfg.HTTP.Mode)
	require.False(t, cfg.GRPC.Enabled)
	require.False(t, cfg.RMQ.Enabled)
	require.False(t, cfg.NATS.Enabled)
	require.Equal(t, 5*time.Second, cfg.Outbox.PublishTimeout)
}

func TestNewConfigRejectsNoEnabledTransports(t *testing.T) {
	t.Setenv("APP_NAME", "gin-clean-template")
	t.Setenv("APP_VERSION", "1.0.0")
	t.Setenv("HTTP_ENABLED", "false")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("PG_POOL_MAX", "2")
	t.Setenv("PG_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("JWT_SECRET", "secret")

	cfg, err := config.NewConfig()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, "at least one transport must be enabled")
}

func TestNewConfigRejectsUnsafeProductionSwagger(t *testing.T) {
	t.Setenv("APP_NAME", "gin-clean-template")
	t.Setenv("APP_VERSION", "1.0.0")
	t.Setenv("APP_ENV", "production")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("PG_POOL_MAX", "2")
	t.Setenv("PG_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("SWAGGER_ENABLED", "true")

	cfg, err := config.NewConfig()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, "swagger must be disabled in production")
}
