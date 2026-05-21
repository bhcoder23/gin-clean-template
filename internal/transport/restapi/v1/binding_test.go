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

func TestBindJSONRejectsBadRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Parallel()

	tests := []struct {
		name      string
		body      string
		requestID string
		logLabel  string
	}{
		{
			name:      "invalid json",
			body:      `{`,
			requestID: "req-bind",
			logLabel:  "test - bind",
		},
		{
			name:      "validation error",
			body:      `{"username":"ab","email":"invalid","password":"123"}`,
			requestID: "req-validate",
			logLabel:  "test - validate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			req := httptest.NewRequestWithContext(
				requestid.WithContext(context.Background(), tt.requestID),
				http.MethodPost,
				"/",
				strings.NewReader(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req

			r := &V1{
				l: logger.New("error"),
				v: validator.New(validator.WithRequiredStructEnabled()),
			}

			var body request.RegisterReq

			ok := r.bindJSON(ctx, &body, tt.logLabel)
			require.False(t, ok)
			require.Equal(t, http.StatusBadRequest, recorder.Code)

			var got response.Error
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &got))
			require.Equal(t, apperror.CodeInvalidRequest, got.Error.Code)
			require.Equal(t, "invalid request body", got.Error.Message)
			require.Equal(t, tt.requestID, got.Error.RequestID)
		})
	}
}
