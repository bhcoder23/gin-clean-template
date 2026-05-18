package persistent

import (
	"errors"
	"testing"
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

var errRows = errors.New("rows error")

type fakeRows struct {
	index int
	err   error
	tasks []entity.Task
	trs   []entity.Translation
}

func (r *fakeRows) Close() {}

func (r *fakeRows) Err() error {
	return r.err
}

func (r *fakeRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (r *fakeRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r *fakeRows) Next() bool {
	if len(r.tasks) > 0 {
		if r.index < len(r.tasks) {
			return true
		}

		return false
	}

	if r.index < len(r.trs) {
		return true
	}

	return false
}

func (r *fakeRows) Scan(dest ...any) error {
	if len(r.tasks) > 0 {
		task := r.tasks[r.index]
		r.index++

		*(dest[0].(*string)) = task.ID
		*(dest[1].(*string)) = task.UserID
		*(dest[2].(*string)) = task.Title
		*(dest[3].(*string)) = task.Description
		*(dest[4].(*entity.TaskStatus)) = task.Status
		*(dest[5].(*time.Time)) = task.CreatedAt
		*(dest[6].(*time.Time)) = task.UpdatedAt

		return nil
	}

	translation := r.trs[r.index]
	r.index++

	*(dest[0].(*string)) = translation.Source
	*(dest[1].(*string)) = translation.Destination
	*(dest[2].(*string)) = translation.Original
	*(dest[3].(*string)) = translation.Translation

	return nil
}

func (r *fakeRows) Values() ([]any, error) {
	return nil, nil
}

func (r *fakeRows) RawValues() [][]byte {
	return nil
}

func (r *fakeRows) Conn() *pgx.Conn {
	return nil
}

var _ pgx.Rows = (*fakeRows)(nil)

func TestCollectTaskRowsChecksIteratorError(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	rows := &fakeRows{
		err: errRows,
		tasks: []entity.Task{
			{
				ID:        "task-1",
				UserID:    "user-1",
				Title:     "Task",
				Status:    entity.TaskStatusTodo,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	tasks, err := collectTaskRows(rows, 10)

	require.Nil(t, tasks)
	require.ErrorIs(t, err, errRows)
}

func TestCollectTranslationRowsChecksIteratorError(t *testing.T) {
	t.Parallel()

	rows := &fakeRows{
		err: errRows,
		trs: []entity.Translation{
			{
				Source:      "auto",
				Destination: "en",
				Original:    "你好",
				Translation: "hello",
			},
		},
	}

	translations, err := collectTranslationRows(rows)

	require.Nil(t, translations)
	require.ErrorIs(t, err, errRows)
}
