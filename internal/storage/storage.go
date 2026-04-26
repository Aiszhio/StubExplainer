package storage

import "github.com/Aiszhio/StubExplainer/internal/model"

// ExplanationStorage describes storage operations required by the service layer.
type ExplanationStorage interface {
	SaveBatch(items []model.Explanation) error
	GetByIncidentID(incidentID string) (model.Explanation, error)
}
