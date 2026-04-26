package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	reader "github.com/Aiszhio/StubExplainer/internal/kafka"
	"github.com/Aiszhio/StubExplainer/internal/service"
	"github.com/Aiszhio/StubExplainer/internal/storage"
	"github.com/Aiszhio/StubExplainer/pkg/config"
	"github.com/Aiszhio/StubExplainer/pkg/logger"
)

func main() {
	logg := logger.New()
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	clickHouseStorage, err := storage.NewClickHouseStorage(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to connect clickhouse: %v", err)
	}
	defer clickHouseStorage.Close()

	explainerService := service.NewExplainerService(clickHouseStorage)
	consumer := reader.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID)
	defer consumer.Close()

	logg.Printf("service started: topic=%s brokers=%v batch_size=%d batch_interval=%s", cfg.KafkaTopic, cfg.KafkaBrokers, cfg.BatchSize, cfg.BatchInterval)

	for {
		select {
		case <-ctx.Done():
			logg.Println("service stopped")
			return
		default:
		}

		messages, err := consumer.ReadBatch(ctx, cfg.BatchSize, cfg.BatchInterval)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				continue
			}
			logg.Printf("failed to read kafka batch: %v", err)
			time.Sleep(time.Second)
			continue
		}

		if len(messages) == 0 {
			continue
		}

		if err := explainerService.ProcessBatch(messages); err != nil {
			logg.Printf("failed to process batch: %v", err)
			continue
		}

		logg.Printf("saved batch: size=%d", len(messages))
	}
}
