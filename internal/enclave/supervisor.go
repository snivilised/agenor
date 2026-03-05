package enclave

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
)

type (
	// Supervisor is responsible for collecting metrics during a traversal. It
	// provides	 methods to load metrics, retrieve or create individual metrics,
	// and count the value of a specific metric. The Supervisor is used by the
	// kernel to collect and report metrics during a traversal, which allows
	// the caller to analyze the performance and behavior of the traversal
	// after it is completed.
	Supervisor struct {
		metrics core.Metrics
	}

	// Crate is a simple struct that contains a core.Metrics field. It is used to
	// hold the metrics collected during a traversal, which allows the kernel to
	// return the collected metrics to the caller when the traversal is completed.
	// This is useful for providing insights into the performance and behavior of
	// the traversal, and allows the caller to analyze the results after the
	// traversal is finished.
	Crate struct {
		// Metrics is the field that holds the collected metrics during a traversal.
		Metrics core.Metrics
	}
)

// NewSupervisor creates a new Supervisor with an empty set of metrics. This is
// used to initialize a new Supervisor before starting a traversal, which
// allows the kernel to collect metrics during the traversal and report
// them to the caller when the traversal is completed. The Supervisor is
// responsible for managing
// the metrics collected during the traversal, and provides methods to load,
// retrieve, and count metrics as needed.
func NewSupervisor() *Supervisor {
	return &Supervisor{
		metrics: make(core.Metrics),
	}
}

// Load merges the given metrics into the Supervisor's existing metrics.
// This is used to load metrics collected during a traversal into the
// Supervisor, which allows the kernel to keep track of the metrics
// collected during the traversal and report them to the caller when the
// traversal is completed. The Load method allows the kernel to update the
// Supervisor's metrics with new data as it is collected during the traversal,
// which ensures that the Supervisor has an accurate and up-to-date
// set of metrics to report at the end of the traversal.
func (s *Supervisor) Load(metrics core.Metrics) {
	s.metrics.Merge(metrics)
}

// Single returns the metric for the given type, creating it if it
// doesn't exist. This is used to retrieve or create a specific metric
// during a traversal, which allows the kernel to track specific metrics
// as needed during the traversal.
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

// Many returns a set of metrics for the given types, creating them if they
// don't exist. This is used to retrieve or create multiple specific metrics
// during a traversal, which allows the kernel to track multiple metrics as
// needed during the traversal. The Many method allows the kernel to efficiently
// retrieve or create multiple metrics at once, which can be useful for tracking
// a set of related metrics during the traversal.
func (s *Supervisor) Many(metrics ...enums.Metric) core.Metrics {
	result := make(core.Metrics)

	for _, mt := range metrics {
		result[mt] = s.Single(mt)
	}

	return result
}

// Count returns the value of the metric for the given type, or 0 if it doesn't exist.
func (s *Supervisor) Count(mt enums.Metric) core.MetricValue {
	if metric, exists := s.metrics[mt]; exists {
		return metric.Value()
	}

	return 0
}
