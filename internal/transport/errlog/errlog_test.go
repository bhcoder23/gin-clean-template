package errlog

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/transport/rpcerror"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

type validateMe struct {
	Name string `validate:"required"`
}

func TestExpected(t *testing.T) {
	t.Parallel()

	v := validator.New()
	errValidation := v.Struct(validateMe{})

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "validation", err: errValidation, want: true},
		{name: "json syntax", err: &json.SyntaxError{}, want: true},
		{name: "domain", err: domain.ErrTaskNotFound, want: true},
		{name: "rpc", err: rpcerror.ErrUnauthorized, want: true},
		{name: "unknown", err: errors.New("boom"), want: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, Expected(tc.err))
		})
	}
}
