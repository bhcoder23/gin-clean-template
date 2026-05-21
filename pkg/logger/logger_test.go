//nolint:paralleltest // These tests mutate zerolog's package-level global log level.
package logger

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

var errTest = errors.New("test error")

// newBufferedLogger. Helper to replace underlying zerolog writer with a buffer and capture logs.
func newBufferedLogger(level string) (*Logger, *bytes.Buffer) {
	l := New(level)
	buf := &bytes.Buffer{}
	zl := zerolog.New(buf).With().Timestamp().Logger()
	l.logger = &zl

	return l, buf
}

func TestNewSetsGlobalLevel(t *testing.T) {
	cases := []struct {
		in   string
		want zerolog.Level
	}{
		{"debug", zerolog.DebugLevel},
		{"info", zerolog.InfoLevel},
		{"warn", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"unknown", zerolog.InfoLevel}, // default path
	}

	for _, tc := range cases {
		l := New(tc.in)

		if l == nil || l.logger == nil {
			t.Fatalf("New(%q) returned nil logger", tc.in)
		}

		if got := zerolog.GlobalLevel(); got != tc.want {
			t.Fatalf("New(%q) global level = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestInfoAndWarn_LogMessageWithAndWithoutArgs(t *testing.T) {
	l, buf := newBufferedLogger("info")

	l.Info("hello")
	l.Info("hello %s", "world")
	l.Warn("warn %d", 7)

	out := buf.String()

	// Expect level fields and messages present
	if !strings.Contains(out, "\"level\":\"info\"") || !strings.Contains(out, "\"message\":\"hello\"") {
		t.Fatalf("info without args not found in output: %s", out)
	}

	if !strings.Contains(out, "hello world") {
		t.Fatalf("formatted info not found in output: %s", out)
	}

	if !strings.Contains(out, "\"level\":\"warn\"") || !strings.Contains(out, "warn 7") {
		t.Fatalf("warn log not found in output: %s", out)
	}
}

func TestWarn_ErrorValueWithContext(t *testing.T) {
	l, buf := newBufferedLogger("info")
	l.Warn(errTest, "known request issue")

	out := buf.String()

	if !strings.Contains(out, "\"level\":\"warn\"") {
		t.Fatalf("expected warn level, got: %s", out)
	}

	if !strings.Contains(out, "\"message\":\"known request issue\"") {
		t.Fatalf("expected context message, got: %s", out)
	}

	if !strings.Contains(out, "\"error\":\"test error\"") {
		t.Fatalf("expected structured error field, got: %s", out)
	}
}

func TestDebug_RespectsLevel(t *testing.T) {
	// when level is info, debug should not emit
	l, buf := newBufferedLogger("info")
	l.Debug("dbg %d", 1)

	if got := buf.String(); got != "" {
		// zerolog may still emit entries depending on global level, ensure global level is info
		if zerolog.GlobalLevel() == zerolog.InfoLevel && strings.Contains(got, "\"level\":\"debug\"") {
			t.Fatalf("debug should be suppressed at info level, got: %s", got)
		}
	}

	// when level is debug, debug should emit
	l, buf = newBufferedLogger("debug")
	l.Debug("dbg %d", 2)

	out := buf.String()

	if !strings.Contains(out, "\"level\":\"debug\"") || !strings.Contains(out, "dbg 2") {
		t.Fatalf("debug should appear at debug level, got: %s", out)
	}
}

func TestError_StringMessage(t *testing.T) {
	l, buf := newBufferedLogger("info")
	l.Error("boom")

	out := buf.String()

	if !strings.Contains(out, "\"level\":\"error\"") {
		t.Fatalf("expected error level, got: %s", out)
	}

	if !strings.Contains(out, "\"message\":\"boom\"") {
		t.Fatalf("expected message field, got: %s", out)
	}
}

func TestError_ErrorValueWithoutContext(t *testing.T) {
	l, buf := newBufferedLogger("info")
	l.Error(errTest)

	out := buf.String()

	if !strings.Contains(out, "\"level\":\"error\"") {
		t.Fatalf("expected error level, got: %s", out)
	}

	if !strings.Contains(out, "\"message\":\"test error\"") {
		t.Fatalf("expected error message, got: %s", out)
	}

	if !strings.Contains(out, "\"error\":\"test error\"") {
		t.Fatalf("expected structured error field, got: %s", out)
	}
}

func TestError_ErrorValueWithContext(t *testing.T) {
	l, buf := newBufferedLogger("info")
	l.Error(errTest, "restapi - v1 - login")

	out := buf.String()

	if !strings.Contains(out, "\"message\":\"restapi - v1 - login\"") {
		t.Fatalf("expected context message, got: %s", out)
	}

	if !strings.Contains(out, "\"error\":\"test error\"") {
		t.Fatalf("expected structured error field, got: %s", out)
	}

	if strings.Contains(out, "%!(EXTRA") {
		t.Fatalf("unexpected fmt extra marker in output: %s", out)
	}
}

func TestError_ErrorValueWithFormattedContext(t *testing.T) {
	l, buf := newBufferedLogger("info")
	l.Error(errTest, "task %s failed", "task-id-123")

	out := buf.String()

	if !strings.Contains(out, "\"message\":\"task task-id-123 failed\"") {
		t.Fatalf("expected formatted context message, got: %s", out)
	}

	if !strings.Contains(out, "\"error\":\"test error\"") {
		t.Fatalf("expected structured error field, got: %s", out)
	}
}

func TestMsg_UnknownTypeFallsBack(t *testing.T) {
	l, buf := newBufferedLogger("debug")
	l.msg(zerolog.ErrorLevel, 12345)

	out := buf.String()

	if !strings.Contains(out, "\"level\":\"error\"") {
		t.Fatalf("expected error level, got: %s", out)
	}

	if !strings.Contains(out, "message 12345 has unknown type int") {
		t.Fatalf("expected fallback message, got: %s", out)
	}
}

func TestFatal_ExitsAndLogs(t *testing.T) {
	if os.Getenv("LOGGER_FATAL_SUBPROC") == "1" {
		l := New("debug")
		l.Fatal(errTest, "fatal now")

		return
	}

	cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run", t.Name()) //nolint:gosec // it's ok to exec self in tests

	cmd.Env = append(os.Environ(), "LOGGER_FATAL_SUBPROC=1")

	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-nil error due to os.Exit in Fatal, got nil; output: %s", string(out))
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected exec.ExitError, got %T", err)
	}

	if status := exitErr.ExitCode(); status != 1 {
		t.Fatalf("expected exit code 1, got %d; output: %s", status, string(out))
	}

	output := string(out)
	if !strings.Contains(output, "\"message\":\"fatal now\"") {
		t.Fatalf("expected fatal message in output, got: %s", output)
	}

	if !strings.Contains(output, "\"error\":\"test error\"") {
		t.Fatalf("expected structured error field in output, got: %s", output)
	}
}
