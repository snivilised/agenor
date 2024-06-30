package measure

import (
	"github.com/snivilised/traverse/enums"
)

// the measure package defines facilities for counting things
// represented by metrics.

type (
	MetricValue = uint

	// Metric represents query access to the metric. The client
	// registering the metric should maintain it's mutate access
	// to the metric so they can update it.
	Metric interface {
		Type() enums.Metric
		Value() MetricValue
	}

	// Mutable represents write access to the metric
	Mutable interface {
		Metric
		Tick() MetricValue
	}

	// Reporter represents query access to the metrics Supervisor
	Reporter interface {
		Count(enums.Metric) MetricValue
	}

	BaseMetric struct {
		t       enums.Metric
		counter MetricValue
	}
)

func (m *BaseMetric) Type() enums.Metric {
	return m.t
}

func (m *BaseMetric) Value() MetricValue {
	return m.counter
}

func (m *BaseMetric) Tick() MetricValue {
	m.counter++

	return m.counter
}
