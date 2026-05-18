package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi/middleware"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestApp(t *testing.T) (*gin.Engine, *jwt.Manager) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	jwtManager := jwt.New("test-secret", time.Hour)

	app := gin.New()
	app.Use(middleware.Auth(jwtManager))
	app.GET("/test", func(ctx *gin.Context) {
		if _, ok := ctx.Get("userID"); !ok {
			ctx.Status(http.StatusUnauthorized)

			return
		}

		ctx.Status(http.StatusOK)
	})

	return app, jwtManager
}

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()

	app, jwtManager := newTestApp(t)

	validToken, err := jwtManager.GenerateToken("user-id-123")
	require.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "missing header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid format",
			authHeader:     "Basic xxx",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		localTC := tc

		t.Run(localTC.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/test", http.NoBody)
			if localTC.authHeader != "" {
				req.Header.Set("Authorization", localTC.authHeader)
			}

			app.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			assert.Equal(t, localTC.expectedStatus, resp.StatusCode)
		})
	}
}
