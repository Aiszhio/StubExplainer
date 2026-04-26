package kafka

import (
	"context"
	"time"
)

// BatchReader describes a source that can return raw messages in batches.
type BatchReader interface {
	ReadBatch(ctx context.Context, maxSize int, timeout time.Duration) ([][]byte, error)
	Close() error
}
