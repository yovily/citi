// Package ldap provides LDAP (Lightweight Directory Access Protocol) client functionality
// for authentication and user management.
//
// It provides:
//   - LDAP connection management
//   - Authentication methods
//   - Secure TLS connections
//   - Error handling for LDAP operations
//
// Example usage:
//
//	client := ldap.NewClient("3269")
//	conn, err := client.Connect("domain.com")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer conn.Close()
// package ldap