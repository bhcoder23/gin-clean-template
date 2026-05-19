package outbox

import "time"

const (
	StatusPending    = "pending"
	StatusPublishing = "publishing"
	StatusPublished  = "published"
	StatusFailed     = "failed"
)

// Event is a durable integration event stored in Postgres before publishing.
type Event struct {
	ID            string
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       []byte
	Headers       map[string]string
	Status        string
	Attempts      int
	AvailableAt   time.Time
	CreatedAt     time.Time
	LockedAt      *time.Time
	PublishedAt   *time.Time
	LastError     *string
}
