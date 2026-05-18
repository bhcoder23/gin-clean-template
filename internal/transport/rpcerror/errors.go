package rpcerror

import (
	"errors"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
)

var (
	ErrInvalidRequest       = errors.New("invalid request body")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserNotFound         = errors.New("user not found")
	ErrTaskNotFound         = errors.New("task not found")
	ErrTaskTitleRequired    = errors.New("task title is required")
	ErrTaskCompleted        = errors.New("completed task cannot be modified")
	ErrInvalidTransition    = errors.New("invalid status transition")
	ErrNotificationNotFound = errors.New("notification not found")
)

func Normalize(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return ErrUserAlreadyExists
	case errors.Is(err, domain.ErrInvalidCredentials):
		return ErrInvalidCredentials
	case errors.Is(err, domain.ErrUserNotFound):
		return ErrUserNotFound
	case errors.Is(err, domain.ErrTaskNotFound):
		return ErrTaskNotFound
	case errors.Is(err, domain.ErrTaskTitleRequired):
		return ErrTaskTitleRequired
	case errors.Is(err, domain.ErrTaskCompleted):
		return ErrTaskCompleted
	case errors.Is(err, domain.ErrInvalidTransition):
		return ErrInvalidTransition
	case errors.Is(err, domain.ErrNotificationNotFound):
		return ErrNotificationNotFound
	default:
		return err
	}
}

func IsKnown(err error) bool {
	switch {
	case errors.Is(err, ErrInvalidRequest):
		return true
	case errors.Is(err, ErrUnauthorized):
		return true
	case errors.Is(err, ErrUserAlreadyExists):
		return true
	case errors.Is(err, ErrInvalidCredentials):
		return true
	case errors.Is(err, ErrUserNotFound):
		return true
	case errors.Is(err, ErrTaskNotFound):
		return true
	case errors.Is(err, ErrTaskTitleRequired):
		return true
	case errors.Is(err, ErrTaskCompleted):
		return true
	case errors.Is(err, ErrInvalidTransition):
		return true
	case errors.Is(err, ErrNotificationNotFound):
		return true
	default:
		return false
	}
}
