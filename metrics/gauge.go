package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Gauge struct {
	GV  *prometheus.GaugeVec
	LVS LabelValues
}

func NewGaugeFrom(opts prometheus.GaugeOpts, labelNames []string) *Gauge {
	gv := promauto.NewGaugeVec(opts, labelNames)

	return &Gauge{
		GV: gv,
	}
}

func (g *Gauge) With(labelValues ...string) *Gauge {
	return &Gauge{
		GV:  g.GV,
		LVS: g.LVS.With(labelValues...),
	}
}

// Set adds the given value to the gauge
func (g *Gauge) Set(value float64) {
	labels := g.LVS.BuildPrometheusLabels()
	g.GV.With(labels).Set(value)
}
