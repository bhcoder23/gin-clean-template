package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

var (
	errNoEnabledTransports      = errors.New("at least one transport must be enabled")
	errInvalidPGPoolMax         = errors.New("PG_POOL_MAX must be greater than zero")
	errProductionSwaggerEnabled = errors.New("swagger must be disabled in production")
	errProductionDefaultJWT     = errors.New("JWT_SECRET must be changed in production")
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
		Trace   Trace
		Outbox  Outbox
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
		Env     string `env:"APP_ENV"              envDefault:"local"`
	}

	// HTTP -.
	HTTP struct {
		Enabled bool   `env:"HTTP_ENABLED" envDefault:"true"`
		Mode    string `env:"GIN_MODE"     envDefault:"release"`
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

	// Trace -.
	Trace struct {
		Enabled     bool   `env:"TRACE_ENABLED"      envDefault:"false"`
		Exporter    string `env:"TRACE_EXPORTER"     envDefault:"stdout"`
		ServiceName string `env:"TRACE_SERVICE_NAME" envDefault:""`
	}

	// Outbox -.
	Outbox struct {
		Enabled        bool          `env:"OUTBOX_ENABLED"        envDefault:"false"`
		Publisher      string        `env:"OUTBOX_PUBLISHER"      envDefault:"nats"`
		NATSURL        string        `env:"OUTBOX_NATS_URL"       envDefault:""`
		SubjectPrefix  string        `env:"OUTBOX_SUBJECT_PREFIX" envDefault:"events"`
		PollInterval   time.Duration `env:"OUTBOX_POLL_INTERVAL"  envDefault:"1s"`
		BatchSize      int           `env:"OUTBOX_BATCH_SIZE"     envDefault:"20"`
		MaxAttempts    int           `env:"OUTBOX_MAX_ATTEMPTS"   envDefault:"10"`
		LockTimeout    time.Duration `env:"OUTBOX_LOCK_TIMEOUT"   envDefault:"5m"`
		PublishTimeout time.Duration `env:"OUTBOX_PUBLISH_TIMEOUT" envDefault:"5s"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	return cfg, nil
}

// Validate checks template-level invariants that env tags cannot express.
func (c *Config) Validate() error {
	var validationErrors []error

	if !c.HTTP.Enabled && !c.GRPC.Enabled && !c.RMQ.Enabled && !c.NATS.Enabled {
		validationErrors = append(validationErrors, errNoEnabledTransports)
	}

	if c.PG.PoolMax <= 0 {
		validationErrors = append(validationErrors, errInvalidPGPoolMax)
	}

	if c.production() {
		if c.Swagger.Enabled {
			validationErrors = append(validationErrors, errProductionSwaggerEnabled)
		}

		if c.JWT.Secret == "your-secret-key-change-in-production" {
			validationErrors = append(validationErrors, errProductionDefaultJWT)
		}
	}

	return errors.Join(validationErrors...)
}

func (c *Config) production() bool {
	return strings.EqualFold(strings.TrimSpace(c.App.Env), "production")
}
