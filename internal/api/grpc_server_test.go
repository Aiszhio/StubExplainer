package api

import (
	"context"
	"testing"

	"github.com/Aiszhio/StubExplainer/internal/service"
	"github.com/Aiszhio/StubExplainer/internal/storage"
	explainerv1 "github.com/Aiszhio/StubExplainer/pkg/gen/explainer/v1"
)

func TestGRPCServerGetExplanation(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	explainerService := service.NewExplainerService(memoryStorage)

	message := []byte(`{
		"incident_id": "inc_1001",
		"user_id": "user_501",
		"decision": "blocked",
		"score": 0.97
	}`)

	if err := explainerService.ProcessBatch([][]byte{message}); err != nil {
		t.Fatalf("failed to process batch: %v", err)
	}

	server := NewGRPCServer(explainerService)
	response, err := server.GetExplanation(context.Background(), &explainerv1.GetExplanationRequest{IncidentId: "inc_1001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Explanation.Decision != "blocked" {
		t.Fatalf("expected decision blocked, got %s", response.Explanation.Decision)
	}
}
