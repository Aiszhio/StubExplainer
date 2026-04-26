package api

import (
	"context"

	explainerv1 "github.com/Aiszhio/StubExplainer/pkg/gen/explainer/v1"
)

// LocalClient adapts the server implementation to the generated-like client interface.
// It is used by the HTTP gateway in the educational stub.
type LocalClient struct {
	server explainerv1.ExplainerServiceServer
}

func NewLocalClient(server explainerv1.ExplainerServiceServer) *LocalClient {
	return &LocalClient{server: server}
}

func (c *LocalClient) GetExplanation(ctx context.Context, req *explainerv1.GetExplanationRequest) (*explainerv1.GetExplanationResponse, error) {
	return c.server.GetExplanation(ctx, req)
}
