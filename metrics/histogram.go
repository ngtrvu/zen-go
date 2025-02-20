package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Histogram struct {
	HV  *prometheus.HistogramVec
	LVS LabelValues
}

func NewHistogramFrom(opts prometheus.HistogramOpts, labelNames []string) *Histogram {
	hv := promauto.NewHistogramVec(opts, labelNames)
	return &Histogram{
		HV: hv,
	}
}

func (h *Histogram) With(labelValues ...string) *Histogram {
	return &Histogram{
		HV:  h.HV,
		LVS: h.LVS.With(labelValues...),
	}
}

func (h *Histogram) Observe(value float64) {
	labels := h.LVS.BuildPrometheusLabels()
	h.HV.With(labels).Observe(value)
}
