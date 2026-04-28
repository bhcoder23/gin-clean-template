package restapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bhcoder23/gin-clean-template/config"
	"github.com/bhcoder23/gin-clean-template/internal/controller/restapi"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestNewRouterRegistersHealthz(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	app := gin.New()
	cfg := &config.Config{}
	jwtManager := jwt.New("test-secret", time.Hour)
	l := logger.New("error")

	restapi.NewRouter(app, cfg, nil, nil, nil, jwtManager, l)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/healthz", http.NoBody)

	app.ServeHTTP(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
