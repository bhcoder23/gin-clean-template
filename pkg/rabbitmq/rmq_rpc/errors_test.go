package rmqrpc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorFromStatus(t *testing.T) {
	t.Parallel()

	require.Nil(t, ErrorFromStatus(Success))
	require.Equal(t, ErrUnauthorized, ErrorFromStatus(ErrUnauthorized.Error()))
	require.Equal(t, ErrTaskNotFound, ErrorFromStatus(ErrTaskNotFound.Error()))
	require.Equal(t, ErrNotificationNotFound, ErrorFromStatus(ErrNotificationNotFound.Error()))
	require.Nil(t, ErrorFromStatus("unknown"))
}

func TestIsKnownStatus(t *testing.T) {
	t.Parallel()

	require.True(t, IsKnownStatus(ErrInvalidRequest.Error()))
	require.True(t, IsKnownStatus(ErrInvalidTransition.Error()))
	require.False(t, IsKnownStatus("unknown"))
	require.False(t, IsKnownStatus(Success))
}
