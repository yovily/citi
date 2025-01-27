// Package auth provides LDAP authentication and JWT token generation functionality
// for enterprise applications. It supports both Windows and Linux environments and
// implements secure token generation with configurable expiration times.
//
// Basic usage:
//
//	config := auth.Config{
//		JWTSecret:     []byte("your-secret"),
//		TokenDuration: 24 * time.Hour,
//		LDAPPort:      "3269",
//	}
//	
//	client := auth.NewClient(config)
//	
//	// Generate a token
//	token, err := client.GenerateToken("user123")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// The package provides:
//   - JWT token generation with configurable expiration
//   - LDAP authentication support
//   - Platform-independent LDAP server resolution
//   - Secure default configurations
//
// Security Considerations:
//   - JWTSecret should be at least 32 bytes long
//   - TokenDuration should be set according to your security requirements
//   - LDAP connections are made over TLS by default
//
// Example usage:
//
//	client := auth.NewClient(auth.Config{
//		JWTSecret:      []byte("secret"),
//		TokenDuration:  24 * time.Hour,
//		SessionManager: mySessionManager,
//		Logger:        myLogger,
//	})

//	client.Logout(w, r, sessionManager, logger)
package auth

// Version information
const (
    Version = "1.0.0"
)


