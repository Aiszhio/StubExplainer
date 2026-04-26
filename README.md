# StubExplainer

Minimal Go implementation of an explainer service for demonstration purposes.

The service shows the main functional flow:

1. Accumulate a batch of raw messages.
2. Deserialize message data.
3. Run primary validation and initialization.
4. Save explanations into storage.
5. Read an explanation by incident ID.

The current storage implementation is in-memory and is intended only for demonstration. In the target system it can be replaced with ClickHouse storage without changing the service flow.

## Project structure

```text
cmd/explainer/main.go          Demo entrypoint
internal/model                 Domain models
internal/ingest                Message parsing
internal/validate              Validation and initialization
internal/storage               Storage implementation
internal/service               Service layer
```

## Run

```bash
go run ./cmd/explainer
```

## Test

```bash
go test ./...
```
