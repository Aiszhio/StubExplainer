package api

import (
	"encoding/json"
	"net/http"
	"strings"

	explainerv1 "github.com/Aiszhio/StubExplainer/pkg/gen/explainer/v1"
)

// HTTPGateway exposes HTTP handlers and delegates business calls to the generated-like gRPC client.
type HTTPGateway struct {
	client explainerv1.ExplainerServiceClient
}

func NewHTTPGateway(client explainerv1.ExplainerServiceClient) *HTTPGateway {
	return &HTTPGateway{client: client}
}

func (g *HTTPGateway) Register(mux *http.ServeMux) {
	mux.HandleFunc("/v1/explanations/", g.handleGetExplanation)
}

func (g *HTTPGateway) handleGetExplanation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	incidentID := strings.TrimPrefix(r.URL.Path, "/v1/explanations/")
	if incidentID == "" || incidentID == r.URL.Path {
		writeError(w, http.StatusBadRequest, "incident_id is required")
		return
	}

	response, err := g.client.GetExplanation(r.Context(), &explainerv1.GetExplanationRequest{IncidentId: incidentID})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{"error": message})
}
