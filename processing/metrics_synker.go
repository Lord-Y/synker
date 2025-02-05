// Package processing provide all requirements to process change data capture
package processing

import (
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func newMetrics() (m *metrics, err error) {
	name := "synker"
	z := &metrics{
		kafka: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: name,
				Subsystem: "kafka",
				Name:      "errors_total",
				Help:      "Number of errors related to kafka messages",
			},
			[]string{"error_type", "topic"},
		),
		elasticsearch: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: name,
				Subsystem: "elaticsearch",
				Name:      "errors_total",
				Help:      "Number of errors related to elasticsearch",
			},
			[]string{"error_type", "topic"},
		),
	}
	if err := prometheus.Register(z.kafka); err != nil {
		_, ok := err.(prometheus.AlreadyRegisteredError)
		if !ok {
			return nil, err
		}
	}
	if err := prometheus.Register(z.elasticsearch); err != nil {
		_, ok := err.(prometheus.AlreadyRegisteredError)
		if !ok {
			return nil, err
		}
	}
	return z, nil
}

func (c *Validate) increaseMetrics(metricType, topic, errorType string) {
	if strings.TrimSpace(os.Getenv("SYNKER_PROMETHEUS")) != "" {
		switch metricType {
		case "kafka":
			c.metrics.kafka.With(prometheus.Labels{
				"topic":      topic,
				"error_type": errorType,
			}).Inc()
		case "elasticsearch":
			c.metrics.elasticsearch.With(prometheus.Labels{
				"topic":      topic,
				"error_type": errorType,
			}).Inc()
		}
	}
}
