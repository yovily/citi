// Package auth provides authentication functionality using LDAP
package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Config holds the configuration for the auth package
type Config struct {
	JWTSecret      []byte
	TokenDuration  time.Duration
	LDAPPort       string
	SessionManager SessionManager
	Logger         Logger
}

// Client handles authentication operations
type Client struct {
	config Config
}

// NewClient creates a new authentication client
func NewClient(config Config) *Client {
	return &Client{
		config: config,
	}
}

// GenerateToken creates a new JWT token for an authenticated user
func (c *Client) GenerateToken(userID string) (string, error) {
	// Validate config
	if len(c.config.JWTSecret) == 0 {
		return "", fmt.Errorf("invalid client configuration: missing JWT secret")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userID
	claims["exp"] = time.Now().Add(c.config.TokenDuration).Unix()

	return token.SignedString(c.config.JWTSecret)
}

// Logout handles user session termination and cleanup
func (c *Client) Logout(w http.ResponseWriter, r *http.Request, sessionManager SessionManager, logger Logger) error {
	// Use request's context instead of separate authCtx and simplified context handling
	ctx := r.Context()

	// Made logging optional and safer
	if logger != nil {
		logger.Info("User Logged Out",
			"message", "User Logged Out",
			"userID", ctx.Value("userID"), // Generic userID instead of soeid
		)
	}

	// Made session management optional
	if sessionManager != nil {
		sessionManager.Put(ctx, "userID", "")
		sessionManager.Put(ctx, "jwt_token", "")
		sessionManager.Put(ctx, "authenticated", false)
	}

	// Same redirect but with error handling potential
	http.Redirect(w, r, "/login", http.StatusFound)
	return nil
}
