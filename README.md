# StubExplainer

Minimal Go implementation of an explainer service for demonstration purposes.

The service shows the main functional flow:

1. Read raw explanation messages from Kafka.
2. Accumulate messages into a batch.
3. Deserialize message data.
4. Run primary validation and initialization.
5. Save explanations into ClickHouse.
6. Read an explanation by incident ID through the HTTP gateway.

## Project structure

```text
api/proto                       gRPC/Protobuf contract
cmd/demo                        In-memory demo entrypoint
cmd/explainer                   Kafka + ClickHouse + HTTP service entrypoint
internal/api                    gRPC adapter and HTTP gateway
internal/consumer               Kafka consumption runner
internal/ingest                 Message parsing
internal/kafka                  Low-level Kafka reader wrapper
internal/model                  Domain models
internal/service                Service layer
internal/storage                Storage interfaces and implementations
internal/validate               Validation and initialization
migrations                      ClickHouse migrations
pkg/config                      Runtime config
pkg/gen                         Generated-like API DTOs and gRPC interfaces
pkg/logger                      Logger helper
```

## Run tests

```bash
go mod tidy
go test ./...
```

Or with Makefile:

```bash
make tidy
make test
```

## Run in-memory demo

```bash
go run ./cmd/demo
```

Or:

```bash
make run-demo
```

## Run infrastructure and service

```bash
docker compose up --build
```

Or:

```bash
make compose-up
```

The HTTP gateway is exposed on port `8080`.

## Send test message to Kafka

```bash
make produce
```

Manual command:

```bash
printf '%s\n' '{"incident_id":"inc_1001","user_id":"user_501","decision":"blocked","score":0.97,"features":[{"name":"request_frequency","value":"high","weight":0.42},{"name":"device_fingerprint","value":"suspicious","weight":0.31}]}' | docker compose exec -T kafka kafka-console-producer --bootstrap-server kafka:29092 --topic explanations
```

## Check data in ClickHouse

```bash
make check-clickhouse
```

Manual command:

```bash
docker compose exec clickhouse clickhouse-client \
  --database explainer \
  --query "SELECT incident_id, user_id, decision, score, features_json, created_at FROM explanations FORMAT Vertical"
```

## Request explanation through HTTP gateway

After sending the test message, wait a few seconds and run:

```bash
curl http://localhost:8080/v1/explanations/inc_1001
```

Expected response format:

```json
{
  "explanation": {
    "incident_id": "inc_1001",
    "user_id": "user_501",
    "decision": "blocked",
    "score": 0.97,
    "features": [
      {
        "name": "request_frequency",
        "value": "high",
        "weight": 0.42
      },
      {
        "name": "device_fingerprint",
        "value": "suspicious",
        "weight": 0.31
      }
    ],
    "created_at": "2026-04-26T10:15:30Z"
  }
}
```

## Configuration

| Variable | Default | Description |
| --- | --- | --- |
| `KAFKA_BROKERS` | `localhost:9092` | Kafka broker list separated by commas |
| `KAFKA_TOPIC` | `explanations` | Topic with explanation messages |
| `KAFKA_GROUP_ID` | `stub-explainer` | Kafka consumer group ID |
| `CLICKHOUSE_ADDR` | `localhost:9000` | ClickHouse native protocol address |
| `CLICKHOUSE_DATABASE` | `default` | ClickHouse database name |
| `CLICKHOUSE_USERNAME` | `default` | ClickHouse username |
| `CLICKHOUSE_PASSWORD` | empty | ClickHouse password |
| `HTTP_ADDR` | `:8080` | HTTP gateway address |
| `BATCH_SIZE` | `100` | Maximum batch size |
| `BATCH_INTERVAL_MS` | `10` | Batch flush interval in milliseconds |
