package restapi_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bhcoder23/gin-clean-template/config"
	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var errListShouldNotRun = errors.New("list should not run")

type taskUseCaseStub struct{}

func (taskUseCaseStub) Create(_ context.Context, _, _, _ string) (domain.Task, error) {
	return domain.Task{}, nil
}

func (taskUseCaseStub) Get(_ context.Context, _, _ string) (domain.Task, error) {
	return domain.Task{}, nil
}

func (taskUseCaseStub) List(_ context.Context, _ string, _ *domain.TaskStatus, _ string, _, _ int) ([]domain.Task, int, error) {
	return []domain.Task{}, 0, nil
}

func (taskUseCaseStub) Update(_ context.Context, _, _, _, _ string) (domain.Task, error) {
	return domain.Task{}, nil
}

func (taskUseCaseStub) Transition(_ context.Context, _, _ string, _ domain.TaskStatus) (domain.Task, error) {
	return domain.Task{}, nil
}

func (taskUseCaseStub) Delete(_ context.Context, _, _ string) error {
	return nil
}

type taskUseCaseWithListError struct{}

func (taskUseCaseWithListError) Create(_ context.Context, _, _, _ string) (domain.Task, error) {
	return domain.Task{}, nil
}

func (taskUseCaseWithListError) Get(_ context.Context, _, _ string) (domain.Task, error) {
	return domain.Task{}, nil
}

func (taskUseCaseWithListError) List(_ context.Context, _ string, _ *domain.TaskStatus, _ string, _, _ int) ([]domain.Task, int, error) {
	return nil, 0, errListShouldNotRun
}

func (taskUseCaseWithListError) Update(_ context.Context, _, _, _, _ string) (domain.Task, error) {
	return domain.Task{}, nil
}

func (taskUseCaseWithListError) Transition(_ context.Context, _, _ string, _ domain.TaskStatus) (domain.Task, error) {
	return domain.Task{}, nil
}

func (taskUseCaseWithListError) Delete(_ context.Context, _, _ string) error {
	return nil
}

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

func TestNewRouterRegistersTaskCollectionWithoutTrailingSlash(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	app := gin.New()
	cfg := &config.Config{}
	jwtManager := jwt.New("test-secret", time.Hour)
	l := logger.New("error")

	restapi.NewRouter(app, cfg, nil, nil, taskUseCaseStub{}, jwtManager, l)

	token, err := jwtManager.GenerateToken("user-id-123")
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/v1/tasks", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+token)

	app.ServeHTTP(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Empty(t, resp.Header.Get("Location"))
}

func TestNewRouterRejectsInvalidTaskPagination(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	app := gin.New()
	cfg := &config.Config{}
	jwtManager := jwt.New("test-secret", time.Hour)
	l := logger.New("error")

	restapi.NewRouter(app, cfg, nil, nil, taskUseCaseWithListError{}, jwtManager, l)

	token, err := jwtManager.GenerateToken("user-id-123")
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/v1/tasks?limit=abc", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+token)

	app.ServeHTTP(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
