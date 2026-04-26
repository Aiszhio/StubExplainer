package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// Consumer wraps Kafka reader logic for reading explanation messages.
type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			CommitInterval: time.Second,
			MinBytes:       1,
			MaxBytes:       10e6,
		}),
	}
}

func (c *Consumer) ReadBatch(ctx context.Context, maxSize int, timeout time.Duration) ([][]byte, error) {
	if maxSize <= 0 {
		maxSize = 1
	}

	batch := make([][]byte, 0, maxSize)
	deadlineCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for len(batch) < maxSize {
		message, err := c.reader.FetchMessage(deadlineCtx)
		if err != nil {
			if len(batch) > 0 && deadlineCtx.Err() != nil {
				return batch, nil
			}
			return batch, err
		}

		batch = append(batch, message.Value)

		if err := c.reader.CommitMessages(ctx, message); err != nil {
			return batch, err
		}
	}

	return batch, nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
