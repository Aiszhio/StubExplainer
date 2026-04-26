package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contains runtime service settings.
type Config struct {
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string

	ClickHouseAddr     string
	ClickHouseDatabase string
	ClickHouseUsername string
	ClickHousePassword string

	GRPCAddr string
	HTTPAddr string

	BatchSize     int
	BatchInterval time.Duration
}

func Load() Config {
	return Config{
		KafkaBrokers:       splitEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopic:         getEnv("KAFKA_TOPIC", "explanations"),
		KafkaGroupID:       getEnv("KAFKA_GROUP_ID", "stub-explainer"),
		ClickHouseAddr:     getEnv("CLICKHOUSE_ADDR", "localhost:9000"),
		ClickHouseDatabase: getEnv("CLICKHOUSE_DATABASE", "default"),
		ClickHouseUsername: getEnv("CLICKHOUSE_USERNAME", "default"),
		ClickHousePassword: getEnv("CLICKHOUSE_PASSWORD", ""),
		GRPCAddr:           getEnv("GRPC_ADDR", ":50051"),
		HTTPAddr:           getEnv("HTTP_ADDR", ":8080"),
		BatchSize:          intEnv("BATCH_SIZE", 100),
		BatchInterval:      time.Duration(intEnv("BATCH_INTERVAL_MS", 10)) * time.Millisecond,
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func intEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func splitEnv(key, fallback string) []string {
	value := getEnv(key, fallback)
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
