package persistence

import (
	"errors"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgtype"
)

var errValueExceedsInt64Max = errors.New("value exceeds int64 max")

func textParam(value string) pgtype.Text {
	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}

func textValue(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}

	return value.String
}

func sqlInt64(value uint64) (int64, error) {
	if value > uint64(math.MaxInt64) {
		return 0, fmt.Errorf("%w: %d", errValueExceedsInt64Max, value)
	}

	return int64(value), nil
}

func paginationCounts(limit, offset uint64) (limitCount, offsetCount int64, err error) {
	limitCount, err = sqlInt64(limit)
	if err != nil {
		return 0, 0, fmt.Errorf("limit: %w", err)
	}

	offsetCount, err = sqlInt64(offset)
	if err != nil {
		return 0, 0, fmt.Errorf("offset: %w", err)
	}

	return limitCount, offsetCount, nil
}
