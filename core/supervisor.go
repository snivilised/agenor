package core

import (
	"github.com/snivilised/traverse/enums"
)

type (
	Supervisor struct {
		metrics Metrics
	}

	Crate struct {
		Metrics Metrics
	}
)

func NewSupervisor() *Supervisor {
	return &Supervisor{
		metrics: make(Metrics),
	}
}

func (s *Supervisor) Single(mt enums.Metric) *NavigationMetric {
	if _, exists := s.metrics[mt]; !exists {
		metric := &NavigationMetric{
			T: mt,
		}

		s.metrics[mt] = metric

		return metric
	}

	return s.metrics[mt]
}

func (s *Supervisor) Many(metrics ...enums.Metric) Metrics {
	result := make(Metrics)

	for _, mt := range metrics {
		result[mt] = s.Single(mt)
	}

	return result
}

func (s *Supervisor) Count(mt enums.Metric) MetricValue {
	if metric, exists := s.metrics[mt]; exists {
		return metric.Value()
	}

	return 0
}
