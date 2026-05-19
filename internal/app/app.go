// Package app configures and runs application.
package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bhcoder23/gin-clean-template/config"
	"github.com/bhcoder23/gin-clean-template/internal/infra/outbox"
	"github.com/bhcoder23/gin-clean-template/internal/infra/persistence"
	amqprpc "github.com/bhcoder23/gin-clean-template/internal/transport/amqp_rpc"
	"github.com/bhcoder23/gin-clean-template/internal/transport/grpc"
	grpcmw "github.com/bhcoder23/gin-clean-template/internal/transport/grpc/middleware"
	natsrpc "github.com/bhcoder23/gin-clean-template/internal/transport/nats_rpc"
	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi"
	"github.com/bhcoder23/gin-clean-template/internal/usecase/notification"
	"github.com/bhcoder23/gin-clean-template/internal/usecase/task"
	"github.com/bhcoder23/gin-clean-template/internal/usecase/user"
	"github.com/bhcoder23/gin-clean-template/pkg/grpcserver"
	"github.com/bhcoder23/gin-clean-template/pkg/httpserver"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	natsRPCServer "github.com/bhcoder23/gin-clean-template/pkg/nats/nats_rpc/server"
	"github.com/bhcoder23/gin-clean-template/pkg/observability"
	"github.com/bhcoder23/gin-clean-template/pkg/postgres"
	rmqRPCServer "github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc/server"
	pbgrpc "google.golang.org/grpc"
)

type useCases struct {
	notification *notification.UseCase
	user         *user.UseCase
	task         *task.UseCase
}

type servers struct {
	rmq  *rmqRPCServer.Server
	nats *natsRPCServer.Server
	grpc *grpcserver.Server
	http *httpserver.Server
}

var errUnsupportedOutboxPublisher = errors.New("unsupported outbox publisher")

type transportSet struct {
	http bool
	grpc bool
	rmq  bool
	nats bool
}

func enabledTransports(cfg *config.Config) transportSet {
	return transportSet{
		http: cfg.HTTP.Enabled,
		grpc: cfg.GRPC.Enabled,
		rmq:  cfg.RMQ.Enabled,
		nats: cfg.NATS.Enabled,
	}
}

func (t transportSet) any() bool {
	return t.http || t.grpc || t.rmq || t.nats
}

func initUseCases(pg *postgres.Postgres, jwtManager *jwt.Manager) useCases {
	stores := persistence.NewStores(pg)
	userRepo := stores.Users()
	taskRepo := stores.Tasks()
	notificationRepo := stores.Notifications()
	transactor := persistence.NewTransactor(pg)

	return useCases{
		user:         user.New(userRepo, jwtManager),
		task:         task.New(taskRepo, notificationRepo, transactor),
		notification: notification.New(notificationRepo),
	}
}

func initServers(cfg *config.Config, uc useCases, jwtManager *jwt.Manager, pg *postgres.Postgres, l logger.Interface) servers {
	enabled := enabledTransports(cfg)
	if !enabled.any() {
		l.Fatal("app - Run - initServers: at least one transport must be enabled")
	}

	var s servers

	if enabled.rmq {
		rmqRouter := amqprpc.NewRouter(uc.notification, uc.user, uc.task, jwtManager, l)

		rmqServer, err := rmqRPCServer.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
		if err != nil {
			l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
		}

		s.rmq = rmqServer
	}

	if enabled.nats {
		natsRouter := natsrpc.NewRouter(uc.notification, uc.user, uc.task, jwtManager, l)

		natsServer, err := natsRPCServer.New(cfg.NATS.URL, cfg.NATS.ServerExchange, natsRouter, l)
		if err != nil {
			l.Fatal(fmt.Errorf("app - Run - natsServer - server.New: %w", err))
		}

		s.nats = natsServer
	}

	if enabled.grpc {
		grpcServer := grpcserver.New(l,
			grpcserver.Port(cfg.GRPC.Port),
			grpcserver.ServerOptions(pbgrpc.ChainUnaryInterceptor(
				grpcmw.RequestIDInterceptor(),
				grpcmw.AuthInterceptor(jwtManager),
			)),
		)
		grpc.NewRouter(grpcServer.App, uc.notification, uc.user, uc.task, l)

		s.grpc = grpcServer
	}

	if enabled.http {
		httpServer := httpserver.New(l, httpserver.GinMode(cfg.HTTP.Mode), httpserver.Port(cfg.HTTP.Port))
		restapi.NewRouter(httpServer.App, cfg, uc.notification, uc.user, uc.task, jwtManager, l, restapi.ReadinessCheck(pg.Ping))

		s.http = httpServer
	}

	return s
}

func (s *servers) startServers() {
	if s.rmq != nil {
		s.rmq.Start()
	}

	if s.nats != nil {
		s.nats.Start()
	}

	if s.grpc != nil {
		s.grpc.Start()
	}

	if s.http != nil {
		s.http.Start()
	}
}

func (s *servers) httpNotify() <-chan error {
	if s.http == nil {
		return nil
	}

	return s.http.Notify()
}

func (s *servers) grpcNotify() <-chan error {
	if s.grpc == nil {
		return nil
	}

	return s.grpc.Notify()
}

func (s *servers) rmqNotify() <-chan error {
	if s.rmq == nil {
		return nil
	}

	return s.rmq.Notify()
}

func (s *servers) natsNotify() <-chan error {
	if s.nats == nil {
		return nil
	}

	return s.nats.Notify()
}

func (s *servers) waitForShutdown(l logger.Interface) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Info("app - Run - signal: %s", sig.String())
	case err = <-s.httpNotify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-s.grpcNotify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	case err = <-s.rmqNotify():
		l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	case err = <-s.natsNotify():
		l.Error(fmt.Errorf("app - Run - natsServer.Notify: %w", err))
	}

	s.shutdownServers(l)
}

func (s *servers) shutdownServers(l logger.Interface) {
	if s.http != nil {
		if err := s.http.Shutdown(); err != nil {
			l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		}
	}

	if s.grpc != nil {
		if err := s.grpc.Shutdown(); err != nil {
			l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
		}
	}

	if s.rmq != nil {
		if err := s.rmq.Shutdown(); err != nil {
			l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
		}
	}

	if s.nats != nil {
		if err := s.nats.Shutdown(); err != nil {
			l.Error(fmt.Errorf("app - Run - natsServer.Shutdown: %w", err))
		}
	}
}

func startOutboxRelay(cfg *config.Config, pg *postgres.Postgres, l logger.Interface) (func() error, error) {
	if !cfg.Outbox.Enabled {
		return func() error { return nil }, nil
	}

	if cfg.Outbox.Publisher != "nats" {
		return nil, fmt.Errorf("%w: %s", errUnsupportedOutboxPublisher, cfg.Outbox.Publisher)
	}

	natsURL := cfg.Outbox.NATSURL
	if natsURL == "" {
		natsURL = cfg.NATS.URL
	}

	publisher, err := outbox.NewNATSPublisher(natsURL, cfg.Outbox.SubjectPrefix)
	if err != nil {
		return nil, fmt.Errorf("app - startOutboxRelay - outbox.NewNATSPublisher: %w", err)
	}

	relay := outbox.NewRelay(outbox.NewStore(pg.Pool), publisher, l, outbox.RelayConfig{
		PollInterval:   cfg.Outbox.PollInterval,
		BatchSize:      cfg.Outbox.BatchSize,
		MaxAttempts:    cfg.Outbox.MaxAttempts,
		LockTimeout:    cfg.Outbox.LockTimeout,
		PublishTimeout: cfg.Outbox.PublishTimeout,
	})
	relay.Start(context.Background())

	return func() error {
		defer publisher.Shutdown()

		return relay.Shutdown()
	}, nil
}

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	shutdownTracing, err := observability.InitTracing(observability.TraceConfig{
		Enabled:     cfg.Trace.Enabled,
		Exporter:    cfg.Trace.Exporter,
		ServiceName: cfg.Trace.ServiceName,
	})
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - observability.InitTracing: %w", err))
	}

	defer func() {
		if err := shutdownTracing(context.Background()); err != nil {
			l.Error(fmt.Errorf("app - Run - shutdownTracing: %w", err))
		}
	}()

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax), postgres.Logger(l))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	shutdownOutbox, err := startOutboxRelay(cfg, pg, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - startOutboxRelay: %w", err))
	}

	defer func() {
		if err := shutdownOutbox(); err != nil {
			l.Error(fmt.Errorf("app - Run - shutdownOutbox: %w", err))
		}
	}()

	// JWT
	jwtManager := jwt.New(cfg.JWT.Secret, cfg.JWT.TokenExpiry)

	uc := initUseCases(pg, jwtManager)
	s := initServers(cfg, uc, jwtManager, pg, l)
	s.startServers()
	s.waitForShutdown(l)
}
