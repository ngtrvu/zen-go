package http

import (
	"fmt"
	"time"

	"github.com/ngtrvu/zen-go/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	LabelNamespace    = "namespace"
	LabelEndpointName = "endpoint"
	LabelStatusCode   = "status_code"
	LabelMethod       = "method"
	LabelSuccess      = "success"
	LabelHost         = "host"
)

type Metrics struct {
	Namespace              string
	AppName                string
	InBoundRequestDuration *metrics.Histogram
}

type OutboutMetrics struct {
	Namespace               string
	AppName                 string
	OutBoundRequestDuration *metrics.Histogram
}

func NewHTTPMetrics(appName string) *Metrics {
	return &Metrics{
		AppName:                appName,
		InBoundRequestDuration: newInBoundRequestDurationHistogram(),
	}
}

func NewHTTPOutboundMetrics(appName string) *OutboutMetrics {
	return &OutboutMetrics{
		AppName:                 appName,
		OutBoundRequestDuration: newOutBoundRequestDurationHistogram(),
	}
}

func newInBoundRequestDurationHistogram() *metrics.Histogram {
	durationBuckets := []float64{
		0.005,
		0.01, 0.02, 0.05,
		0.1, 0.2, 0.5,
		1, 2, 5,
		10, 30, 60,
		120, 180, 300, 600,
	}
	requestDuration := metrics.NewHistogramFrom(prometheus.HistogramOpts{
		Name:    "http_server_handle_request_duration",
		Help:    "inbound request duration in second",
		Buckets: durationBuckets,
	}, []string{LabelEndpointName, LabelStatusCode, LabelMethod})
	return requestDuration
}

func newOutBoundRequestDurationHistogram() *metrics.Histogram {
	durationBuckets := []float64{
		0.005,
		0.01, 0.02, 0.05,
		0.1, 0.2, 0.5,
		1, 2, 5,
		10, 30, 60,
		120, 180, 300, 600,
	}
	requestDuration := metrics.NewHistogramFrom(prometheus.HistogramOpts{
		Name:    "http_out_request_duration",
		Help:    "outbound request duration in second",
		Buckets: durationBuckets,
	}, []string{LabelHost, LabelMethod, LabelStatusCode, LabelSuccess})
	return requestDuration
}

// ObserveInBoundRequest observe inbound request duration in second
func (m *Metrics) ObserveRequest(start time.Time, statusCode int, method, endpoint string) {
	m.InBoundRequestDuration.With(
		LabelEndpointName, endpoint,
		LabelMethod, method,
		LabelStatusCode, fmt.Sprintf("%v", statusCode),
	).Observe(time.Since(start).Seconds())
}

// ObserveInBoundRequest observe inbound request duration in second
func (m *OutboutMetrics) ObserveRequest(start time.Time, success bool, statusCode int, method, host string) {
	m.OutBoundRequestDuration.With(
		LabelHost, host,
		LabelMethod, method,
		LabelStatusCode, fmt.Sprintf("%v", statusCode),
		LabelSuccess, fmt.Sprintf("%v", success),
	).Observe(time.Since(start).Seconds())
}
