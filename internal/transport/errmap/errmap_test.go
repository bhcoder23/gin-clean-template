package errmap

import (
	"net/http"
	"testing"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLookup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err        error
		httpStatus int
		grpcCode   codes.Code
		message    string
	}{
		{err: domain.ErrUserAlreadyExists, httpStatus: http.StatusConflict, grpcCode: codes.AlreadyExists, message: "user already exists"},
		{err: domain.ErrInvalidCredentials, httpStatus: http.StatusUnauthorized, grpcCode: codes.Unauthenticated, message: "invalid credentials"},
		{err: domain.ErrUserNotFound, httpStatus: http.StatusNotFound, grpcCode: codes.NotFound, message: "user not found"},
		{err: domain.ErrTaskNotFound, httpStatus: http.StatusNotFound, grpcCode: codes.NotFound, message: "task not found"},
		{err: domain.ErrTaskTitleRequired, httpStatus: http.StatusBadRequest, grpcCode: codes.InvalidArgument, message: "task title is required"},
		{err: domain.ErrTaskCompleted, httpStatus: http.StatusBadRequest, grpcCode: codes.InvalidArgument, message: "completed task cannot be modified"},
		{err: domain.ErrInvalidTransition, httpStatus: http.StatusBadRequest, grpcCode: codes.InvalidArgument, message: "invalid status transition"},
		{err: domain.ErrNotificationNotFound, httpStatus: http.StatusNotFound, grpcCode: codes.NotFound, message: "notification not found"},
		{err: assertErr("boom"), httpStatus: http.StatusInternalServerError, grpcCode: codes.Internal, message: "internal server error"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()

			mapping := Lookup(tc.err)
			require.Equal(t, tc.httpStatus, mapping.HTTPStatus)
			require.Equal(t, tc.grpcCode, mapping.GRPCCode)
			require.Equal(t, tc.message, mapping.Message)
		})
	}
}

func TestGRPC(t *testing.T) {
	t.Parallel()

	err := GRPC(domain.ErrTaskNotFound)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, st.Code())
	require.Equal(t, "task not found", st.Message())
}

type assertErr string

func (e assertErr) Error() string {
	return string(e)
}
