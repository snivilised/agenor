package measure

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
)

// 📦 pkg: measure - defines facilities for counting things
// represented by metrics.

type (
	// MutableMetric represents write access to the metric
	MutableMetric interface {
		core.Metric
		Tick() core.MetricValue
		Times(increment uint) core.MetricValue
	}

	BaseMetric struct {
		t       enums.Metric
		counter core.MetricValue
	}
)

func (m *BaseMetric) Type() enums.Metric {
	return m.t
}

func (m *BaseMetric) Value() core.MetricValue {
	return m.counter
}

func (m *BaseMetric) Tick() core.MetricValue {
	m.counter++

	return m.counter
}

func (m *BaseMetric) Times(increment uint) core.MetricValue {
	m.counter += increment

	return m.counter
}
