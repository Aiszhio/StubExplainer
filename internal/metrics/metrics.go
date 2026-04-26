package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics contains service metrics and hides Prometheus implementation details from business code.
type Metrics struct {
	registry *prometheus.Registry

	processedMessages prometheus.Counter
	failedMessages    prometheus.Counter
	savedBatches      prometheus.Counter
	savedMessages     prometheus.Counter
	kafkaReadErrors   prometheus.Counter
	processErrors     prometheus.Counter
	httpRequests      *prometheus.CounterVec
	batchSize         prometheus.Histogram
	batchDuration     prometheus.Histogram
	httpDuration      *prometheus.HistogramVec
}

// Recorder is the metric interface used by internal components.
type Recorder interface {
	ObserveKafkaReadError()
	ObserveProcessError()
	ObserveMessageProcessed()
	ObserveMessageFailed()
	ObserveBatchSaved(size int, duration time.Duration)
	ObserveHTTPRequest(method string, path string, statusCode string, duration time.Duration)
	Handler() http.Handler
}

func New() *Metrics {
	registry := prometheus.NewRegistry()

	m := &Metrics{
		registry: registry,
		processedMessages: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "explainer_processed_messages_total",
			Help: "Total number of successfully processed messages.",
		}),
		failedMessages: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "explainer_failed_messages_total",
			Help: "Total number of failed messages.",
		}),
		savedBatches: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "explainer_saved_batches_total",
			Help: "Total number of saved batches.",
		}),
		savedMessages: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "explainer_saved_messages_total",
			Help: "Total number of saved messages.",
		}),
		kafkaReadErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "explainer_kafka_read_errors_total",
			Help: "Total number of Kafka read errors.",
		}),
		processErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "explainer_process_errors_total",
			Help: "Total number of batch processing errors.",
		}),
		httpRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "explainer_http_requests_total",
			Help: "Total number of HTTP gateway requests.",
		}, []string{"method", "path", "status"}),
		batchSize: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "explainer_batch_size",
			Help:    "Observed Kafka batch sizes.",
			Buckets: []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000},
		}),
		batchDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "explainer_batch_processing_duration_seconds",
			Help:    "Batch processing duration in seconds.",
			Buckets: prometheus.DefBuckets,
		}),
		httpDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "explainer_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path", "status"}),
	}

	registry.MustRegister(
		m.processedMessages,
		m.failedMessages,
		m.savedBatches,
		m.savedMessages,
		m.kafkaReadErrors,
		m.processErrors,
		m.httpRequests,
		m.batchSize,
		m.batchDuration,
		m.httpDuration,
	)

	return m
}

func (m *Metrics) ObserveKafkaReadError() {
	m.kafkaReadErrors.Inc()
}

func (m *Metrics) ObserveProcessError() {
	m.processErrors.Inc()
}

func (m *Metrics) ObserveMessageProcessed() {
	m.processedMessages.Inc()
}

func (m *Metrics) ObserveMessageFailed() {
	m.failedMessages.Inc()
}

func (m *Metrics) ObserveBatchSaved(size int, duration time.Duration) {
	m.savedBatches.Inc()
	m.savedMessages.Add(float64(size))
	m.batchSize.Observe(float64(size))
	m.batchDuration.Observe(duration.Seconds())
}

func (m *Metrics) ObserveHTTPRequest(method string, path string, statusCode string, duration time.Duration) {
	m.httpRequests.WithLabelValues(method, path, statusCode).Inc()
	m.httpDuration.WithLabelValues(method, path, statusCode).Observe(duration.Seconds())
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
