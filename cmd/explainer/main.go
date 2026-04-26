package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Aiszhio/StubExplainer/internal/service"
	"github.com/Aiszhio/StubExplainer/internal/storage"
)

func main() {
	storage := storage.NewMemoryStorage()
	explainerService := service.NewExplainerService(storage)

	message := []byte(`{
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
		]
	}`)

	if err := explainerService.ProcessBatch([][]byte{message}); err != nil {
		log.Fatal(err)
	}

	result, err := explainerService.GetExplanation("inc_1001")
	if err != nil {
		log.Fatal(err)
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))
}
