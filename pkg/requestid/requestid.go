// Package requestid carries request/correlation identifiers across transports.
package requestid

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

const (
	// Header is the canonical HTTP header for request correlation.
	Header = "X-Request-ID"
	// MetadataKey is the canonical lowercase metadata key for RPC transports.
	MetadataKey = "x-request-id"
)

type contextKey struct{}

// New creates a new request identifier.
func New() string {
	return uuid.NewString()
}

// Normalize returns id trimmed, or a generated id when empty.
func Normalize(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return New()
	}

	return id
}

// WithContext stores id in ctx.
func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKey{}, Normalize(id))
}

// FromContext returns the id stored in ctx.
func FromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(contextKey{}).(string)

	return id, ok && id != ""
}

// Ensure returns a context with a request id and the id.
func Ensure(ctx context.Context) (next context.Context, id string) {
	if id, ok := FromContext(ctx); ok {
		return ctx, id
	}

	id = New()

	return WithContext(ctx, id), id
}
