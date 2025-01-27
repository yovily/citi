// Add these interfaces to a new file: pkg/auth/interfaces.go
package auth

import "context"

type SessionManager interface {
	Put(ctx context.Context, key string, val interface{})
}

type Logger interface {
	Info(msg string, keyvals ...interface{})
}
