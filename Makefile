.PHONY: test run-demo compose-up compose-down compose-build tidy produce check-clickhouse

APP_NAME=stub-explainer
KAFKA_CONTAINER=stubexplainer-kafka-1
CLICKHOUSE_CONTAINER=stubexplainer-clickhouse-1

# Run unit and integration tests that do not require external services.
test:
	go test ./...

# Run local in-memory demo without Kafka and ClickHouse.
run-demo:
	go run ./cmd/demo

# Download modules and refresh go.sum.
tidy:
	go mod tidy

# Build and start Kafka, ClickHouse and the service.
compose-up:
	docker compose up --build

# Build containers without starting them.
compose-build:
	docker compose build

# Stop local stack.
compose-down:
	docker compose down -v

# Send one demo message to Kafka.
produce:
	printf '%s\n' '{"incident_id":"inc_1001","user_id":"user_501","decision":"blocked","score":0.97,"features":[{"name":"request_frequency","value":"high","weight":0.42},{"name":"device_fingerprint","value":"suspicious","weight":0.31}]}' | docker compose exec -T kafka kafka-console-producer --bootstrap-server kafka:29092 --topic explanations

# Read saved rows from ClickHouse.
check-clickhouse:
	docker compose exec clickhouse clickhouse-client --database explainer --query "SELECT incident_id, user_id, decision, score, features_json, created_at FROM explanations FORMAT Vertical"
