package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/yovily/citi/internal/auth/service"
)

// Mock Auth Service
type mockAuthService struct {
	shouldAuthSucceed bool
	token            string
	shouldTokenErr   bool
}

func (m *mockAuthService) Authenticate(username, password, domain string) (*service.AuthResult, error) {
	if m.shouldAuthSucceed {
		return &service.AuthResult{
			Success: true,
			UserID:  username,
			Role:    "user",
		}, nil
	}
	return nil, fmt.Errorf("authentication failed")
}

func (m *mockAuthService) GenerateToken(userID string) (string, error) {
	if m.shouldTokenErr {
		return "", fmt.Errorf("token generation failed")
	}
	return m.token, nil
}

// Mock Logger
type mockLogger struct {
	errorMsgs []string
}

func (m *mockLogger) Error(msg string, args ...interface{}) {
	m.errorMsgs = append(m.errorMsgs, fmt.Sprintf(msg, args...))
}

func TestHandleAuthentication(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		request      *AuthRequest
		authSuccess  bool
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
			authSuccess:  true,
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
			name:   "authentication failure",
			method: http.MethodPost,
			request: &AuthRequest{
				UserID:   "testuser",
				Password: "wrongpass",
				Domain:   "example.com",
				Role:     "user",
			},
			authSuccess: false,
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
			authSuccess:  true,
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
			authService := &mockAuthService{
				shouldAuthSucceed: tt.authSuccess,
				token:            tt.mockToken,
				shouldTokenErr:   !tt.tokenSuccess,
			}
			logger := &mockLogger{}

			// Create handler
			handler := NewAuthHandler(authService, logger)

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

			// Check error logging for failure cases
			if !tt.authSuccess && len(logger.errorMsgs) == 0 {
				t.Error("Expected error to be logged for authentication failure")
			}
		})
	}
}
