// internal/platform/lookup_test.go

package platform

import (
    "strings"
    "testing"
)

func TestLookupService(t *testing.T) {
    service := NewLookupService()
    if service == nil {
        t.Fatal("Expected non-nil service")
    }
}

func TestLookupServer(t *testing.T) {
    tests := []struct {
        name      string
        domain    string
        wantError bool
    }{
        {
            name:      "valid domain",
            domain:    "example.com",
            wantError: false,
        },
        {
            name:      "empty domain",
            domain:    "",
            wantError: true,
        },
    }

    service := NewLookupService()

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            host, err := service.LookupServer(tt.domain)
            if (err != nil) != tt.wantError {
                t.Errorf("LookupServer() error = %v, wantError %v", err, tt.wantError)
                return
            }

            if !tt.wantError && !strings.Contains(host, tt.domain) {
                t.Errorf("LookupServer() = %v, want to contain domain %v", host, tt.domain)
            }
        })
    }
}