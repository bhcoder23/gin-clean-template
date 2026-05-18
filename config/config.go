package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	// Config -.
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		PG      PG
		GRPC    GRPC
		RMQ     RMQ
		NATS    NATS
		JWT     JWT
		Metrics Metrics
		Swagger Swagger
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	HTTP struct {
		Enabled bool   `env:"HTTP_ENABLED" envDefault:"true"`
		Port    string `env:"HTTP_PORT"    envDefault:"8080"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	// GRPC -.
	GRPC struct {
		Enabled bool   `env:"GRPC_ENABLED" envDefault:"false"`
		Port    string `env:"GRPC_PORT"    envDefault:"8081"`
	}

	// RMQ -.
	RMQ struct {
		Enabled        bool   `env:"RMQ_ENABLED"    envDefault:"false"`
		ServerExchange string `env:"RMQ_RPC_SERVER" envDefault:"rpc_server"`
		ClientExchange string `env:"RMQ_RPC_CLIENT" envDefault:"rpc_client"`
		URL            string `env:"RMQ_URL"        envDefault:"amqp://guest:guest@localhost:5672/"`
	}

	// NATS -.
	NATS struct {
		Enabled        bool   `env:"NATS_ENABLED"    envDefault:"false"`
		ServerExchange string `env:"NATS_RPC_SERVER" envDefault:"rpc_server"`
		URL            string `env:"NATS_URL"        envDefault:"nats://guest:guest@localhost:4222/"`
	}

	// JWT -.
	JWT struct {
		Secret      string        `env:"JWT_SECRET,required"`
		TokenExpiry time.Duration `env:"JWT_TOKEN_EXPIRY"    envDefault:"24h"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
