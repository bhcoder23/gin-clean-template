// Package server implements RabbitMQ RPC server.
package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	rmqrpc "github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc"
	"github.com/bhcoder23/gin-clean-template/pkg/requestid"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

const (
	_defaultWaitTime = 5 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second
)

// CallHandler -.
type CallHandler func(context.Context, *amqp.Delivery) (any, error)

// Server -.
type Server struct {
	ctx context.Context
	eg  *errgroup.Group

	conn   *rmqrpc.Connection
	router map[string]CallHandler
	stop   chan struct{}
	notify chan error

	timeout time.Duration

	logger logger.Interface
}

// New -.
func New(url, serverExchange string, router map[string]CallHandler, l logger.Interface, opts ...Option) (*Server, error) {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1) // Run only one goroutine

	cfg := rmqrpc.Config{
		URL:      url,
		WaitTime: _defaultWaitTime,
		Attempts: _defaultAttempts,
		Logger:   l,
	}

	s := &Server{
		ctx:     ctx,
		eg:      group,
		conn:    rmqrpc.New(serverExchange, cfg),
		router:  router,
		stop:    make(chan struct{}),
		notify:  make(chan error, 1),
		timeout: _defaultTimeout,
		logger:  l,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	err := s.conn.AttemptConnect()
	if err != nil {
		return nil, fmt.Errorf("rmq_rpc server - NewServer - s.conn.AttemptConnect: %w", err)
	}

	return s, nil
}

// Start -.
func (s *Server) Start() {
	s.eg.Go(func() error {
		err := s.handleMessages()
		if err != nil {
			s.notify <- err

			close(s.notify)

			return err
		}

		return nil
	})

	s.logger.Info("rmq_rpc server - Server - Started")
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	var shutdownErrors []error

	close(s.stop)

	// Wait for all goroutines to finish and get any error
	err := s.eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(err, "rmq_rpc server - Server - Shutdown - s.eg.Wait")

		shutdownErrors = append(shutdownErrors, err)
	}

	// Close connection

	err = s.conn.Connection.Close()
	if err != nil {
		s.logger.Error(err, "rmq_rpc server - Server - Shutdown - s.Connection.Close")

		shutdownErrors = append(shutdownErrors, err)
	}

	s.logger.Info("rmq_rpc server - Server - Shutdown")

	return errors.Join(shutdownErrors...)
}

func (s *Server) handleMessages() error {
	for {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case <-s.stop:
			return nil
		case d, opened := <-s.conn.Delivery:
			if !opened {
				err := s.reconnect()
				if err != nil {
					return err
				}

				break
			}

			s.serveCall(&d)
		}
	}
}

func (s *Server) reconnect() error {
	return s.conn.AttemptConnect()
}

func (s *Server) serveCall(d *amqp.Delivery) {
	defer s.ack(d, false)

	callHandler, ok := s.router[d.Type]
	if !ok {
		s.publish(d, nil, rmqrpc.CodeBadHandler, rmqrpc.ErrBadHandler.Error())

		return
	}

	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()

	ctx = requestid.WithContext(ctx, requestIDFromDelivery(d))

	response, err := callHandler(ctx, d)
	if err != nil {
		rpcErr := rmqrpc.ErrorFromError(err)

		s.publish(d, nil, rpcErr.Code, rpcErr.Message)

		if rpcErr.Code == rmqrpc.CodeInternalServer {
			s.logger.Error(err, "rmq_rpc server - Server - serveCall - callHandler")

			return
		}

		s.logger.Warn(err, "rmq_rpc server - Server - serveCall - callHandler")

		return
	}

	body, status, message, err := encodeResponse(response)
	if err != nil {
		s.logger.Error(err, "rmq_rpc server - Server - serveCall - json.Marshal")
	}

	s.publish(d, body, status, message)
}

func encodeResponse(response any) (body []byte, status, message string, err error) {
	body, err = json.Marshal(response)
	if err != nil {
		return nil, rmqrpc.CodeInternalServer, rmqrpc.ErrInternalServer.Error(), err
	}

	return body, rmqrpc.Success, "", nil
}

func (s *Server) ack(d *amqp.Delivery, multiple bool) {
	err := d.Ack(multiple)
	if err != nil {
		s.logger.Error(err, "rmq_rpc server - Server - ack - d.Ack")
	}
}

func (s *Server) publish(d *amqp.Delivery, body []byte, status, message string) {
	err := s.conn.Channel.Publish(
		d.ReplyTo,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: d.CorrelationId,
			Headers:       responseHeaders(d, message),
			Type:          status,
			Body:          body,
		},
	)
	if err != nil {
		s.logger.Error(err, "rmq_rpc server - Server - publish - s.conn.Channel.Publish")
	}
}

func responseHeaders(d *amqp.Delivery, message string) amqp.Table {
	headers := requestIDHeaders(d)
	if message != "" {
		headers[rmqrpc.HeaderErrorMessage] = message
	}

	return headers
}

func requestIDFromDelivery(d *amqp.Delivery) string {
	if d.Headers != nil {
		if value, ok := d.Headers[requestid.MetadataKey].(string); ok {
			return requestid.Normalize(value)
		}
	}

	return requestid.Normalize(d.CorrelationId)
}

func requestIDHeaders(d *amqp.Delivery) amqp.Table {
	return amqp.Table{requestid.MetadataKey: requestIDFromDelivery(d)}
}
