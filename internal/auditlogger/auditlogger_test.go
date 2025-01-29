package auditlogger

import (
    "net/http"
    "testing"
)

func TestInitCtx(t *testing.T) {
    // Create a request with test headers
    req, err := http.NewRequest("GET", "/test", nil)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }

    // Add test headers
    req.Header.Set("UUID", "test-uuid")
    req.Header.Set("SOEID", "test-soeid")

    // Create writer without actual MongoDB connection
    mw := &MongoWriter{DB: nil}

    // Test InitCtx
    ctx := mw.InitCtx(req)

    // Get values from context
    values, ok := ctx.Value("auditValues").(Values)
    if !ok {
        t.Error("Failed to get Values from context")
        return
    }

    // Verify values
    if values.Uuid != "test-uuid" {
        t.Errorf("Expected UUID test-uuid, got %s", values.Uuid)
    }

    if values.Soeid != "test-soeid" {
        t.Errorf("Expected SOEID test-soeid, got %s", values.Soeid)
    }

    if values.Logger == nil {
        t.Error("Logger should not be nil")
    }
}

func TestNewAuditLogger(t *testing.T) {
    logger := New(nil)
    if logger == nil {
        t.Error("Expected non-nil logger")
    }
}

func TestWrite(t *testing.T) {
    // Skip actual MongoDB writes in tests
    t.Skip("Skipping MongoDB write tests")
}