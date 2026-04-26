package service

import (
	"github.com/Aiszhio/StubExplainer/internal/ingest"
	"github.com/Aiszhio/StubExplainer/internal/model"
	"github.com/Aiszhio/StubExplainer/internal/storage"
	"github.com/Aiszhio/StubExplainer/internal/validate"
)

// ExplainerService coordinates message processing and explanation retrieval.
type ExplainerService struct {
	storage storage.ExplanationStorage
}

func NewExplainerService(storage storage.ExplanationStorage) *ExplainerService {
	return &ExplainerService{
		storage: storage,
	}
}

func (s *ExplainerService) ProcessBatch(messages [][]byte) error {
	explanations := make([]model.Explanation, 0, len(messages))

	for _, message := range messages {
		parsed, err := ingest.ParseMessage(message)
		if err != nil {
			return err
		}

		validated, err := validate.ValidateAndInit(parsed)
		if err != nil {
			return err
		}

		explanations = append(explanations, validated)
	}

	return s.storage.SaveBatch(explanations)
}

func (s *ExplainerService) GetExplanation(incidentID string) (model.Explanation, error) {
	return s.storage.GetByIncidentID(incidentID)
}
