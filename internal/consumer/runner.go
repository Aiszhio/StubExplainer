package consumer

import (
	"context"
	"errors"
	"log"
	"time"

	reader "github.com/Aiszhio/StubExplainer/internal/kafka"
	"github.com/Aiszhio/StubExplainer/internal/metrics"
	"github.com/Aiszhio/StubExplainer/internal/service"
	"github.com/Aiszhio/StubExplainer/pkg/config"
)

// Runner owns the background Kafka consumption loop.
type Runner struct {
	consumer reader.BatchReader
	service  service.Processor
	metrics  metrics.Recorder
	cfg      config.Config
	logg     *log.Logger
}

func NewRunner(consumer reader.BatchReader, service service.Processor, metrics metrics.Recorder, cfg config.Config, logg *log.Logger) *Runner {
	return &Runner{
		consumer: consumer,
		service:  service,
		metrics:  metrics,
		cfg:      cfg,
		logg:     logg,
	}
}

func (r *Runner) Run(ctx context.Context) {
	r.logg.Printf("consumer started: topic=%s brokers=%v batch_size=%d batch_interval=%s", r.cfg.KafkaTopic, r.cfg.KafkaBrokers, r.cfg.BatchSize, r.cfg.BatchInterval)

	for {
		select {
		case <-ctx.Done():
			r.logg.Println("consumer stopped")
			return
		default:
		}

		messages, err := r.consumer.ReadBatch(ctx, r.cfg.BatchSize, r.cfg.BatchInterval)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				continue
			}
			r.metrics.ObserveKafkaReadError()
			r.logg.Printf("failed to read kafka batch: %v", err)
			time.Sleep(time.Second)
			continue
		}

		if len(messages) == 0 {
			continue
		}

		startedAt := time.Now()
		if err := r.service.ProcessBatch(messages); err != nil {
			r.metrics.ObserveProcessError()
			r.logg.Printf("failed to process batch: %v", err)
			continue
		}

		r.metrics.ObserveBatchSaved(len(messages), time.Since(startedAt))
		r.logg.Printf("saved batch: size=%d", len(messages))
	}
}
