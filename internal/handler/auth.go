// internal/handler/auth.go
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yovily/customers/citi/auth-service/pkg/ldap"
)

type Logger interface {
	Error(msg string, args ...interface{})
}

type AuthRequest struct {
	UserID   string
	Password string
	Domain   string
	Role     string
}

type AuthResponse struct {
	UserID          string
	IsAuthenticated bool
	Role            string
	Token           string
}

type ErrorResponse struct {
	Error string
}

type LDAPClient interface {
	Authenticate(username, password string) (*ldap.AuthResult, error)
}

type AuthClient interface {
	GenerateToken(userID string) (string, error)
}

type AuthHandler struct {
	ldapClient LDAPClient
	authClient AuthClient
	logger     Logger
}

func NewAuthHandler(ldapClient LDAPClient, authClient AuthClient, logger Logger) *AuthHandler {
	return &AuthHandler{
		ldapClient: ldapClient,
		authClient: authClient,
		logger:     logger,
	}
}

func (h *AuthHandler) HandleAuthentication(w http.ResponseWriter, r *http.Request) {
	// Check method first
	if r.Method != http.MethodPost {
		h.respondError(w, http.StatusMethodNotAllowed, "invalid request")
		return
	}

	var request AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// Format username for LDAP
	username := fmt.Sprintf("%s@%s", request.UserID, request.Domain)

	// Authenticate with LDAP
	result, err := h.ldapClient.Authenticate(username, request.Password)
	if err != nil || !result.Success {
		h.logger.Error("LDAP authentication failed", "error", err)
		h.respondError(w, http.StatusUnauthorized, "authentication failed")
		return
	}

	// Generate JWT token
	token, err := h.authClient.GenerateToken(request.UserID)
	if err != nil {
		h.logger.Error("Token generation failed", "error", err)
		h.respondError(w, http.StatusInternalServerError, "token generation failed")
		return
	}

	// Create response
	response := AuthResponse{
		UserID:          request.UserID,
		IsAuthenticated: true,
		Role:            request.Role,
		Token:           token,
	}

	h.respondJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
	}
}

func (h *AuthHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, ErrorResponse{
		Error: message,
	})
}
