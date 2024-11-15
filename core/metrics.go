package core

import (
	"github.com/snivilised/agenor/enums"
)

type (
	MetricValue = uint

	// Metric represents query access to the metric. The client
	// registering the metric should maintain it's mutate access
	// to the metric so they can update it.
	Metric interface {
		Type() enums.Metric
		Value() MetricValue
	}

	// MutableMetric represents write access to the metric
	MutableMetric interface {
		Metric
		Tick() MetricValue
		Times(increment uint) MetricValue
	}

	NavigationMetric struct {
		T       enums.Metric
		Counter MetricValue
	}

	Metrics map[enums.Metric]*NavigationMetric

	// Reporter represents query access to the metrics Supervisor
	Reporter interface {
		Count(enums.Metric) MetricValue
	}
)

func (m *NavigationMetric) Type() enums.Metric {
	return m.T
}

func (m *NavigationMetric) Value() MetricValue {
	return m.Counter
}

func (m *NavigationMetric) Tick() MetricValue {
	m.Counter++

	return m.Counter
}

func (m *NavigationMetric) Times(increment uint) MetricValue {
	m.Counter += increment

	return m.Counter
}
