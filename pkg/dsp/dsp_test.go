package dsp

import (
    "context"
    "io"
    "testing"
    "time"
    
    "github.com/citi/dsp/internal/auditlogger"
    "log/slog"
)

func TestNew(t *testing.T) {
    dsp := New()
    if dsp == nil {
        t.Error("Expected non-nil DSP client")
    }
}

func setupTestContext() context.Context {
    // Create a context with the required audit values
    ctx := context.Background()
    logger := slog.New(slog.NewTextHandler(io.Discard, nil))
    values := auditlogger.Values{
        Uuid:   "test-uuid",
        Soeid:  "test-soeid",
        Logger: logger,
    }
    return context.WithValue(ctx, "auditValues", values)
}

func TestLookupDataserverResourceWithTimeout(t *testing.T) {
    // Create a context with timeout to prevent hanging
    ctx, cancel := context.WithTimeout(setupTestContext(), 5*time.Second)
    defer cancel()

    dsp := New()
    _, err := dsp.LookupDataserverResource(ctx, "test-resource")
    if err == nil {
        t.Error("Expected an error when calling external service in test")
    }
}