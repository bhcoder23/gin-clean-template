// Package httpserver implements HTTP server.
package httpserver

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
	_defaultGinMode         = gin.ReleaseMode
)

// Server -.
type Server struct {
	eg *errgroup.Group

	App        *gin.Engine
	httpServer *http.Server
	notify     chan error

	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
	mode            string

	logger logger.Interface
}

// New -.
func New(l logger.Interface, opts ...Option) *Server {
	group := new(errgroup.Group)
	group.SetLimit(1)

	s := &Server{
		eg:              group,
		App:             nil,
		httpServer:      nil,
		notify:          make(chan error, 1),
		address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
		mode:            _defaultGinMode,
		logger:          l,
	}

	for _, opt := range opts {
		opt(s)
	}

	mode := normalizeGinMode(s.mode)
	if s.mode != "" && mode != strings.ToLower(s.mode) {
		s.logger.Warn("restapi server - Server - invalid gin mode %q, fallback to %q", s.mode, mode)
	}

	gin.SetMode(mode)

	app := gin.New()
	s.App = app
	s.httpServer = &http.Server{
		Addr:              s.address,
		Handler:           app,
		ReadHeaderTimeout: s.readTimeout,
		ReadTimeout:       s.readTimeout,
		WriteTimeout:      s.writeTimeout,
	}

	return s
}

func normalizeGinMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case gin.DebugMode:
		return gin.DebugMode
	case gin.TestMode:
		return gin.TestMode
	case "", gin.ReleaseMode:
		return gin.ReleaseMode
	default:
		return gin.ReleaseMode
	}
}

// Start -.
func (s *Server) Start() {
	s.eg.Go(func() error {
		err := s.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.notify <- err

			close(s.notify)

			return err
		}

		return nil
	})

	s.logger.Info("restapi server - Server - Started")
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	var shutdownErrors []error

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	err := s.httpServer.Shutdown(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(err, "restapi server - Server - Shutdown - s.httpServer.Shutdown")

		shutdownErrors = append(shutdownErrors, err)
	}

	err = s.eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(err, "restapi server - Server - Shutdown - s.eg.Wait")

		shutdownErrors = append(shutdownErrors, err)
	}

	s.logger.Info("restapi server - Server - Shutdown")

	return errors.Join(shutdownErrors...)
}
