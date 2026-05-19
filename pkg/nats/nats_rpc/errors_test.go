package natsrpc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var errUnexpected = errors.New("unexpected")

func TestErrorFromStatus(t *testing.T) {
	t.Parallel()

	require.Nil(t, ErrorFromStatus(Success, ""))

	err := ErrorFromStatus("TASK_NOT_FOUND", "task not found")
	require.EqualError(t, err, "task not found")
	require.Equal(t, "TASK_NOT_FOUND", ErrorCode(err))
	require.Equal(t, "task not found", ErrorMessage(err))

	err = ErrorFromStatus("UNKNOWN_CODE", "")
	require.EqualError(t, err, "UNKNOWN_CODE")
	require.Equal(t, "UNKNOWN_CODE", ErrorCode(err))
	require.Equal(t, "UNKNOWN_CODE", ErrorMessage(err))
}

func TestErrorFromError(t *testing.T) {
	t.Parallel()

	err := ErrorFromError(rpcErrorStub{code: "UNAUTHORIZED", message: "unauthorized"})
	require.Equal(t, "UNAUTHORIZED", err.Code)
	require.Equal(t, "unauthorized", err.Message)

	err = ErrorFromError(errUnexpected)
	require.Equal(t, CodeInternalServer, err.Code)
	require.Equal(t, "internal server error", err.Message)
}

type rpcErrorStub struct {
	code    string
	message string
}

func (e rpcErrorStub) Error() string { return e.message }

func (e rpcErrorStub) RPCCode() string { return e.code }

func (e rpcErrorStub) RPCMessage() string { return e.message }
