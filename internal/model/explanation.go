package model

import "time"

// FeatureImpact describes one factor that affected a model decision.
type FeatureImpact struct {
	Name   string  `json:"name"`
	Value  string  `json:"value"`
	Weight float64 `json:"weight"`
}

// Explanation is the main domain object stored and returned by the service.
type Explanation struct {
	IncidentID string          `json:"incident_id"`
	UserID     string          `json:"user_id"`
	Decision   string          `json:"decision"`
	Score      float64         `json:"score"`
	Features   []FeatureImpact `json:"features"`
	CreatedAt  time.Time       `json:"created_at"`
}
