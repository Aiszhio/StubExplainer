package api

import (
	"time"

	"github.com/Aiszhio/StubExplainer/internal/model"
	explainerv1 "github.com/Aiszhio/StubExplainer/pkg/gen/explainer/v1"
)

func ToProto(item model.Explanation) *explainerv1.Explanation {
	features := make([]*explainerv1.FeatureImpact, 0, len(item.Features))
	for _, feature := range item.Features {
		features = append(features, &explainerv1.FeatureImpact{
			Name:   feature.Name,
			Value:  feature.Value,
			Weight: feature.Weight,
		})
	}

	return &explainerv1.Explanation{
		IncidentId: item.IncidentID,
		UserId:     item.UserID,
		Decision:   item.Decision,
		Score:      item.Score,
		Features:   features,
		CreatedAt:  item.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}
