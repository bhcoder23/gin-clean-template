package outbox

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5"
)

func scanEvent(rows pgx.Rows) (Event, error) {
	var (
		event      Event
		headersRaw []byte
	)

	err := rows.Scan(
		&event.ID,
		&event.AggregateType,
		&event.AggregateID,
		&event.EventType,
		&event.Payload,
		&headersRaw,
		&event.Status,
		&event.Attempts,
		&event.AvailableAt,
		&event.CreatedAt,
		&event.LockedAt,
		&event.PublishedAt,
		&event.LastError,
	)
	if err != nil {
		return Event{}, fmt.Errorf("outbox scanEvent - rows.Scan: %w", err)
	}

	if len(headersRaw) > 0 {
		if err = json.Unmarshal(headersRaw, &event.Headers); err != nil {
			return Event{}, fmt.Errorf("outbox scanEvent - unmarshal headers: %w", err)
		}
	}

	return event, nil
}
