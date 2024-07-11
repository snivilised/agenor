package measure

import (
	"github.com/snivilised/traverse/enums"
)

type (
	Metrics        map[enums.Metric]Metric
	MutableMetrics map[enums.Metric]MutableMetric

	Supervisor struct {
		metrics Metrics
	}

	Owned struct {
		Mums MutableMetrics
	}
)

func New() *Supervisor {
	return &Supervisor{
		metrics: make(Metrics),
	}
}

func (s *Supervisor) Single(mt enums.Metric) MutableMetric {
	if _, exists := s.metrics[mt]; !exists {
		metric := &BaseMetric{
			t: mt,
		}

		s.metrics[mt] = metric
		return metric
	}

	return s.metrics[mt].(MutableMetric)
}

func (s *Supervisor) Many(metrics ...enums.Metric) MutableMetrics {
	result := make(MutableMetrics)

	for _, mt := range metrics {
		metric := &BaseMetric{
			t: mt,
		}
		if _, exists := s.metrics[mt]; !exists {
			s.metrics[mt] = metric
		}
		result[mt] = metric
	}

	return result
}

func (s *Supervisor) Count(mt enums.Metric) MetricValue {
	if metric, exists := s.metrics[mt]; exists {
		return metric.Value()
	}

	return 0
}
