package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Aiszhio/StubExplainer/internal/metrics"
	explainerv1 "github.com/Aiszhio/StubExplainer/pkg/gen/explainer/v1"
)

// HTTPGateway exposes HTTP handlers and delegates business calls to the generated-like gRPC client.
type HTTPGateway struct {
	client  explainerv1.ExplainerServiceClient
	metrics metrics.Recorder
}

func NewHTTPGateway(client explainerv1.ExplainerServiceClient, metrics metrics.Recorder) *HTTPGateway {
	return &HTTPGateway{client: client, metrics: metrics}
}

func (g *HTTPGateway) Register(mux *http.ServeMux) {
	mux.HandleFunc("/v1/explanations/", g.handleGetExplanation)
}

func (g *HTTPGateway) handleGetExplanation(w http.ResponseWriter, r *http.Request) {
	startedAt := time.Now()
	statusCode := http.StatusOK
	defer func() {
		g.metrics.ObserveHTTPRequest(r.Method, "/v1/explanations/{incident_id}", strconv.Itoa(statusCode), time.Since(startedAt))
	}()

	if r.Method != http.MethodGet {
		statusCode = http.StatusMethodNotAllowed
		writeError(w, statusCode, "method not allowed")
		return
	}

	incidentID := strings.TrimPrefix(r.URL.Path, "/v1/explanations/")
	if incidentID == "" || incidentID == r.URL.Path {
		statusCode = http.StatusBadRequest
		writeError(w, statusCode, "incident_id is required")
		return
	}

	response, err := g.client.GetExplanation(r.Context(), &explainerv1.GetExplanationRequest{IncidentId: incidentID})
	if err != nil {
		statusCode = http.StatusNotFound
		writeError(w, statusCode, err.Error())
		return
	}

	writeJSON(w, statusCode, response)
}

func writeJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{"error": message})
}
