package enclave

import (
	"maps"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
)

type (
	Supervisor struct {
		metrics core.Metrics
	}

	Crate struct {
		Metrics core.Metrics
	}
)

func NewSupervisor() *Supervisor {
	return &Supervisor{
		metrics: make(core.Metrics),
	}
}

func (s *Supervisor) Load(metrics core.Metrics) {
	s.metrics = maps.Clone(metrics)
}

func (s *Supervisor) Single(mt enums.Metric) *core.NavigationMetric {
	if _, exists := s.metrics[mt]; !exists {
		metric := &core.NavigationMetric{
			T: mt,
		}

		s.metrics[mt] = metric

		return metric
	}

	return s.metrics[mt]
}

func (s *Supervisor) Many(metrics ...enums.Metric) core.Metrics {
	result := make(core.Metrics)

	for _, mt := range metrics {
		result[mt] = s.Single(mt)
	}

	return result
}

func (s *Supervisor) Count(mt enums.Metric) core.MetricValue {
	if metric, exists := s.metrics[mt]; exists {
		return metric.Value()
	}

	return 0
}
