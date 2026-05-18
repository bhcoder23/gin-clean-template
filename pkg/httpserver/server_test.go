package httpserver

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestNormalizeGinMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "release default", in: "", want: gin.ReleaseMode},
		{name: "debug", in: "debug", want: gin.DebugMode},
		{name: "test", in: "test", want: gin.TestMode},
		{name: "trimmed", in: " release ", want: gin.ReleaseMode},
		{name: "invalid fallback", in: "prod", want: gin.ReleaseMode},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, normalizeGinMode(tc.in))
		})
	}
}
