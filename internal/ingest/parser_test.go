package ingest

import "testing"

func TestParseMessage(t *testing.T) {
	input := []byte(`{
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

	result, err := ParseMessage(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.IncidentID != "inc_1001" {
		t.Fatalf("expected incident_id inc_1001, got %s", result.IncidentID)
	}
}

func TestParseMessageEmpty(t *testing.T) {
	_, err := ParseMessage(nil)
	if err == nil {
		t.Fatal("expected error for empty message")
	}
}
