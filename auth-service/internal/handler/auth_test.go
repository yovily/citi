// internal/handler/auth_test.go
package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/yovily/customers/citi/auth-service/pkg/ldap"
)

// Add after the imports
type AuthResult struct {
	Success bool
}

// Mock LDAP Client
type mockLDAPClient struct {
	shouldSucceed bool
	lastUsername  string
	lastPassword  string
}

func (m *mockLDAPClient) Authenticate(username, password string) (*ldap.AuthResult, error) {
	m.lastUsername = username
	m.lastPassword = password
	if m.shouldSucceed {
		return &ldap.AuthResult{Success: true}, nil
	}
	return &ldap.AuthResult{Success: false}, nil
}

// Mock Auth Client
type mockAuthClient struct {
	token     string
	shouldErr bool
}

func (m *mockAuthClient) GenerateToken(userID string) (string, error) {
	if m.shouldErr {
		return "", fmt.Errorf("token generation failed")
	}
	return m.token, nil
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

func TestHandleAuthentication(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		request      *AuthRequest
		ldapSuccess  bool
		tokenSuccess bool
		mockToken    string
		wantStatus   int
		wantResponse interface{}
	}{
		{
			name:   "successful authentication",
			method: http.MethodPost,
			request: &AuthRequest{
				UserID:   "testuser",
				Password: "testpass",
				Domain:   "example.com",
				Role:     "user",
			},
			ldapSuccess:  true,
			tokenSuccess: true,
			mockToken:    "valid.jwt.token",
			wantStatus:   http.StatusOK,
			wantResponse: AuthResponse{
				UserID:          "testuser",
				IsAuthenticated: true,
				Role:            "user",
				Token:           "valid.jwt.token",
			},
		},
		{
			name:       "invalid method",
			method:     http.MethodGet,
			wantStatus: http.StatusMethodNotAllowed,
			wantResponse: ErrorResponse{
				Error: "invalid request",
			},
		},
		{
			name:   "ldap authentication failure",
			method: http.MethodPost,
			request: &AuthRequest{
				UserID:   "testuser",
				Password: "wrongpass",
				Domain:   "example.com",
				Role:     "user",
			},
			ldapSuccess: false,
			wantStatus:  http.StatusUnauthorized,
			wantResponse: ErrorResponse{
				Error: "authentication failed",
			},
		},
		{
			name:   "token generation failure",
			method: http.MethodPost,
			request: &AuthRequest{
				UserID:   "testuser",
				Password: "testpass",
				Domain:   "example.com",
				Role:     "user",
			},
			ldapSuccess:  true,
			tokenSuccess: false,
			wantStatus:   http.StatusInternalServerError,
			wantResponse: ErrorResponse{
				Error: "token generation failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			ldapClient := &mockLDAPClient{shouldSucceed: tt.ldapSuccess}
			authClient := &mockAuthClient{
				token:     tt.mockToken,
				shouldErr: !tt.tokenSuccess,
			}
			logger := &mockLogger{}

			// Create handler
			handler := NewAuthHandler(ldapClient, authClient, logger)

			// Create request
			var body []byte
			if tt.request != nil {
				body, _ = json.Marshal(tt.request)
			}
			req := httptest.NewRequest(tt.method, "/auth", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Handle request
			handler.HandleAuthentication(rr, req)

			// Check status code
			if rr.Code != tt.wantStatus {
				t.Errorf("HandleAuthentication() status = %v, want %v", rr.Code, tt.wantStatus)
			}

			// Check response body
			var got interface{}
			if tt.wantResponse != nil {
				switch tt.wantResponse.(type) {
				case AuthResponse:
					var resp AuthResponse
					json.NewDecoder(rr.Body).Decode(&resp)
					got = resp
				case ErrorResponse:
					var resp ErrorResponse
					json.NewDecoder(rr.Body).Decode(&resp)
					got = resp
				}

				if !reflect.DeepEqual(got, tt.wantResponse) {
					t.Errorf("HandleAuthentication() response = %v, want %v", got, tt.wantResponse)
				}
			}

			// Additional checks for successful auth
			if tt.ldapSuccess && tt.request != nil {
				expectedUsername := tt.request.UserID + "@" + tt.request.Domain
				if ldapClient.lastUsername != expectedUsername {
					t.Errorf("LDAP username = %v, want %v", ldapClient.lastUsername, expectedUsername)
				}
			}
		})
	}
}
