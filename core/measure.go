package core

import (
	"github.com/snivilised/traverse/enums"
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

	// Reporter represents query access to the metrics Supervisor
	Reporter interface {
		Count(enums.Metric) MetricValue
	}
)
