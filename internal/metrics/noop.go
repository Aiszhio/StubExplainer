package metrics

import (
	"net/http"
	"time"
)

// Noop is used in unit tests and demos where metrics collection is not required.
type Noop struct{}

func (Noop) ObserveKafkaReadError() {}
func (Noop) ObserveProcessError() {}
func (Noop) ObserveMessageProcessed() {}
func (Noop) ObserveMessageFailed() {}
func (Noop) ObserveBatchSaved(size int, duration time.Duration) {}
func (Noop) ObserveHTTPRequest(method string, path string, statusCode string, duration time.Duration) {}
func (Noop) Handler() http.Handler { return http.NotFoundHandler() }
