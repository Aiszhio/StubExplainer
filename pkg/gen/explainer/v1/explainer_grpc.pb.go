// Code generated manually for the educational stub. DO NOT EDIT.
package explainerv1

import "context"

// ExplainerServiceServer is the generated-like server interface used by adapters.
type ExplainerServiceServer interface {
	GetExplanation(ctx context.Context, req *GetExplanationRequest) (*GetExplanationResponse, error)
}

// ExplainerServiceClient is the generated-like client interface used by the HTTP gateway.
type ExplainerServiceClient interface {
	GetExplanation(ctx context.Context, req *GetExplanationRequest) (*GetExplanationResponse, error)
}
