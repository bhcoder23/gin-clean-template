package domain_test

import (
	"testing"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTask_Transition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		from      domain.TaskStatus
		to        domain.TaskStatus
		wantErr   bool
		wantState domain.TaskStatus
	}{
		{"todo to in_progress", domain.TaskStatusTodo, domain.TaskStatusInProgress, false, domain.TaskStatusInProgress},
		{"in_progress to done", domain.TaskStatusInProgress, domain.TaskStatusDone, false, domain.TaskStatusDone},
		{"in_progress to todo", domain.TaskStatusInProgress, domain.TaskStatusTodo, false, domain.TaskStatusTodo},
		{"todo to done (invalid)", domain.TaskStatusTodo, domain.TaskStatusDone, true, domain.TaskStatusTodo},
		{"done to todo (invalid)", domain.TaskStatusDone, domain.TaskStatusTodo, true, domain.TaskStatusDone},
		{"done to in_progress (invalid)", domain.TaskStatusDone, domain.TaskStatusInProgress, true, domain.TaskStatusDone},
		{"unknown status (invalid)", domain.TaskStatus("unknown"), domain.TaskStatusTodo, true, domain.TaskStatus("unknown")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			task := domain.Task{Status: tt.from}
			err := task.Transition(tt.to)

			if tt.wantErr {
				require.ErrorIs(t, err, domain.ErrInvalidTransition)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantState, task.Status)
		})
	}
}
