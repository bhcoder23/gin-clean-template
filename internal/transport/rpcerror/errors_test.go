package rpcerror

import (
	"errors"
	"testing"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestNormalize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want error
	}{
		{name: "nil", err: nil, want: nil},
		{name: "user exists", err: errors.Join(domain.ErrUserAlreadyExists), want: ErrUserAlreadyExists},
		{name: "invalid credentials", err: errors.Join(domain.ErrInvalidCredentials), want: ErrInvalidCredentials},
		{name: "user not found", err: errors.Join(domain.ErrUserNotFound), want: ErrUserNotFound},
		{name: "task not found", err: errors.Join(domain.ErrTaskNotFound), want: ErrTaskNotFound},
		{name: "task title required", err: errors.Join(domain.ErrTaskTitleRequired), want: ErrTaskTitleRequired},
		{name: "task completed", err: errors.Join(domain.ErrTaskCompleted), want: ErrTaskCompleted},
		{name: "invalid transition", err: errors.Join(domain.ErrInvalidTransition), want: ErrInvalidTransition},
		{name: "notification not found", err: errors.Join(domain.ErrNotificationNotFound), want: ErrNotificationNotFound},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, Normalize(tc.err))
		})
	}
}
