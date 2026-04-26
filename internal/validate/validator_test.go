package validate

import (
	"testing"

	"github.com/Aiszhio/StubExplainer/internal/model"
)

func TestValidateAndInit(t *testing.T) {
	input := model.Explanation{
		IncidentID: "inc_1001",
		UserID:     "user_501",
		Decision:   "blocked",
		Score:      0.97,
	}

	result, err := ValidateAndInit(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.CreatedAt.IsZero() {
		t.Fatal("created_at must be initialized")
	}
}

func TestValidateAndInitMissingIncidentID(t *testing.T) {
	input := model.Explanation{
		UserID:   "user_501",
		Decision: "blocked",
	}

	_, err := ValidateAndInit(input)
	if err == nil {
		t.Fatal("expected error for missing incident_id")
	}
}
