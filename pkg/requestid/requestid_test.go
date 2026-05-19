package requestid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnsureKeepsExistingRequestID(t *testing.T) {
	t.Parallel()

	ctx := WithContext(t.Context(), "request-123")

	nextCtx, id := Ensure(ctx)

	require.Equal(t, "request-123", id)

	stored, ok := FromContext(nextCtx)
	require.True(t, ok)
	require.Equal(t, "request-123", stored)
}

func TestEnsureCreatesRequestID(t *testing.T) {
	t.Parallel()

	ctx, id := Ensure(t.Context())

	require.NotEmpty(t, id)

	stored, ok := FromContext(ctx)
	require.True(t, ok)
	require.Equal(t, id, stored)
}
