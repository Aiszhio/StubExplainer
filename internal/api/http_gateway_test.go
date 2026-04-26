package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aiszhio/StubExplainer/internal/service"
	"github.com/Aiszhio/StubExplainer/internal/storage"
)

func TestHTTPGatewayGetExplanation(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	explainerService := service.NewExplainerService(memoryStorage)

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
	gateway := NewHTTPGateway(grpcClient)

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
	memoryStorage := storage.NewMemoryStorage()
	explainerService := service.NewExplainerService(memoryStorage)
	grpcServer := NewGRPCServer(explainerService)
	grpcClient := NewLocalClient(grpcServer)
	gateway := NewHTTPGateway(grpcClient)

	mux := http.NewServeMux()
	gateway.Register(mux)

	request := httptest.NewRequest(http.MethodGet, "/v1/explanations/unknown", nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
}
