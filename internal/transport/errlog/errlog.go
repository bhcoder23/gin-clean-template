package errlog

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/bhcoder23/gin-clean-template/internal/transport/errmap"
	"github.com/bhcoder23/gin-clean-template/internal/transport/rpcerror"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
)

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
	case errmap.Known(err):
		return true
	case rpcerror.IsKnown(err):
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
