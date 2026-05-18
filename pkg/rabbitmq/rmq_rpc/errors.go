package rmqrpc

import "errors"

var (
	// ErrTimeout -.
	ErrTimeout = errors.New("timeout")
	// ErrInternalServer -.
	ErrInternalServer = errors.New("internal server error")
	// ErrBadHandler -.
	ErrBadHandler = errors.New("unregistered handler")
	// ErrInvalidRequest -.
	ErrInvalidRequest = errors.New("invalid request body")
	// ErrUnauthorized -.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrUserAlreadyExists -.
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrInvalidCredentials -.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserNotFound -.
	ErrUserNotFound = errors.New("user not found")
	// ErrTaskNotFound -.
	ErrTaskNotFound = errors.New("task not found")
	// ErrTaskTitleRequired -.
	ErrTaskTitleRequired = errors.New("task title is required")
	// ErrTaskCompleted -.
	ErrTaskCompleted = errors.New("completed task cannot be modified")
	// ErrInvalidTransition -.
	ErrInvalidTransition = errors.New("invalid status transition")
	// ErrNotificationNotFound -.
	ErrNotificationNotFound = errors.New("notification not found")
)

// Success -.
const Success = "success"

func ErrorFromStatus(status string) error {
	switch status {
	case Success:
		return nil
	case ErrBadHandler.Error():
		return ErrBadHandler
	case ErrInternalServer.Error():
		return ErrInternalServer
	case ErrInvalidRequest.Error():
		return ErrInvalidRequest
	case ErrUnauthorized.Error():
		return ErrUnauthorized
	case ErrUserAlreadyExists.Error():
		return ErrUserAlreadyExists
	case ErrInvalidCredentials.Error():
		return ErrInvalidCredentials
	case ErrUserNotFound.Error():
		return ErrUserNotFound
	case ErrTaskNotFound.Error():
		return ErrTaskNotFound
	case ErrTaskTitleRequired.Error():
		return ErrTaskTitleRequired
	case ErrTaskCompleted.Error():
		return ErrTaskCompleted
	case ErrInvalidTransition.Error():
		return ErrInvalidTransition
	case ErrNotificationNotFound.Error():
		return ErrNotificationNotFound
	default:
		return nil
	}
}

func IsKnownStatus(status string) bool {
	return ErrorFromStatus(status) != nil
}
