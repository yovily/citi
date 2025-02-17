// internal/auth/handler/auth.go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/yovily/citi/internal/auth/service"
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

type AuthService interface {
	Authenticate(username, password, domain string) (*service.AuthResult, error)
	GenerateToken(userID string) (string, error)
}

type AuthHandler struct {
	authService AuthService
	logger      Logger
}

func NewAuthHandler(authService AuthService, logger Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) HandleAuthentication(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight
	if r.Method == http.MethodOptions {
		h.handleCORS(w)
		return
	}

	// Check method
	if r.Method != http.MethodPost {
		h.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var request AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if request.UserID == "" || request.Password == "" || request.Domain == "" {
		h.respondError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	// Authenticate user
	authResult, err := h.authService.Authenticate(request.UserID, request.Password, request.Domain)
	if err != nil {
		h.logger.Error("Authentication failed", "error", err)
		h.respondError(w, http.StatusUnauthorized, "authentication failed")
		return
	}

	// Use role from auth result if available
	if request.Role == "" {
		request.Role = authResult.Role
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(request.UserID)
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

func (h *AuthHandler) handleCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
		// If we failed to encode the response, try to send a basic error
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to encode response"})
	}
}

func (h *AuthHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, ErrorResponse{
		Error: message,
	})
}
