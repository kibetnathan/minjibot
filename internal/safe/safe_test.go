package safe

import (
	"bytes"
	"log/slog"
	"testing"
	"time"
)

// TestRecover_SwallowsPanic verifies a panic inside a function guarded by
// Recover does not propagate to the caller.
func TestRecover_SwallowsPanic(t *testing.T) {
	// If Recover fails to catch, this deferred check re-panics and fails.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic escaped Recover: %v", r)
		}
	}()

	func() {
		defer Recover(nil, "test")
		panic("boom")
	}()
	// Reaching this line means the panic was recovered.
}

// TestRecover_LogsWithContext verifies the recovery is logged with the "where"
// tag and the panic value.
func TestRecover_LogsWithContext(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	func() {
		defer Recover(logger, "myHandler")
		panic("kaboom")
	}()

	out := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("myHandler")) {
		t.Errorf("log missing where tag: %s", out)
	}
	if !bytes.Contains(buf.Bytes(), []byte("kaboom")) {
		t.Errorf("log missing panic value: %s", out)
	}
}

// TestRecover_NoPanicIsNoop verifies Recover does nothing when there is no panic.
func TestRecover_NoPanicIsNoop(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	func() {
		defer Recover(logger, "clean")
		// no panic
	}()

	if buf.Len() != 0 {
		t.Errorf("expected no log output when nothing panics, got: %s", buf.String())
	}
}

// TestGo_RecoversInGoroutine verifies a panic in a goroutine started via Go is
// recovered (the test process would crash otherwise).
func TestGo_RecoversInGoroutine(t *testing.T) {
	done := make(chan struct{})
	Go(nil, "worker", func() {
		defer close(done)
		panic("goroutine boom")
	})

	select {
	case <-done:
		// The goroutine ran and its panic was recovered by Go's deferred Recover.
	case <-time.After(2 * time.Second):
		t.Fatal("goroutine did not complete")
	}
}
