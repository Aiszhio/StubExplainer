package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Aiszhio/StubExplainer/internal/metrics"
	"github.com/Aiszhio/StubExplainer/internal/service"
	"github.com/Aiszhio/StubExplainer/internal/storage"
)

func TestHTTPGatewayGetExplanation(t *testing.T) {
	metricsRecorder := metrics.New()
	memoryStorage := storage.NewMemoryStorage()
	explainerService := service.NewExplainerService(memoryStorage, metricsRecorder)

	message := []byte(`{
		"incident_id": "inc_1001",
		"user_id": "user_501",
		"decision": "blocked",
		"score": 0.97,
		"features": [
			{"name": "request_frequency", "value": "high", "weight": 0.42}
		]
	}`)

	if err := explainerService.ProcessBatch([][]byte{message}); err != nil {
		t.Fatalf("failed to process batch: %v", err)
	}

	grpcServer := NewGRPCServer(explainerService)
	grpcClient := NewLocalClient(grpcServer)
	gateway := NewHTTPGateway(grpcClient, metricsRecorder)

	mux := http.NewServeMux()
	gateway.Register(mux)

	request := httptest.NewRequest(http.MethodGet, "/v1/explanations/inc_1001", nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	explanation, ok := body["explanation"].(map[string]any)
	if !ok {
		t.Fatalf("expected explanation object, got %#v", body)
	}

	if explanation["incident_id"] != "inc_1001" {
		t.Fatalf("expected incident_id inc_1001, got %#v", explanation["incident_id"])
	}
}

func TestHTTPGatewayGetExplanationNotFound(t *testing.T) {
	metricsRecorder := metrics.New()
	memoryStorage := storage.NewMemoryStorage()
	explainerService := service.NewExplainerService(memoryStorage, metricsRecorder)
	grpcServer := NewGRPCServer(explainerService)
	grpcClient := NewLocalClient(grpcServer)
	gateway := NewHTTPGateway(grpcClient, metricsRecorder)

	mux := http.NewServeMux()
	gateway.Register(mux)

	request := httptest.NewRequest(http.MethodGet, "/v1/explanations/unknown", nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
}

func TestMetricsEndpointContainsHTTPMetric(t *testing.T) {
	metricsRecorder := metrics.New()
	memoryStorage := storage.NewMemoryStorage()
	explainerService := service.NewExplainerService(memoryStorage, metricsRecorder)
	grpcServer := NewGRPCServer(explainerService)
	grpcClient := NewLocalClient(grpcServer)
	gateway := NewHTTPGateway(grpcClient, metricsRecorder)

	mux := http.NewServeMux()
	gateway.Register(mux)
	mux.Handle("/metrics", metricsRecorder.Handler())

	request := httptest.NewRequest(http.MethodGet, "/v1/explanations/unknown", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	metricsRequest := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metricsResponse := httptest.NewRecorder()
	mux.ServeHTTP(metricsResponse, metricsRequest)

	if metricsResponse.Code != http.StatusOK {
		t.Fatalf("expected metrics status 200, got %d", metricsResponse.Code)
	}

	if !strings.Contains(metricsResponse.Body.String(), "explainer_http_requests_total") {
		t.Fatal("expected http requests metric in metrics response")
	}
}
