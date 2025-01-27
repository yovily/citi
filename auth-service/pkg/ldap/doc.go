// pkg/ldap/doc.go

// Package ldap provides LDAP (Lightweight Directory Access Protocol) client functionality
// for authentication and user directory services.
//
// The package provides:
//   - LDAP server connection management
//   - User authentication
//   - Secure TLS connections
//   - Platform-independent server resolution
//
// Basic usage:
//
//	config := ldap.Config{
//	    Port:   "3269",
//	    Domain: "example.com",
//	}
//	
//	client := ldap.NewClient(config, logger)
//	
//	result, err := client.Authenticate("username", "password")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Security Considerations:
//   - All connections use LDAPS (LDAP over TLS)
//   - Credentials are never logged
//   - Connection timeouts are enforced
package ldap