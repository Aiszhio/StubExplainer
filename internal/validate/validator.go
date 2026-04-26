package validate

import (
	"errors"
	"time"

	"github.com/Aiszhio/StubExplainer/internal/model"
)

// ValidateAndInit validates required fields and fills default values.
func ValidateAndInit(explanation model.Explanation) (model.Explanation, error) {
	if explanation.IncidentID == "" {
		return model.Explanation{}, errors.New("incident_id is required")
	}

	if explanation.UserID == "" {
		return model.Explanation{}, errors.New("user_id is required")
	}

	if explanation.Decision == "" {
		return model.Explanation{}, errors.New("decision is required")
	}

	if explanation.CreatedAt.IsZero() {
		explanation.CreatedAt = time.Now().UTC()
	}

	return explanation, nil
}
