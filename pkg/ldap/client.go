// pkg/ldap/client.go
package ldap

import (
	"fmt"

	ldapv3 "github.com/go-ldap/ldap/v3"
)

type LookupService interface {
	LookupServer(domain string) (string, error)
}

type Config struct {
	Port      string
	Domain    string
	LookupSvc LookupService
}

// Add LDAP interface for mocking
type ldapConnection interface {
	Bind(username, password string) error
	Close() error
}

// Add factory function for LDAP connections
type ldapDialer func(addr string) (ldapConnection, error)

type Client struct {
	config   Config
	logger   Logger
	dialLDAP ldapDialer // Add dialer function
}

func NewClient(config Config, logger Logger) *Client {
	if config.Port == "" || config.Domain == "" || config.LookupSvc == nil {
		return nil
	}

	return &Client{
		config: config,
		logger: logger,
		dialLDAP: func(addr string) (ldapConnection, error) { // Default implementation
			return ldapv3.DialURL(addr)
		},
	}
}

func (c *Client) Authenticate(username, password string) (*AuthResult, error) {
	if username == "" || password == "" {
		c.logger.Error("Empty credentials provided")
		return &AuthResult{Success: false}, fmt.Errorf("empty credentials")
	}

	// Get LDAP server
	host, err := c.config.LookupSvc.LookupServer(c.config.Domain)
	if err != nil {
		c.logger.Error("LDAP lookup failed", "error", err)
		return &AuthResult{Success: false}, fmt.Errorf("failed to lookup LDAP server: %w", err)
	}

	// Connect to LDAP
	ldapURL := fmt.Sprintf("ldaps://%s:%s", host, c.config.Port)
	conn, err := c.dialLDAP(ldapURL)
	if err != nil {
		c.logger.Error("Failed to connect to LDAP", "error", err)
		return &AuthResult{Success: false}, fmt.Errorf("failed to connect to LDAP: %w", err)
	}
	defer conn.Close()

	// Bind with credentials
	err = conn.Bind(username, password)
	if err != nil {
		c.logger.Error("Authentication failed", "error", err)
		return &AuthResult{Success: false}, fmt.Errorf("authentication failed: %w", err)
	}

	c.logger.Info("Authentication successful", "username", username)
	return &AuthResult{
		Username: username,
		Success:  true,
	}, nil
}

type Logger interface {
	Error(msg string, args ...interface{})
	Info(msg string, args ...interface{})
}

type AuthResult struct {
	Username string
	Success  bool
}
