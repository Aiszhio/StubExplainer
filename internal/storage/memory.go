package storage

import (
	"errors"
	"sync"

	"github.com/Aiszhio/StubExplainer/internal/model"
)

// MemoryStorage is a demonstration in-memory storage implementation.
type MemoryStorage struct {
	mu   sync.RWMutex
	data map[string]model.Explanation
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]model.Explanation),
	}
}

func (s *MemoryStorage) SaveBatch(items []model.Explanation) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, item := range items {
		s.data[item.IncidentID] = item
	}

	return nil
}

func (s *MemoryStorage) GetByIncidentID(incidentID string) (model.Explanation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.data[incidentID]
	if !ok {
		return model.Explanation{}, errors.New("explanation not found")
	}

	return item, nil
}
