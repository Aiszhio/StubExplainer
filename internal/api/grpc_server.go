package api

import (
	"context"

	"github.com/Aiszhio/StubExplainer/internal/service"
	explainerv1 "github.com/Aiszhio/StubExplainer/pkg/gen/explainer/v1"
)

// GRPCServer implements the generated-like ExplainerServiceServer contract.
type GRPCServer struct {
	service service.Reader
}

func NewGRPCServer(service service.Reader) *GRPCServer {
	return &GRPCServer{service: service}
}

func (s *GRPCServer) GetExplanation(ctx context.Context, req *explainerv1.GetExplanationRequest) (*explainerv1.GetExplanationResponse, error) {
	item, err := s.service.GetExplanation(req.IncidentId)
	if err != nil {
		return nil, err
	}

	return &explainerv1.GetExplanationResponse{
		Explanation: ToProto(item),
	}, nil
}
