package service

import (
	"fmt"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/golang-jwt/jwt/v5"
)

type AuthResult struct {
	Success bool
	UserID  string
	Role    string
}

type Config struct {
	LDAPServer   string
	LDAPPort     int
	JWTSecret    string
	TokenExpiry  time.Duration
}

type AuthService struct {
	config Config
}

func NewAuthService(config Config) *AuthService {
	return &AuthService{
		config: config,
	}
}

func (s *AuthService) Authenticate(username, password, domain string) (*AuthResult, error) {
	// Set default timeout for LDAP operations
	ldap.DefaultTimeout = time.Second * 10

	// Connect to LDAP server with timeout
	l, err := ldap.DialURL(fmt.Sprintf("ldap://%s:%d", s.config.LDAPServer, s.config.LDAPPort))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to LDAP server: %w", err)
	}
	defer l.Close()

	// Set read timeout for operations
	l.SetTimeout(time.Second * 5)

	// Bind with user credentials
	bindDN := fmt.Sprintf("%s@%s", username, domain)
	err = l.Bind(bindDN, password)
	if err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok {
			switch ldapErr.ResultCode {
			case ldap.LDAPResultInvalidCredentials:
				return nil, fmt.Errorf("invalid credentials")
			case ldap.LDAPResultTimeLimitExceeded:
				return nil, fmt.Errorf("operation timed out")
			default:
				return nil, fmt.Errorf("LDAP authentication failed: %w", err)
			}
		}
		return nil, fmt.Errorf("LDAP authentication failed: %w", err)
	}

	// Search for user attributes
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("dc=%s", domain), // Base DN
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 30, false,
		fmt.Sprintf("(&(objectClass=user)(userPrincipalName=%s))", bindDN),
		[]string{"memberOf"}, // Attributes to retrieve
		nil,
	)

	result, err := l.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user attributes: %w", err)
	}

	role := "user"
	if len(result.Entries) > 0 {
		// Check group membership for role determination
		for _, group := range result.Entries[0].GetAttributeValues("memberOf") {
			if group == "cn=admins" {
				role = "admin"
				break
			}
		}
	}

	return &AuthResult{
		Success: true,
		UserID:  username,
		Role:    role,
	}, nil
}

func (s *AuthService) GenerateToken(userID string) (string, error) {
	now := time.Now()
	exp := now.Add(s.config.TokenExpiry)

	claims := jwt.MapClaims{
		"sub": userID,                    // Subject (user ID)
		"exp": exp.Unix(),               // Expiration time
		"iat": now.Unix(),               // Issued at
		"nbf": now.Unix(),               // Not before
		"iss": "citi-auth-service",      // Issuer
		"aud": []string{"citi-services"}, // Audience
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = "1" // Key ID for future key rotation support

	signedToken, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

type AuthService interface {
	Authenticate(username, password, domain string) (*AuthResult, error)
	GenerateToken(userID string) (string, error)
}
