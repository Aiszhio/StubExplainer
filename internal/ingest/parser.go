package ingest

import (
	"encoding/json"
	"errors"

	"github.com/Aiszhio/StubExplainer/internal/model"
)

// ParseMessage deserializes a raw queue message into an Explanation object.
func ParseMessage(data []byte) (model.Explanation, error) {
	if len(data) == 0 {
		return model.Explanation{}, errors.New("empty message")
	}

	var explanation model.Explanation
	if err := json.Unmarshal(data, &explanation); err != nil {
		return model.Explanation{}, err
	}

	return explanation, nil
}
