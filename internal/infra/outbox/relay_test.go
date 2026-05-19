package outbox

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakePublisher struct{}

func (fakePublisher) Publish(context.Context, *Event) error {
	return nil
}

func TestNewRelayDefaultsPublishTimeout(t *testing.T) {
	t.Parallel()

	relay := NewRelay(&Store{}, fakePublisher{}, nil, RelayConfig{})

	require.Equal(t, defaultPublishTimeout, relay.publishTimeout)
}

func TestNewRelayUsesConfiguredPublishTimeout(t *testing.T) {
	t.Parallel()

	relay := NewRelay(&Store{}, fakePublisher{}, nil, RelayConfig{PublishTimeout: 2 * time.Second})

	require.Equal(t, 2*time.Second, relay.publishTimeout)
}
