// Code generated manually for the educational stub. DO NOT EDIT.
package explainerv1

// GetExplanationRequest describes a request for one explanation by incident ID.
type GetExplanationRequest struct {
	IncidentId string `json:"incident_id"`
}

// GetExplanationResponse contains the found explanation.
type GetExplanationResponse struct {
	Explanation *Explanation `json:"explanation,omitempty"`
}

// Explanation is the API representation of a stored explanation.
type Explanation struct {
	IncidentId string           `json:"incident_id"`
	UserId     string           `json:"user_id"`
	Decision   string           `json:"decision"`
	Score      float64          `json:"score"`
	Features   []*FeatureImpact `json:"features"`
	CreatedAt  string           `json:"created_at"`
}

// FeatureImpact describes one factor that affected a model decision.
type FeatureImpact struct {
	Name   string  `json:"name"`
	Value  string  `json:"value"`
	Weight float64 `json:"weight"`
}
