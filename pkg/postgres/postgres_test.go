package postgres

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSafeIntToInt32ClampsInvalidValues(t *testing.T) {
	t.Parallel()

	require.Equal(t, int32(1), safeIntToInt32(-10))
	require.Equal(t, int32(1), safeIntToInt32(0))
	require.Equal(t, int32(math.MaxInt32), safeIntToInt32(math.MaxInt64))
	require.Equal(t, int32(32), safeIntToInt32(32))
}
