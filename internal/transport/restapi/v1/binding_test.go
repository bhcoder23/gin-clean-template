package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi/v1/request"
	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi/v1/response"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/bhcoder23/gin-clean-template/pkg/requestid"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestBindJSONRejectsInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequestWithContext(
		requestid.WithContext(context.Background(), "req-bind"),
		http.MethodPost,
		"/",
		strings.NewReader(`{`),
	)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	r := &V1{
		l: logger.New("error"),
		v: validator.New(validator.WithRequiredStructEnabled()),
	}

	var body request.RegisterReq
	ok := r.bindJSON(ctx, &body, "test - bind")
	require.False(t, ok)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

	var got response.Error
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &got))
	require.Equal(t, apperror.CodeInvalidRequest, got.Error.Code)
	require.Equal(t, "invalid request body", got.Error.Message)
	require.Equal(t, "req-bind", got.Error.RequestID)
}

func TestBindJSONRejectsValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequestWithContext(
		requestid.WithContext(context.Background(), "req-validate"),
		http.MethodPost,
		"/",
		strings.NewReader(`{"username":"ab","email":"invalid","password":"123"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	r := &V1{
		l: logger.New("error"),
		v: validator.New(validator.WithRequiredStructEnabled()),
	}

	var body request.RegisterReq
	ok := r.bindJSON(ctx, &body, "test - validate")
	require.False(t, ok)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

	var got response.Error
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &got))
	require.Equal(t, apperror.CodeInvalidRequest, got.Error.Code)
	require.Equal(t, "invalid request body", got.Error.Message)
	require.Equal(t, "req-validate", got.Error.RequestID)
}
