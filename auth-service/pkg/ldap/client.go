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

type Client struct {
	config Config
	logger Logger
}

func NewClient(config Config, logger Logger) *Client {
	return &Client{
		config: config,
		logger: logger,
	}
}

func (c *Client) Authenticate(username, password string) (*AuthResult, error) {
	// Get LDAP server
	host, err := c.config.LookupSvc.LookupServer(c.config.Domain)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup LDAP server: %w", err)
	}

	// Connect to LDAP
	ldapURL := fmt.Sprintf("ldaps://%s:%s", host, c.config.Port)
	conn, err := ldapv3.DialURL(ldapURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to LDAP: %w", err)
	}
	defer conn.Close()

	// Bind with credentials
	err = conn.Bind(username, password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return &AuthResult{
		Username: username,
		Success:  true,
	}, nil
}

type Logger interface {
	Error(msg string, args ...interface{})
}

type AuthResult struct {
	Username string
	Success  bool
}
