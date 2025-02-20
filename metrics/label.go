package metrics

import "github.com/prometheus/client_golang/prometheus"

type LabelValues []string

func (lvs LabelValues) With(labelValues ...string) LabelValues {
	if len(labelValues)%2 != 0 {
		labelValues = append(labelValues, "unknown")
	}
	return append(lvs, labelValues...)
}

func (lvs LabelValues) BuildPrometheusLabels() prometheus.Labels {
	labels := prometheus.Labels{}
	for i := 0; i < len(lvs); i += 2 {
		labels[lvs[i]] = lvs[i+1]
	}
	return labels
}
