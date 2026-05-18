package errmap

import (
	"errors"
	"net/http"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Mapping struct {
	HTTPStatus int
	GRPCCode   codes.Code
	Message    string
}

func Lookup(err error) Mapping {
	switch {
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return Mapping{HTTPStatus: http.StatusConflict, GRPCCode: codes.AlreadyExists, Message: "user already exists"}
	case errors.Is(err, domain.ErrInvalidCredentials):
		return Mapping{HTTPStatus: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated, Message: "invalid credentials"}
	case errors.Is(err, domain.ErrUserNotFound):
		return Mapping{HTTPStatus: http.StatusNotFound, GRPCCode: codes.NotFound, Message: "user not found"}
	case errors.Is(err, domain.ErrTaskNotFound):
		return Mapping{HTTPStatus: http.StatusNotFound, GRPCCode: codes.NotFound, Message: "task not found"}
	case errors.Is(err, domain.ErrTaskTitleRequired):
		return Mapping{HTTPStatus: http.StatusBadRequest, GRPCCode: codes.InvalidArgument, Message: "task title is required"}
	case errors.Is(err, domain.ErrTaskCompleted):
		return Mapping{HTTPStatus: http.StatusBadRequest, GRPCCode: codes.InvalidArgument, Message: "completed task cannot be modified"}
	case errors.Is(err, domain.ErrInvalidTransition):
		return Mapping{HTTPStatus: http.StatusBadRequest, GRPCCode: codes.InvalidArgument, Message: "invalid status transition"}
	case errors.Is(err, domain.ErrNotificationNotFound):
		return Mapping{HTTPStatus: http.StatusNotFound, GRPCCode: codes.NotFound, Message: "notification not found"}
	default:
		return Mapping{HTTPStatus: http.StatusInternalServerError, GRPCCode: codes.Internal, Message: "internal server error"}
	}
}

func Known(err error) bool {
	if err == nil {
		return false
	}

	return Lookup(err).GRPCCode != codes.Internal
}

func HTTP(err error) (int, string) {
	mapping := Lookup(err)

	return mapping.HTTPStatus, mapping.Message
}

func GRPC(err error) error {
	mapping := Lookup(err)

	return status.Error(mapping.GRPCCode, mapping.Message)
}
