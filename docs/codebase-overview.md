# Authentication Service Codebase Overview

## Architecture Overview

The codebase implements a robust authentication service with LDAP integration, JWT token management, and service communication through NATS messaging.

## Core Components

### 1. Authentication Handler (`internal/auth/handler/auth.go`)

#### Key Types
```go
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
```

#### Features
- HTTP request processing
- CORS support
- Input validation
- LDAP authentication
- JWT token generation
- Structured error responses

### 2. Authentication Service (`internal/auth/service/auth.go`)

#### Configuration
```go
type Config struct {
    LDAPServer   string
    LDAPPort     int
    JWTSecret    string
    TokenExpiry  time.Duration
}
```

#### Features
- LDAP integration with timeout handling
- Role-based access control
- JWT token generation with claims:
  - Subject (user ID)
  - Expiration time
  - Issued at
  - Not before
  - Issuer
  - Audience
  - Key ID for rotation

### 3. NATS Client (`internal/platform/nats/client.go`)

#### Features
- Message queue integration
- Publish/Subscribe patterns
- Request/Reply functionality
- Automatic reconnection
- Error handling and logging

### 4. Gateway Handler (`internal/gateway/handler/gateway.go`)

#### Features
- HTTP API exposure
- Swagger UI integration
- Request forwarding via NATS
- Error handling
- Response formatting

## Platform Support

### Lookup Service (`internal/platform/lookup.go`)

#### Features
- Platform-specific LDAP server discovery
- Multi-platform support:
  - Linux
  - Windows
  - macOS
- SRV record resolution
- Load balancing through random selection

### Resolver (`internal/platform/resolver/client.go`)

#### Features
- Random host selection
- Thread-safe operation
- Load distribution

## Testing Infrastructure

### Test Suite (`internal/auth/handler/auth_test.go`)

#### Mock Implementations
```go
type mockAuthService struct {
    shouldAuthSucceed bool
    token            string
    shouldTokenErr   bool
}

type mockLogger struct {
    errorMsgs []string
}
```

#### Test Coverage
- Successful authentication
- Invalid method handling
- Authentication failures
- Token generation failures
- Error logging verification

## Configuration Management

### Features
- Environment variable support
- Command-line flags
- Default values
- Runtime configuration
- Flexible overrides

## Security Features

### Implementation
- CORS protection
- Request validation
- Secure token generation
- Timeout handling
- Input sanitization
- Error handling

## Integration Points

### External Systems
- LDAP authentication
- NATS messaging
- HTTP API
- JWT management
- Swagger documentation

## Error Handling

### Comprehensive Coverage
- HTTP status codes
- Detailed error messages
- Error logging
- Client-safe responses
- Timeout management

## Deployment Considerations

### Key Points
- Environment-based configuration
- LDAP server discovery
- Automatic reconnection
- Platform adaptations
- Documentation availability

## Best Practices

### Implementation
1. **Security**
   - Input validation
   - Secure communication
   - Token management
   - Error handling

2. **Scalability**
   - Message queue architecture
   - Load balancing
   - Connection pooling

3. **Maintainability**
   - Clear code structure
   - Comprehensive testing
   - Documentation
   - Error handling

4. **Reliability**
   - Timeout handling
   - Reconnection logic
   - Error recovery
   - Logging

## Code Examples

### Authentication Flow
```go
func (h *AuthHandler) HandleAuthentication(w http.ResponseWriter, r *http.Request) {
    // 1. Validate request
    if r.Method != http.MethodPost {
        h.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }

    // 2. Parse request
    var request AuthRequest
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    // 3. Authenticate
    authResult, err := h.authService.Authenticate(request.UserID, request.Password, request.Domain)
    if err != nil {
        h.logger.Error("Authentication failed", "error", err)
        h.respondError(w, http.StatusUnauthorized, "authentication failed")
        return
    }

    // 4. Generate token
    token, err := h.authService.GenerateToken(request.UserID)
    if err != nil {
        h.logger.Error("Token generation failed", "error", err)
        h.respondError(w, http.StatusInternalServerError, "token generation failed")
        return
    }

    // 5. Send response
    response := AuthResponse{
        UserID:          request.UserID,
        IsAuthenticated: true,
        Role:            authResult.Role,
        Token:           token,
    }
    h.respondJSON(w, http.StatusOK, response)
}
```

This codebase provides a production-ready authentication service suitable for enterprise deployment, with proper attention to security, scalability, and maintainability considerations.
