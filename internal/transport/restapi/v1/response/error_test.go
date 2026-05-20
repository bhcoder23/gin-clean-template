package response

import (
	"encoding/json"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError("TASK_NOT_FOUND", "task not found", "req-123")

	if err.Error.Code != "TASK_NOT_FOUND" {
		t.Fatalf("unexpected code: %s", err.Error.Code)
	}
	if err.Error.Message != "task not found" {
		t.Fatalf("unexpected message: %s", err.Error.Message)
	}
	if err.Error.RequestID != "req-123" {
		t.Fatalf("unexpected request id: %s", err.Error.RequestID)
	}
}

func TestNewErrorOmitsEmptyRequestID(t *testing.T) {
	err := NewError("INTERNAL_SERVER_ERROR", "internal server error", "")

	payload, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("marshal error: %v", marshalErr)
	}

	if string(payload) != `{"error":{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}}` {
		t.Fatalf("unexpected payload: %s", payload)
	}
}
