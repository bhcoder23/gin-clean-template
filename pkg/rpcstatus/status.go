package rpcstatus

import "errors"

const (
	Success            = "success"
	CodeInternalServer = "INTERNAL_SERVER_ERROR"
	CodeBadHandler     = "BAD_HANDLER"
	HeaderErrorMessage = "Error-Message"
)

type Error struct {
	Code    string
	Message string
}

type CodedError interface {
	RPCCode() string
	RPCMessage() string
}

func FromStatus(status, message, defaultCode, defaultMessage string) error {
	if status == Success {
		return nil
	}

	if status == "" {
		status = defaultCode
		message = defaultMessage
	}

	if message == "" {
		message = status
	}

	return Error{Code: status, Message: message}
}

func FromError(err error, defaultCode, defaultMessage string) Error {
	var coded CodedError
	if errors.As(err, &coded) {
		return Error{
			Code:    coded.RPCCode(),
			Message: coded.RPCMessage(),
		}
	}

	return Error{Code: defaultCode, Message: defaultMessage}
}

func Code(err error) string {
	var rpcErr Error
	if errors.As(err, &rpcErr) {
		return rpcErr.Code
	}

	return ""
}

func Message(err error) string {
	var rpcErr Error
	if errors.As(err, &rpcErr) {
		return rpcErr.Message
	}

	return ""
}

func (e Error) Error() string {
	return e.Message
}
