package apperror

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error struct {
	HTTPStatus int
	GRPCCode   codes.Code
	Code       string
	Message    string
	Expected   bool
}

const (
	CodeInternalServer       = "INTERNAL_SERVER_ERROR"
	CodeInvalidRequest       = "INVALID_REQUEST"
	CodeUnauthorized         = "UNAUTHORIZED"
	CodeUserAlreadyExists    = "USER_ALREADY_EXISTS"
	CodeInvalidCredentials   = "INVALID_CREDENTIALS" // #nosec G101 -- client-facing error code, not a secret.
	CodeUserNotFound         = "USER_NOT_FOUND"
	CodeInvalidUsername      = "INVALID_USERNAME"
	CodeInvalidEmail         = "INVALID_EMAIL"
	CodePasswordTooShort     = "PASSWORD_TOO_SHORT" // #nosec G101 -- client-facing error code, not a secret.
	CodeTaskNotFound         = "TASK_NOT_FOUND"
	CodeTaskTitleRequired    = "TASK_TITLE_REQUIRED"
	CodeTaskCompleted        = "TASK_COMPLETED"
	CodeInvalidTransition    = "INVALID_STATUS_TRANSITION"
	CodeNotificationNotFound = "NOTIFICATION_NOT_FOUND"
)

var (
	ErrInvalidRequest = errors.New("invalid request body")
	ErrUnauthorized   = errors.New("unauthorized")
)

func From(err error) Error {
	for _, known := range knownErrorMappings() {
		if errors.Is(err, known.err) {
			return known.app
		}
	}

	return Error{HTTPStatus: http.StatusInternalServerError, GRPCCode: codes.Internal, Code: CodeInternalServer, Message: "internal server error"}
}

func knownErrorMappings() []struct {
	err error
	app Error
} {
	return []struct {
		err error
		app Error
	}{
		{err: ErrInvalidRequest, app: Error{HTTPStatus: http.StatusBadRequest, GRPCCode: codes.InvalidArgument, Code: CodeInvalidRequest, Message: "invalid request body", Expected: true}},
		{err: ErrUnauthorized, app: Error{HTTPStatus: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated, Code: CodeUnauthorized, Message: "unauthorized", Expected: true}},
		{err: domain.ErrUserAlreadyExists, app: Error{HTTPStatus: http.StatusConflict, GRPCCode: codes.AlreadyExists, Code: CodeUserAlreadyExists, Message: "user already exists", Expected: true}},
		{err: domain.ErrInvalidCredentials, app: Error{HTTPStatus: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated, Code: CodeInvalidCredentials, Message: "invalid credentials", Expected: true}},
		{err: domain.ErrUserNotFound, app: Error{HTTPStatus: http.StatusNotFound, GRPCCode: codes.NotFound, Code: CodeUserNotFound, Message: "user not found", Expected: true}},
		{err: domain.ErrInvalidUsername, app: Error{HTTPStatus: http.StatusBadRequest, GRPCCode: codes.InvalidArgument, Code: CodeInvalidUsername, Message: "invalid username", Expected: true}},
		{err: domain.ErrInvalidEmail, app: Error{HTTPStatus: http.StatusBadRequest, GRPCCode: codes.InvalidArgument, Code: CodeInvalidEmail, Message: "invalid email", Expected: true}},
		{err: domain.ErrPasswordTooShort, app: Error{HTTPStatus: http.StatusBadRequest, GRPCCode: codes.InvalidArgument, Code: CodePasswordTooShort, Message: "password too short", Expected: true}},
		{err: domain.ErrTaskNotFound, app: Error{HTTPStatus: http.StatusNotFound, GRPCCode: codes.NotFound, Code: CodeTaskNotFound, Message: "task not found", Expected: true}},
		{err: domain.ErrTaskTitleRequired, app: Error{HTTPStatus: http.StatusBadRequest, GRPCCode: codes.InvalidArgument, Code: CodeTaskTitleRequired, Message: "task title is required", Expected: true}},
		{err: domain.ErrTaskCompleted, app: Error{HTTPStatus: http.StatusBadRequest, GRPCCode: codes.InvalidArgument, Code: CodeTaskCompleted, Message: "completed task cannot be modified", Expected: true}},
		{err: domain.ErrInvalidTransition, app: Error{HTTPStatus: http.StatusBadRequest, GRPCCode: codes.InvalidArgument, Code: CodeInvalidTransition, Message: "invalid status transition", Expected: true}},
		{err: domain.ErrNotificationNotFound, app: Error{HTTPStatus: http.StatusNotFound, GRPCCode: codes.NotFound, Code: CodeNotificationNotFound, Message: "notification not found", Expected: true}},
	}
}

func Known(err error) bool {
	return err != nil && From(err).Expected
}

func HTTP(err error) (statusCode int, message string) {
	appErr := From(err)

	return appErr.HTTPStatus, appErr.Message
}

func GRPC(err error) error {
	appErr := From(err)

	return status.Error(appErr.GRPCCode, appErr.Message)
}

func RPC(err error) error {
	if err == nil {
		return nil
	}

	appErr := From(err)

	return rpcError{
		code:    appErr.Code,
		message: appErr.Message,
	}
}

func Expected(err error) bool {
	if err == nil {
		return false
	}

	var validationErrs validator.ValidationErrors

	var syntaxErr *json.SyntaxError

	var unmarshalErr *json.UnmarshalTypeError

	switch {
	case errors.As(err, &validationErrs):
		return true
	case errors.As(err, &syntaxErr):
		return true
	case errors.As(err, &unmarshalErr):
		return true
	case errors.Is(err, io.EOF):
		return true
	case errors.Is(err, io.ErrUnexpectedEOF):
		return true
	case Known(err):
		return true
	default:
		return false
	}
}

func Log(l logger.Interface, err error, message string, args ...any) {
	logArgs := append([]any{message}, args...)

	if Expected(err) {
		l.Warn(err, logArgs...)

		return
	}

	l.Error(err, logArgs...)
}

type rpcError struct {
	code    string
	message string
}

func (e rpcError) Error() string {
	return e.message
}

func (e rpcError) RPCCode() string {
	return e.code
}

func (e rpcError) RPCMessage() string {
	return e.message
}
