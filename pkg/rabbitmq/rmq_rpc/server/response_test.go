package server

import (
	"testing"

	rmqrpc "github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc"
	"github.com/stretchr/testify/require"
)

func TestEncodeResponseReturnsInternalErrorWhenMarshalFails(t *testing.T) {
	t.Parallel()

	body, status, message, err := encodeResponse(make(chan int))

	require.Error(t, err)
	require.Nil(t, body)
	require.Equal(t, rmqrpc.CodeInternalServer, status)
	require.Equal(t, rmqrpc.ErrInternalServer.Error(), message)
}

func TestEncodeResponseReturnsSuccessWhenMarshalSucceeds(t *testing.T) {
	t.Parallel()

	body, status, message, err := encodeResponse(map[string]string{"status": "ok"})

	require.NoError(t, err)
	require.JSONEq(t, `{"status":"ok"}`, string(body))
	require.Equal(t, rmqrpc.Success, status)
	require.Empty(t, message)
}
