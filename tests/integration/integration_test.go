package module_test

import (
    "context"
    "testing"
    "time"

    "github.com/yourusername/module"
)

func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    client := module.New(
        module.WithTimeout(5 * time.Second),
    )

    ctx := context.Background()
    err := client.DoSomething(ctx)
    if err != nil {
        t.Errorf("integration test failed: %v", err)
    }
}