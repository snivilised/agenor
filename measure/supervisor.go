package measure

import (
	"github.com/snivilised/traverse/enums"
)

type (
	Metrics  map[enums.Metric]Metric
	Mutables map[enums.Metric]Mutable

	Supervisor struct {
		metrics Metrics
	}
)

func New() *Supervisor {
	return &Supervisor{
		metrics: make(Metrics),
	}
}

func (s *Supervisor) Single(mt enums.Metric) Mutable {
	if _, exists := s.metrics[mt]; !exists {
		metric := &BaseMetric{
			t: mt,
		}

		s.metrics[mt] = metric
		return metric
	}

	return s.metrics[mt].(Mutable)
}

func (s *Supervisor) Many(metrics ...enums.Metric) Mutables {
	result := make(Mutables)

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
