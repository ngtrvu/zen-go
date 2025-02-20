package grpc

import (
	"time"

	"github.com/ngtrvu/zen-go/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	LabelNamespace   = "namespace"
	LabelServiceName = "service_name"
	LabelGrpcType    = "grpc_type"
	LabelMethod      = "method"
	LabelStatusCode  = "status_code"
)

type Metrics struct {
	Namespace       string
	AppName         string
	RequestDuration *metrics.Histogram
}

func NewMetrics(namespace, appName string) *Metrics {
	return &Metrics{
		Namespace:       namespace,
		AppName:         appName,
		RequestDuration: newRequestDurationHistogram(appName),
	}
}

func newRequestDurationHistogram(appName string) *metrics.Histogram {
	durationBuckets := []float64{
		0.00001, 0.00002, 0.00005,
		0.0001, 0.0002, 0.0005,
		0.001, 0.002, 0.005,
		0.01, 0.02, 0.05,
		0.1, 0.2, 0.5,
		1, 2, 5,
		10, 30, 60,
		120, 180, 300, 600,
	}
	requestDuration := metrics.NewHistogramFrom(prometheus.HistogramOpts{
		Name:    "grpc_server_handle_request_duration",
		Help:    "Grpc server handle request duration in second",
		Buckets: durationBuckets,
	}, []string{LabelServiceName, LabelGrpcType, LabelMethod, LabelStatusCode})
	return requestDuration
}

// ObserveRequest observe request duration in second
func (m *Metrics) ObserveRequest(start time.Time, serviceName, grpcType, method, statusCode string) {
	m.RequestDuration.With(
		LabelServiceName, serviceName,
		LabelGrpcType, grpcType,
		LabelMethod, method,
		LabelStatusCode, statusCode,
	).Observe(time.Since(start).Seconds())
}
