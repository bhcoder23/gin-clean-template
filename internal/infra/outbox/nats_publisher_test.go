package outbox

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNATSPublisherSubjectUsesPrefixAndEventType(t *testing.T) {
	t.Parallel()

	publisher := &NATSPublisher{subjectPrefix: "events"}

	require.Equal(t, "events.task.created", publisher.subject(&Event{EventType: "task.created"}))
}
