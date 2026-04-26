package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aiszhio/StubExplainer/internal/api"
	"github.com/Aiszhio/StubExplainer/internal/consumer"
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
	kafkaConsumer := reader.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID)
	defer kafkaConsumer.Close()

	consumerRunner := consumer.NewRunner(kafkaConsumer, explainerService, cfg, logg)
	go consumerRunner.Run(ctx)

	grpcServer := api.NewGRPCServer(explainerService)
	grpcClient := api.NewLocalClient(grpcServer)
	httpGateway := api.NewHTTPGateway(grpcClient)

	mux := http.NewServeMux()
	httpGateway.Register(mux)

	httpServer := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logg.Printf("http gateway started: addr=%s", cfg.HTTPAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Printf("http gateway error: %v", err)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logg.Printf("http gateway shutdown error: %v", err)
	}

	logg.Println("service stopped")
}
