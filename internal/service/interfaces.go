package service

import (
	"github.com/Aiszhio/StubExplainer/internal/model"
)

// Processor describes write-path use cases.
type Processor interface {
	ProcessBatch(messages [][]byte) error
}

// Reader describes read-path use cases.
type Reader interface {
	GetExplanation(incidentID string) (model.Explanation, error)
}

// Explainer describes all service use cases.
type Explainer interface {
	Processor
	Reader
}
