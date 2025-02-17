package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/yovily/citi/internal/platform/nats"
)

type Config struct {
	AuthServiceSubject string
	RequestTimeout    time.Duration
}

type GatewayHandler struct {
	natsClient *nats.Client
	config     Config
	apiDir     string
}

func NewGatewayHandler(natsClient *nats.Client, config Config, apiDir string) *GatewayHandler {
	return &GatewayHandler{
		natsClient: natsClient,
		config:     config,
		apiDir:     apiDir,
	}
}

// ServeSwaggerUI serves the Swagger UI HTML page
func (h *GatewayHandler) ServeSwaggerUI(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(h.apiDir, "swagger-ui.html"))
}

// ServeSwaggerSpec serves the OpenAPI specification file
func (h *GatewayHandler) ServeSwaggerSpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	http.ServeFile(w, r, filepath.Join(h.apiDir, "swagger.yaml"))
}

func (h *GatewayHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var authRequest map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Forward request to auth service via NATS
	response, err := h.natsClient.Request(h.config.AuthServiceSubject, authRequest, h.config.RequestTimeout)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("auth service error: %v", err))
		return
	}

	// Parse and forward the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
