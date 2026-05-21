// Package app configures and runs application.
package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	notification *notification.Usecase
	user         *user.Usecase
	task         *task.Usecase
}

type servers struct {
	rmq  *rmqRPCServer.Server
	nats *natsRPCServer.Server
	grpc *grpcserver.Server
	http *httpserver.Server
}

var errUnsupportedOutboxPublisher = errors.New("unsupported outbox publisher")

const appShutdownTimeout = 10 * time.Second

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
	repos := persistence.NewRepositories(pg)
	userRepo := repos.Users()
	taskRepo := repos.Tasks()
	notificationRepo := repos.Notifications()
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
		s.rmq = initRMQServer(cfg, uc, jwtManager, l)
	}

	if enabled.nats {
		s.nats = initNATSServer(cfg, uc, jwtManager, l)
	}

	if enabled.grpc {
		s.grpc = initGRPCServer(cfg, uc, jwtManager, l)
	}

	if enabled.http {
		s.http = initHTTPServer(cfg, uc, jwtManager, pg, l)
	}

	return s
}

func initRMQServer(cfg *config.Config, uc useCases, jwtManager *jwt.Manager, l logger.Interface) *rmqRPCServer.Server {
	rmqRouter := amqprpc.NewRouter(amqprpc.RouterDeps{
		Notification: uc.notification,
		User:         uc.user,
		Task:         uc.task,
		JWTManager:   jwtManager,
		Logger:       l,
	})

	rmqServer, err := rmqRPCServer.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	}

	return rmqServer
}

func initNATSServer(cfg *config.Config, uc useCases, jwtManager *jwt.Manager, l logger.Interface) *natsRPCServer.Server {
	natsRouter := natsrpc.NewRouter(natsrpc.RouterDeps{
		Notification: uc.notification,
		User:         uc.user,
		Task:         uc.task,
		JWTManager:   jwtManager,
		Logger:       l,
	})

	natsServer, err := natsRPCServer.New(cfg.NATS.URL, cfg.NATS.ServerExchange, natsRouter, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - natsServer - server.New: %w", err))
	}

	return natsServer
}

func initGRPCServer(cfg *config.Config, uc useCases, jwtManager *jwt.Manager, l logger.Interface) *grpcserver.Server {
	grpcServer := grpcserver.New(l,
		grpcserver.Port(cfg.GRPC.Port),
		grpcserver.ServerOptions(pbgrpc.ChainUnaryInterceptor(
			grpcmw.RequestIDInterceptor(),
			grpcmw.AuthInterceptor(jwtManager),
		)),
	)
	grpc.NewRouter(grpcServer.App, grpc.RouterDeps{
		Notification: uc.notification,
		User:         uc.user,
		Task:         uc.task,
		Logger:       l,
	})

	return grpcServer
}

func initHTTPServer(cfg *config.Config, uc useCases, jwtManager *jwt.Manager, pg *postgres.Postgres, l logger.Interface) *httpserver.Server {
	httpServer := httpserver.New(l, httpserver.GinMode(cfg.HTTP.Mode), httpserver.Port(cfg.HTTP.Port))
	restapi.NewRouter(httpServer.App, cfg, restapi.RouterDeps{
		Notification: uc.notification,
		User:         uc.user,
		Task:         uc.task,
		JWTManager:   jwtManager,
		Logger:       l,
	}, restapi.ReadinessCheck(pg.Ping))

	return httpServer
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

func (s *servers) waitForShutdown(ctx context.Context, stopApp context.CancelFunc, l logger.Interface) {
	var err error

	select {
	case <-ctx.Done():
		l.Info("app - Run - shutdown requested")
	case err = <-s.httpNotify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-s.grpcNotify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	case err = <-s.rmqNotify():
		l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	case err = <-s.natsNotify():
		l.Error(fmt.Errorf("app - Run - natsServer.Notify: %w", err))
	}

	stopApp()
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

func startOutboxRelay(ctx context.Context, cfg *config.Config, pg *postgres.Postgres, l logger.Interface) (func() error, error) {
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
	relay.Start(ctx)

	return func() error {
		defer publisher.Shutdown()

		return relay.Shutdown()
	}, nil
}

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	appCtx, stopApp := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopApp()

	shutdownTracing, err := observability.InitTracing(observability.TraceConfig{
		Enabled:     cfg.Trace.Enabled,
		Exporter:    cfg.Trace.Exporter,
		ServiceName: cfg.Trace.ServiceName,
	})
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - observability.InitTracing: %w", err))
	}

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(appCtx), appShutdownTimeout)
		defer cancel()

		if err := shutdownTracing(shutdownCtx); err != nil {
			l.Error(fmt.Errorf("app - Run - shutdownTracing: %w", err))
		}
	}()

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax), postgres.Logger(l))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	shutdownOutbox, err := startOutboxRelay(appCtx, cfg, pg, l)
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
	s.waitForShutdown(appCtx, stopApp, l)
}
