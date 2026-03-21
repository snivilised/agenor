package core

import (
	"maps"

	"github.com/snivilised/jaywalk/src/agenor/enums"
)

type (
	// MetricValue represents the value of a metric
	MetricValue = uint

	// Metric represents query access to the metric. The client
	// registering the metric should maintain it's mutate access
	// to the metric so they can update it.
	Metric interface {
		// Type returns the type of the metric, which can be used to identify
		// and categorize the metric for reporting and analysis purposes.
		Type() enums.Metric

		// Value returns the current value of the metric, allowing the client to
		// access the performance data or other relevant information that the metric
		// represents. This can be used for monitoring, reporting, or making
		// decisions based on the metric's value during the traversal process.
		Value() MetricValue
	}

	// MutableMetric represents write access to the metric
	MutableMetric interface {
		Metric
		// Tick increments the metric by one and returns the new value. This can be used
		// to track occurrences of specific events or conditions during the traversal process,
		// allowing the client to easily update the metric as needed.
		Tick() MetricValue

		// Times increments the metric by a specified amount and returns the new value. This
		// can be used to track occurrences of specific events or conditions during the
		// traversal process in a more flexible way, allowing the client to update the metric
		// by any desired amount as needed.
		Times(increment uint) MetricValue
	}

	// NavigationMetric represents a specific metric being tracked during the traversal process,
	// including its type and current value. This struct can be used to store and manage
	// individual metrics, allowing for easy access and updates to the metric's value as
	// the traversal progresses.
	NavigationMetric struct {
		// T represents the type of the metric, which can be used to identify and categorize
		// the metric for reporting and analysis purposes.
		T enums.Metric

		// Counter holds the current value of the metric, allowing the client to access and update
		// the performance data or other relevant information that the metric represents during
		// the traversal process. This field is mutable and can be updated using the Tick and
		// Times methods of the MutableMetric interface.
		Counter MetricValue
	}

	// Metrics represents a collection of metrics being tracked during the traversal process,
	// where each metric is identified by its type. This allows for easy access and management
	// of multiple metrics, enabling the client to track various aspects of the traversal's
	// performance and behavior in a structured way.
	Metrics map[enums.Metric]*NavigationMetric

	// Reporter represents query access to the metrics Supervisor
	Reporter interface {
		// Count returns the current value of the specified metric, allowing the client to
		// access the performance data or other relevant information that the metric represents.
		// This can be used for monitoring, reporting, or making decisions based on the metric's
		// value during the traversal process.
		Count(enums.Metric) MetricValue
	}
)

// Type returns the type of the metric, which can be used to identify and categorize
// the metric for reporting and analysis purposes.
func (m *NavigationMetric) Type() enums.Metric {
	return m.T
}

// Value returns the current value of the metric, allowing the client to access the performance
// data or other relevant information that the metric represents. This can be used for monitoring,
// reporting, or making decisions based on the metric's value during the traversal process.
func (m *NavigationMetric) Value() MetricValue {
	return m.Counter
}

// Tick increments the metric by one and returns the new value. This can be used to track occurrences
// of specific events or conditions during the traversal process, allowing the client to easily
// update the metric as needed.
func (m *NavigationMetric) Tick() MetricValue {
	m.Counter++

	return m.Counter
}

// Times increments the metric by a specified amount and returns the new value. This
// can be used to track occurrences of specific events or conditions during the traversal
// process in a more flexible way, allowing the client to update the metric by any desired
// amount as needed.
func (m *NavigationMetric) Times(increment uint) MetricValue {
	m.Counter += increment

	return m.Counter
}

// Merge merges another Metrics collection into the current one, combining the values of
// metrics with the same type. If a metric type exists in both collections, their values
// are added together. If a metric type exists only in the other collection, it is added
// to the current collection. This allows for aggregating metrics from different sources
// or stages of the traversal process, providing a comprehensive view of the performance
// and behavior of the traversal.
func (m Metrics) Merge(other Metrics) {
	for mt := range maps.Keys(m) {
		if om, foundOther := other[mt]; foundOther {
			if metric, found := m[mt]; found {
				metric.Times(om.Counter)
			} else {
				m[mt] = &NavigationMetric{
					T:       mt,
					Counter: om.Counter,
				}
			}
		}
	}
}
