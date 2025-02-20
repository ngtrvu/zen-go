package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricDataType string
type MetricCounterStatus string

type Counter struct {
	CV  *prometheus.CounterVec
	LVS LabelValues
}

func NewCounterFrom(opts prometheus.CounterOpts, labelNames []string) *Counter {
	cv := promauto.NewCounterVec(opts, labelNames)

	return &Counter{
		CV: cv,
	}
}

func (c *Counter) With(labelValues ...string) *Counter {
	return &Counter{
		CV:  c.CV,
		LVS: c.LVS.With(labelValues...),
	}
}

// Add adds the given value to the counter
func (c *Counter) Add(delta float64) {
	labels := c.LVS.BuildPrometheusLabels()
	c.CV.With(labels).Add(delta)
}

// Inc increments the counter by 1
func (c *Counter) Inc() {
	labels := c.LVS.BuildPrometheusLabels()
	c.CV.With(labels).Inc()
}
