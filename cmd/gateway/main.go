package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yovily/citi/internal/gateway/handler"
	"github.com/yovily/citi/internal/platform/nats"
)

type config struct {
	httpPort          int
	natsURL           string
	authServiceSubject string
	requestTimeout    time.Duration
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

	// Initialize gateway handler with API directory path
	gatewayHandler := handler.NewGatewayHandler(natsClient, handler.Config{
		AuthServiceSubject: cfg.authServiceSubject,
		RequestTimeout:    cfg.requestTimeout,
	}, "./api") // absolute path to api directory from project root

	// Set up HTTP routes
	http.HandleFunc("/auth", gatewayHandler.HandleAuth)
	
	// Swagger documentation routes
	http.HandleFunc("/docs", gatewayHandler.ServeSwaggerUI)
	http.HandleFunc("/swagger.yaml", gatewayHandler.ServeSwaggerSpec)

	// Start HTTP server
	serverAddr := fmt.Sprintf(":%d", cfg.httpPort)
	log.Printf("Gateway server starting on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}

func parseConfig() config {
	cfg := config{}

	flag.IntVar(&cfg.httpPort, "port", 8080, "HTTP server port")
	flag.StringVar(&cfg.natsURL, "nats-url", "nats://localhost:4222", "NATS server URL")
	flag.StringVar(&cfg.authServiceSubject, "auth-subject", "auth.service", "Auth service NATS subject")
	flag.DurationVar(&cfg.requestTimeout, "timeout", time.Second*5, "Request timeout duration")

	flag.Parse()

	// Allow environment variable overrides
	if port := os.Getenv("GATEWAY_PORT"); port != "" {
		if p, err := fmt.Sscanf(port, "%d", &cfg.httpPort); err != nil || p != 1 {
			log.Fatalf("Invalid GATEWAY_PORT: %s", port)
		}
	}
	if url := os.Getenv("NATS_URL"); url != "" {
		cfg.natsURL = url
	}
	if subject := os.Getenv("AUTH_SERVICE_SUBJECT"); subject != "" {
		cfg.authServiceSubject = subject
	}
	if timeout := os.Getenv("REQUEST_TIMEOUT"); timeout != "" {
		var err error
		cfg.requestTimeout, err = time.ParseDuration(timeout)
		if err != nil {
			log.Fatalf("Invalid REQUEST_TIMEOUT: %s", timeout)
		}
	}

	return cfg
}
