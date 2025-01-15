// pkg/ldap/client_test.go
package ldap

import (
	"fmt"
	"testing"
)

// Mock LookupService
type mockLookupService struct {
	host string
	err  error
}

func (m *mockLookupService) LookupServer(domain string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.host, nil
}

// Mock Logger
type mockLogger struct {
	infoMsgs  []string
	errorMsgs []string
}

func (m *mockLogger) Info(msg string, keyvals ...interface{}) {
	m.infoMsgs = append(m.infoMsgs, msg)
}

func (m *mockLogger) Error(msg string, keyvals ...interface{}) {
	m.errorMsgs = append(m.errorMsgs, msg)
}

// Add mock LDAP connection
type mockLDAPConn struct {
	shouldError bool
}

func (m *mockLDAPConn) Bind(username, password string) error {
	if m.shouldError {
		return fmt.Errorf("bind error")
	}
	return nil
}

func (m *mockLDAPConn) Close() error {
	return nil
}

func TestNewClient(t *testing.T) {
	mockLookup := &mockLookupService{} // Create mock lookup service

	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: Config{
				Port:      "3269",
				Domain:    "example.com",
				LookupSvc: mockLookup, // Add the mock lookup service
			},
			wantErr: false,
		},
		{
			name: "missing port",
			config: Config{
				Domain:    "example.com",
				LookupSvc: mockLookup,
			},
			wantErr: true,
		},
		{
			name: "missing domain",
			config: Config{
				Port:      "3269",
				LookupSvc: mockLookup,
			},
			wantErr: true,
		},
		{
			name: "missing lookup service",
			config: Config{
				Port:   "3269",
				Domain: "example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mockLogger{}
			client := NewClient(tt.config, logger)

			if tt.wantErr {
				if client != nil {
					t.Error("NewClient() should return nil when config is invalid")
				}
			} else {
				if client == nil {
					t.Error("NewClient() should not return nil when config is valid")
				}
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		password    string
		mockHost    string
		mockErr     error
		bindErr     bool
		wantSuccess bool
		wantErr     bool
		wantInfoLog bool
		wantErrLog  bool
	}{
		{
			name:        "successful authentication",
			username:    "testuser",
			password:    "testpass",
			mockHost:    "ldap.example.com",
			bindErr:     false,
			wantSuccess: true,
			wantErr:     false,
			wantInfoLog: true,
			wantErrLog:  false,
		},
		{
			name:        "lookup service error",
			username:    "testuser",
			password:    "testpass",
			mockErr:     fmt.Errorf("lookup failed"),
			wantSuccess: false,
			wantErr:     true,
		},
		{
			name:        "empty credentials",
			username:    "",
			password:    "",
			mockHost:    "ldap.example.com",
			wantSuccess: false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockConn := &mockLDAPConn{shouldError: tt.bindErr}
			mockDialer := func(addr string) (ldapConnection, error) {
				return mockConn, nil
			}
			logger := &mockLogger{} // Create logger here

			client := NewClient(Config{
				Port:      "3269",
				Domain:    "example.com",
				LookupSvc: &mockLookupService{host: tt.mockHost, err: tt.mockErr},
			}, logger) // Pass logger to NewClient
			client.dialLDAP = mockDialer // Override the dialer with mock

			// Perform authentication
			result, err := client.Authenticate(tt.username, tt.password)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check success
			if !tt.wantErr && result.Success != tt.wantSuccess {
				t.Errorf("Authenticate() success = %v, want %v", result.Success, tt.wantSuccess)
			}

			// Check logging
			if tt.wantErrLog && len(logger.errorMsgs) == 0 {
				t.Error("Expected error to be logged")
			}
			if tt.wantInfoLog && len(logger.infoMsgs) == 0 {
				t.Error("Expected success to be logged")
			}
		})
	}
}

// TestAuthenticateIntegration performs integration tests with actual LDAP server
// This test is skipped unless explicitly enabled
func TestAuthenticateIntegration(t *testing.T) {
	t.Skip("Skipping integration test - requires real LDAP server")

	// Setup real config
	config := Config{
		Port:   "3269",
		Domain: "your.actual.domain",
	}
	logger := &mockLogger{}
	client := NewClient(config, logger)

	result, err := client.Authenticate("testuser@domain.com", "testpass")
	if err != nil {
		t.Errorf("Integration test failed: %v", err)
	}
	if !result.Success {
		t.Error("Integration test authentication failed")
	}
}
