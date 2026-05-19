package outbox

import (
	"context"
	"fmt"
	"strings"

	"github.com/bhcoder23/gin-clean-template/pkg/requestid"
	"github.com/nats-io/nats.go"
)

// NATSPublisher publishes outbox events to NATS.
type NATSPublisher struct {
	conn          *nats.Conn
	subjectPrefix string
}

// NewNATSPublisher connects a NATS-backed outbox publisher.
func NewNATSPublisher(url, subjectPrefix string) (*NATSPublisher, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("outbox NewNATSPublisher - nats.Connect: %w", err)
	}

	if subjectPrefix == "" {
		subjectPrefix = "events"
	}

	return &NATSPublisher{
		conn:          conn,
		subjectPrefix: strings.Trim(subjectPrefix, "."),
	}, nil
}

// Publish publishes one event to NATS.
func (p *NATSPublisher) Publish(ctx context.Context, event *Event) error {
	msg := nats.NewMsg(p.subject(event))
	msg.Data = event.Payload
	msg.Header.Set("event-id", event.ID)
	msg.Header.Set("event-type", event.EventType)
	msg.Header.Set("aggregate-type", event.AggregateType)
	msg.Header.Set("aggregate-id", event.AggregateID)

	if id, ok := requestid.FromContext(ctx); ok {
		msg.Header.Set(requestid.MetadataKey, id)
	}

	for key, value := range event.Headers {
		msg.Header.Set(key, value)
	}

	if err := p.conn.PublishMsg(msg); err != nil {
		return fmt.Errorf("outbox NATSPublisher - PublishMsg: %w", err)
	}

	if err := p.conn.FlushWithContext(ctx); err != nil {
		return fmt.Errorf("outbox NATSPublisher - FlushWithContext: %w", err)
	}

	return nil
}

// Shutdown closes the publisher.
func (p *NATSPublisher) Shutdown() {
	p.conn.Close()
}

func (p *NATSPublisher) subject(event *Event) string {
	return p.subjectPrefix + "." + event.EventType
}
