package service

import (
	"testing"

	"github.com/Aiszhio/StubExplainer/internal/metrics"
	"github.com/Aiszhio/StubExplainer/internal/storage"
)

func TestProcessAndGetExplanation(t *testing.T) {
	storage := storage.NewMemoryStorage()
	service := NewExplainerService(storage, metrics.Noop{})

	message := []byte(`{
		"incident_id": "inc_1001",
		"user_id": "user_501",
		"decision": "blocked",
		"score": 0.97,
		"features": [
			{
				"name": "request_frequency",
				"value": "high",
				"weight": 0.42
			}
		]
	}`)

	if err := service.ProcessBatch([][]byte{message}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := service.GetExplanation("inc_1001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Decision != "blocked" {
		t.Fatalf("expected decision blocked, got %s", result.Decision)
	}
}
