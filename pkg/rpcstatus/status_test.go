package rpcstatus

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var errUnexpected = errors.New("unexpected")

func TestFromStatus(t *testing.T) {
	t.Parallel()

	require.Nil(t, FromStatus(Success, "", CodeInternalServer, "internal server error"))

	err := FromStatus("TASK_NOT_FOUND", "task not found", CodeInternalServer, "internal server error")
	require.EqualError(t, err, "task not found")
	require.Equal(t, "TASK_NOT_FOUND", Code(err))
	require.Equal(t, "task not found", Message(err))

	err = FromStatus("UNKNOWN_CODE", "", CodeInternalServer, "internal server error")
	require.EqualError(t, err, "UNKNOWN_CODE")
	require.Equal(t, "UNKNOWN_CODE", Code(err))
	require.Equal(t, "UNKNOWN_CODE", Message(err))
}

func TestFromError(t *testing.T) {
	t.Parallel()

	err := FromError(codedErrorStub{code: "UNAUTHORIZED", message: "unauthorized"}, CodeInternalServer, "internal server error")
	require.Equal(t, "UNAUTHORIZED", err.Code)
	require.Equal(t, "unauthorized", err.Message)

	err = FromError(errUnexpected, CodeInternalServer, "internal server error")
	require.Equal(t, CodeInternalServer, err.Code)
	require.Equal(t, "internal server error", err.Message)
}

type codedErrorStub struct {
	code    string
	message string
}

func (e codedErrorStub) Error() string { return e.message }

func (e codedErrorStub) RPCCode() string { return e.code }

func (e codedErrorStub) RPCMessage() string { return e.message }
