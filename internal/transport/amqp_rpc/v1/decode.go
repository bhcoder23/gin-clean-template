package v1

import (
	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/goccy/go-json"
)

func (r *V1) decodeAndValidate(data []byte, target any) error {
	if err := json.Unmarshal(data, target); err != nil {
		return apperror.RPC(apperror.ErrInvalidRequest)
	}

	if err := r.v.Struct(target); err != nil {
		return apperror.RPC(apperror.ErrInvalidRequest)
	}

	return nil
}

func (r *V1) decodeOptionalAndValidate(data []byte, target any) error {
	if len(data) == 0 {
		return nil
	}

	return r.decodeAndValidate(data, target)
}
