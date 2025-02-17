package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yovily/citi/internal/auth/handler"
	"github.com/yovily/citi/internal/auth/service"
	"github.com/yovily/citi/internal/platform/nats"
)

type config struct {
	natsURL           string
	authServiceSubject string
	ldapServer        string
	ldapPort          int
	jwtSecret         string
	tokenExpiry       time.Duration
}

type logger struct{}

func (l *logger) Error(msg string, args ...interface{}) {
	log.Printf("ERROR: "+msg, args...)
}

func main() {
	cfg := parseConfig()

	// Initialize NATS client
	natsClient, err := nats.NewClient(nats.Config{
		URL:           cfg.natsURL,
		MaxReconnects: 5,
		ReconnectWait: time.Second * 1,
	})
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsClient.Close()

	// Initialize auth service
	authService := service.NewAuthService(service.Config{
		LDAPServer:   cfg.ldapServer,
		LDAPPort:     cfg.ldapPort,
		JWTSecret:    cfg.jwtSecret,
		TokenExpiry:  cfg.tokenExpiry,
	})

	// Initialize auth handler
	authHandler := handler.NewAuthHandler(authService, &logger{})

	// Subscribe to auth requests
	_, err = natsClient.Subscribe(cfg.authServiceSubject, func(data []byte) error {
		var request handler.AuthRequest
		if err := json.Unmarshal(data, &request); err != nil {
			return err
		}

		// Create a mock HTTP response writer to capture the response
		rw := newMockResponseWriter()
		authHandler.HandleAuthentication(rw, newMockRequest(request))

		return natsClient.Publish(cfg.authServiceSubject+".response", rw.body)
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to auth requests: %v", err)
	}

	log.Printf("Auth service listening on subject: %s", cfg.authServiceSubject)
	select {} // Block forever
}

func parseConfig() config {
	cfg := config{}

	flag.StringVar(&cfg.natsURL, "nats-url", "nats://localhost:4222", "NATS server URL")
	flag.StringVar(&cfg.authServiceSubject, "auth-subject", "auth.service", "Auth service NATS subject")
	flag.StringVar(&cfg.ldapServer, "ldap-server", "localhost", "LDAP server address")
	flag.IntVar(&cfg.ldapPort, "ldap-port", 389, "LDAP server port")
	flag.StringVar(&cfg.jwtSecret, "jwt-secret", "your-secret-key", "JWT signing secret")
	flag.DurationVar(&cfg.tokenExpiry, "token-expiry", time.Hour*24, "JWT token expiry duration")

	flag.Parse()

	// Allow environment variable overrides
	if url := os.Getenv("NATS_URL"); url != "" {
		cfg.natsURL = url
	}
	if subject := os.Getenv("AUTH_SERVICE_SUBJECT"); subject != "" {
		cfg.authServiceSubject = subject
	}
	if server := os.Getenv("LDAP_SERVER"); server != "" {
		cfg.ldapServer = server
	}
	if port := os.Getenv("LDAP_PORT"); port != "" {
		if p, err := fmt.Sscanf(port, "%d", &cfg.ldapPort); err != nil || p != 1 {
			log.Fatalf("Invalid LDAP_PORT: %s", port)
		}
	}
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.jwtSecret = secret
	}
	if expiry := os.Getenv("TOKEN_EXPIRY"); expiry != "" {
		var err error
		cfg.tokenExpiry, err = time.ParseDuration(expiry)
		if err != nil {
			log.Fatalf("Invalid TOKEN_EXPIRY: %s", expiry)
		}
	}

	return cfg
}

// Mock HTTP types for handling auth requests
type mockResponseWriter struct {
	headers http.Header
	body    []byte
	status  int
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		headers: make(http.Header),
	}
}

func (w *mockResponseWriter) Header() http.Header {
	return w.headers
}

func (w *mockResponseWriter) Write(b []byte) (int, error) {
	w.body = b
	return len(b), nil
}

func (w *mockResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

func newMockRequest(req handler.AuthRequest) *http.Request {
	body, _ := json.Marshal(req)
	r, _ := http.NewRequest(http.MethodPost, "/auth", bytes.NewReader(body))
	return r
}
