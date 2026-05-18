// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bhcoder23/gin-clean-template/config"
	amqprpc "github.com/bhcoder23/gin-clean-template/internal/controller/amqp_rpc"
	"github.com/bhcoder23/gin-clean-template/internal/controller/grpc"
	grpcmw "github.com/bhcoder23/gin-clean-template/internal/controller/grpc/middleware"
	natsrpc "github.com/bhcoder23/gin-clean-template/internal/controller/nats_rpc"
	"github.com/bhcoder23/gin-clean-template/internal/controller/restapi"
	"github.com/bhcoder23/gin-clean-template/internal/repo/persistent"
	"github.com/bhcoder23/gin-clean-template/internal/usecase/notification"
	"github.com/bhcoder23/gin-clean-template/internal/usecase/task"
	"github.com/bhcoder23/gin-clean-template/internal/usecase/user"
	"github.com/bhcoder23/gin-clean-template/pkg/grpcserver"
	"github.com/bhcoder23/gin-clean-template/pkg/httpserver"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	natsRPCServer "github.com/bhcoder23/gin-clean-template/pkg/nats/nats_rpc/server"
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
	userRepo := persistent.NewUserRepo(pg)
	taskRepo := persistent.NewTaskRepo(pg)
	notificationRepo := persistent.NewNotificationRepo(pg)

	return useCases{
		user:         user.New(userRepo, jwtManager),
		task:         task.New(taskRepo, notificationRepo),
		notification: notification.New(notificationRepo),
	}
}

func initServers(cfg *config.Config, uc useCases, jwtManager *jwt.Manager, l logger.Interface) servers {
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
			grpcserver.ServerOptions(pbgrpc.UnaryInterceptor(grpcmw.AuthInterceptor(jwtManager))),
		)
		grpc.NewRouter(grpcServer.App, uc.notification, uc.user, uc.task, l)

		s.grpc = grpcServer
	}

	if enabled.http {
		httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
		restapi.NewRouter(httpServer.App, cfg, uc.notification, uc.user, uc.task, jwtManager, l)

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

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// JWT
	jwtManager := jwt.New(cfg.JWT.Secret, cfg.JWT.TokenExpiry)

	uc := initUseCases(pg, jwtManager)
	s := initServers(cfg, uc, jwtManager, l)
	s.startServers()
	s.waitForShutdown(l)
}
