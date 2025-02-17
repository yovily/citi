.PHONY: build run-gateway run-auth test clean

# Build settings
BINARY_DIR = bin
GATEWAY_BINARY = $(BINARY_DIR)/gateway
AUTH_BINARY = $(BINARY_DIR)/auth

# Go settings
GO = go
GOFLAGS = -v

build: clean $(GATEWAY_BINARY) $(AUTH_BINARY)

$(GATEWAY_BINARY):
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -o $(GATEWAY_BINARY) ./cmd/gateway

$(AUTH_BINARY):
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -o $(AUTH_BINARY) ./cmd/auth

run-gateway: $(GATEWAY_BINARY)
	@echo "Starting Gateway Service..."
	@$(GATEWAY_BINARY) \
		--port=8080 \
		--nats-url=nats://localhost:4222 \
		--auth-subject=auth.service \
		--timeout=5s

run-auth: $(AUTH_BINARY)
	@echo "Starting Auth Service..."
	@$(AUTH_BINARY) \
		--nats-url=nats://localhost:4222 \
		--auth-subject=auth.service \
		--ldap-server=localhost \
		--ldap-port=389 \
		--jwt-secret=your-secret-key \
		--token-expiry=24h

test:
	$(GO) test -v ./...

clean:
	@rm -rf $(BINARY_DIR)

# Development helpers
dev-deps:
	@echo "Installing development dependencies..."
	$(GO) install github.com/nats-io/nats-server/v2@latest

run-nats:
	@echo "Starting NATS server..."
	nats-server

.DEFAULT_GOAL := build
