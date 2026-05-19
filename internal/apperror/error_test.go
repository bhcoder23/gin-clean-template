package apperror

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errDatabaseUnavailable = errors.New("database unavailable")
	errUnknown             = errors.New("unknown")
	errPostgresDown        = errors.New("postgres down")
	errBoom                = errors.New("boom")
)

func TestFrom(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		httpStatus int
		grpcCode   codes.Code
		code       string
		message    string
		expected   bool
	}{
		{name: "invalid request", err: ErrInvalidRequest, httpStatus: http.StatusBadRequest, grpcCode: codes.InvalidArgument, code: CodeInvalidRequest, message: "invalid request body", expected: true},
		{name: "unauthorized", err: ErrUnauthorized, httpStatus: http.StatusUnauthorized, grpcCode: codes.Unauthenticated, code: CodeUnauthorized, message: "unauthorized", expected: true},
		{name: "user already exists", err: domain.ErrUserAlreadyExists, httpStatus: http.StatusConflict, grpcCode: codes.AlreadyExists, code: CodeUserAlreadyExists, message: "user already exists", expected: true},
		{name: "invalid credentials", err: domain.ErrInvalidCredentials, httpStatus: http.StatusUnauthorized, grpcCode: codes.Unauthenticated, code: CodeInvalidCredentials, message: "invalid credentials", expected: true},
		{name: "user not found", err: domain.ErrUserNotFound, httpStatus: http.StatusNotFound, grpcCode: codes.NotFound, code: CodeUserNotFound, message: "user not found", expected: true},
		{name: "invalid username", err: domain.ErrInvalidUsername, httpStatus: http.StatusBadRequest, grpcCode: codes.InvalidArgument, code: CodeInvalidUsername, message: "invalid username", expected: true},
		{name: "invalid email", err: domain.ErrInvalidEmail, httpStatus: http.StatusBadRequest, grpcCode: codes.InvalidArgument, code: CodeInvalidEmail, message: "invalid email", expected: true},
		{name: "password too short", err: domain.ErrPasswordTooShort, httpStatus: http.StatusBadRequest, grpcCode: codes.InvalidArgument, code: CodePasswordTooShort, message: "password too short", expected: true},
		{name: "task not found", err: domain.ErrTaskNotFound, httpStatus: http.StatusNotFound, grpcCode: codes.NotFound, code: CodeTaskNotFound, message: "task not found", expected: true},
		{name: "task title required", err: domain.ErrTaskTitleRequired, httpStatus: http.StatusBadRequest, grpcCode: codes.InvalidArgument, code: CodeTaskTitleRequired, message: "task title is required", expected: true},
		{name: "task completed", err: domain.ErrTaskCompleted, httpStatus: http.StatusBadRequest, grpcCode: codes.InvalidArgument, code: CodeTaskCompleted, message: "completed task cannot be modified", expected: true},
		{name: "invalid transition", err: domain.ErrInvalidTransition, httpStatus: http.StatusBadRequest, grpcCode: codes.InvalidArgument, code: CodeInvalidTransition, message: "invalid status transition", expected: true},
		{name: "notification not found", err: domain.ErrNotificationNotFound, httpStatus: http.StatusNotFound, grpcCode: codes.NotFound, code: CodeNotificationNotFound, message: "notification not found", expected: true},
		{name: "internal", err: errDatabaseUnavailable, httpStatus: http.StatusInternalServerError, grpcCode: codes.Internal, code: CodeInternalServer, message: "internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := From(errors.Join(tt.err))

			require.Equal(t, tt.httpStatus, got.HTTPStatus)
			require.Equal(t, tt.grpcCode, got.GRPCCode)
			require.Equal(t, tt.code, got.Code)
			require.Equal(t, tt.message, got.Message)
			require.Equal(t, tt.expected, got.Expected)
		})
	}
}

func TestKnown(t *testing.T) {
	t.Parallel()

	require.True(t, Known(errors.Join(domain.ErrTaskNotFound)))
	require.False(t, Known(nil))
	require.False(t, Known(errUnknown))
}

func TestGRPC(t *testing.T) {
	t.Parallel()

	err := GRPC(domain.ErrTaskNotFound)

	require.Equal(t, codes.NotFound, status.Code(err))
	require.Equal(t, "task not found", status.Convert(err).Message())
}

func TestRPC(t *testing.T) {
	t.Parallel()

	require.NoError(t, RPC(nil))

	taskErr := RPC(domain.ErrTaskNotFound)
	require.EqualError(t, taskErr, "task not found")
	require.Equal(t, CodeTaskNotFound, rpcCode(t, taskErr))
	require.Equal(t, "task not found", rpcMessage(t, taskErr))

	internalErr := RPC(errPostgresDown)
	require.EqualError(t, internalErr, "internal server error")
	require.Equal(t, CodeInternalServer, rpcCode(t, internalErr))
	require.Equal(t, "internal server error", rpcMessage(t, internalErr))
}

func rpcCode(t *testing.T, err error) string {
	t.Helper()

	rpcErr, ok := err.(interface{ RPCCode() string })
	require.True(t, ok)

	return rpcErr.RPCCode()
}

func rpcMessage(t *testing.T, err error) string {
	t.Helper()

	rpcErr, ok := err.(interface{ RPCMessage() string })
	require.True(t, ok)

	return rpcErr.RPCMessage()
}

func TestExpected(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "app error", err: domain.ErrInvalidTransition, want: true},
		{name: "invalid request", err: ErrInvalidRequest, want: true},
		{name: "json syntax", err: &json.SyntaxError{}, want: true},
		{name: "json type", err: &json.UnmarshalTypeError{}, want: true},
		{name: "eof", err: io.EOF, want: true},
		{name: "unexpected eof", err: io.ErrUnexpectedEOF, want: true},
		{name: "nil", err: nil, want: false},
		{name: "unknown", err: errBoom, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, Expected(tt.err))
		})
	}
}
