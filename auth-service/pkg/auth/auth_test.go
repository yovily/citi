package auth

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Mock types for testing
type mockSessionManager struct {
	values map[string]interface{}
}

func newMockSessionManager() *mockSessionManager {
	return &mockSessionManager{
		values: make(map[string]interface{}),
	}
}

func (m *mockSessionManager) Put(ctx context.Context, key string, val interface{}) {
	m.values[key] = val
}

func (m *mockSessionManager) Get(key string) interface{} {
	return m.values[key]
}

type mockLogger struct {
	logs []string
}

func newMockLogger() *mockLogger {
	return &mockLogger{
		logs: make([]string, 0),
	}
}

func (m *mockLogger) Info(msg string, keyvals ...interface{}) {
	logEntry := msg
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			logEntry += fmt.Sprintf(" %v=%v", keyvals[i], keyvals[i+1])
		}
	}
	m.logs = append(m.logs, logEntry)
}

func TestNewClient(t *testing.T) {
	// Test client creation
	config := Config{
		JWTSecret:     []byte("test-secret"),
		TokenDuration: 24 * time.Hour,
		LDAPPort:      "3269",
	}

	client := NewClient(config)
	if client == nil {
		t.Error("Expected non-nil client")
	}
}

func TestGenerateToken(t *testing.T) {
	// Test cases to run
	tests := []struct {
		name      string
		userID    string
		secret    []byte
		duration  time.Duration
		wantError bool
	}{
		{
			name:      "Valid token generation",
			userID:    "test-user",
			secret:    []byte("test-secret"),
			duration:  24 * time.Hour,
			wantError: false,
		},
		{
			name:      "Empty userID",
			userID:    "",
			secret:    []byte("test-secret"),
			duration:  24 * time.Hour,
			wantError: false,
		},
		{
			name:      "Empty secret",
			userID:    "test-user",
			secret:    []byte(""),
			duration:  24 * time.Hour,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create client with test configuration
			client := NewClient(Config{
				JWTSecret:     tt.secret,
				TokenDuration: tt.duration,
				LDAPPort:      "3269",
			})

			// Generate token
			token, err := client.GenerateToken(tt.userID)

			// Check error expectation
			if (err != nil) != tt.wantError {
				t.Errorf("GenerateToken() error = %v, wantError %v", err, tt.wantError)
				return
			}

			// If we don't expect an error, validate the token
			if !tt.wantError {
				// Parse and validate the token
				parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return tt.secret, nil
				})

				if err != nil {
					t.Errorf("Failed to parse generated token: %v", err)
					return
				}

				// Verify claims
				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				if !ok {
					t.Error("Failed to parse token claims")
					return
				}

				// Verify username claim
				if username, ok := claims["username"].(string); !ok || username != tt.userID {
					t.Errorf("Token username = %v, want %v", username, tt.userID)
				}

				// Verify expiration time
				if exp, ok := claims["exp"].(float64); !ok {
					t.Error("Token expiration time not found or invalid")
				} else {
					// Check if expiration is roughly correct (within 1 second tolerance)
					expectedExp := time.Now().Add(tt.duration).Unix()
					if int64(exp) < expectedExp-1 || int64(exp) > expectedExp+1 {
						t.Errorf("Token expiration time = %v, want close to %v", int64(exp), expectedExp)
					}
				}
			}
		})
	}

	t.Run("Empty_secret", func(t *testing.T) {
		// Create client with empty config to simulate invalid setup
		client := &Client{
			config: Config{},
		}

		_, err := client.GenerateToken("test-user")

		// Assert that we get an error about invalid configuration
		if err == nil {
			t.Error("expected error for empty config, got nil")
		}
	})
}

// TestTokenExpiration verifies that generated tokens actually expire
func TestTokenExpiration(t *testing.T) {
	client := NewClient(Config{
		JWTSecret:     []byte("test-secret"),
		TokenDuration: 1 * time.Second, // Short duration for testing
		LDAPPort:      "3269",
	})

	// Generate token
	token, err := client.GenerateToken("test-user")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait for token to expire
	time.Sleep(2 * time.Second)

	// Verify token is expired
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})

	if err == nil || parsedToken.Valid {
		t.Error("Token should be expired")
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name         string
		setupContext func(context.Context) context.Context
		wantLogEntry string
		wantRedirect string
		wantError    bool
	}{
		{
			name: "successful logout with user context",
			setupContext: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, "userID", "test-user")
			},
			wantLogEntry: "User Logged Out message=User Logged Out userID=test-user",
			wantRedirect: "/login",
			wantError:    false,
		},
		{
			name: "logout without user context",
			setupContext: func(ctx context.Context) context.Context {
				return ctx
			},
			wantLogEntry: "User Logged Out message=User Logged Out userID=<nil>",
			wantRedirect: "/login",
			wantError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionMgr := newMockSessionManager()
			logger := newMockLogger()

			client := NewClient(Config{
				JWTSecret:     []byte("test-secret"),
				TokenDuration: time.Hour,
			})

			r := httptest.NewRequest("GET", "/logout", nil)
			r = r.WithContext(tt.setupContext(r.Context()))
			w := httptest.NewRecorder()

			err := client.Logout(w, r, sessionMgr, logger)

			// Check error
			if (err != nil) != tt.wantError {
				t.Errorf("Logout() error = %v, wantError %v", err, tt.wantError)
			}

			// Check session cleared
			if sessionMgr.Get("userID") != "" {
				t.Error("userID was not cleared")
			}
			if sessionMgr.Get("jwt_token") != "" {
				t.Error("jwt_token was not cleared")
			}
			if sessionMgr.Get("authenticated") != false {
				t.Error("authenticated was not set to false")
			}

			// Check redirect
			if got := w.Header().Get("Location"); got != tt.wantRedirect {
				t.Errorf("Redirect location = %v, want %v", got, tt.wantRedirect)
			}

			// Check logging
			if len(logger.logs) == 0 {
				t.Error("No log entries were created")
			}
			if logger.logs[0] != tt.wantLogEntry {
				t.Errorf("Log entry = %v, want %v", logger.logs[0], tt.wantLogEntry)
			}
		})
	}
}

func ExampleClient_GenerateToken() {
	client := NewClient(Config{
		JWTSecret:     []byte("example-secret"),
		TokenDuration: 24 * time.Hour,
		LDAPPort:      "3269",
	})

	token, err := client.GenerateToken("example-user")
	if err != nil {
		panic(err)
	}

	// Use the token
	_ = token
}
