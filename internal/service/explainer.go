package service

import (
	"github.com/Aiszhio/StubExplainer/internal/ingest"
	"github.com/Aiszhio/StubExplainer/internal/metrics"
	"github.com/Aiszhio/StubExplainer/internal/model"
	"github.com/Aiszhio/StubExplainer/internal/storage"
	"github.com/Aiszhio/StubExplainer/internal/validate"
)

// ExplainerService coordinates message processing and explanation retrieval.
type ExplainerService struct {
	storage storage.ExplanationStorage
	metrics metrics.Recorder
}

func NewExplainerService(storage storage.ExplanationStorage, metrics metrics.Recorder) *ExplainerService {
	return &ExplainerService{
		storage: storage,
		metrics: metrics,
	}
}

func (s *ExplainerService) ProcessBatch(messages [][]byte) error {
	explanations := make([]model.Explanation, 0, len(messages))

	for _, message := range messages {
		parsed, err := ingest.ParseMessage(message)
		if err != nil {
			s.metrics.ObserveMessageFailed()
			return err
		}

		validated, err := validate.ValidateAndInit(parsed)
		if err != nil {
			s.metrics.ObserveMessageFailed()
			return err
		}

		s.metrics.ObserveMessageProcessed()
		explanations = append(explanations, validated)
	}

	return s.storage.SaveBatch(explanations)
}

func (s *ExplainerService) GetExplanation(incidentID string) (model.Explanation, error) {
	return s.storage.GetByIncidentID(incidentID)
}
