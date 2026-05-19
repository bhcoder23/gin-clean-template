package natsrpc

import (
	"errors"

	"github.com/bhcoder23/gin-clean-template/pkg/rpcstatus"
)

var (
	// ErrTimeout -.
	ErrTimeout = errors.New("timeout")
	// ErrInternalServer -.
	ErrInternalServer = errors.New("internal server error")
	// ErrBadHandler -.
	ErrBadHandler = errors.New("unregistered handler")
)

const (
	Success            = rpcstatus.Success
	CodeInternalServer = rpcstatus.CodeInternalServer
	CodeBadHandler     = rpcstatus.CodeBadHandler
	HeaderErrorMessage = rpcstatus.HeaderErrorMessage
)

type Error = rpcstatus.Error

func ErrorFromStatus(status, message string) error {
	return rpcstatus.FromStatus(status, message, CodeInternalServer, ErrInternalServer.Error())
}

func ErrorFromError(err error) Error {
	return rpcstatus.FromError(err, CodeInternalServer, ErrInternalServer.Error())
}

func ErrorCode(err error) string {
	return rpcstatus.Code(err)
}

func ErrorMessage(err error) string {
	return rpcstatus.Message(err)
}
