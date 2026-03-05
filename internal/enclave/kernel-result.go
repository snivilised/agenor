package enclave

import (
	"github.com/snivilised/agenor/core"
)

// KernelResult is the internal representation of core.TraverseResult
type KernelResult struct {
	session  core.Session
	reporter *Supervisor
	complete bool
}

// NewResult creates a new KernelResult with the given session, supervisor,
// and completion status. This is used to create a new KernelResult when a
// traversal is completed, which allows the kernel to return the session,
// the reporter with the collected metrics, and the completion status of the
// traversal to the caller.
func NewResult(session core.Session,
	supervisor *Supervisor,
	complete bool,
) *KernelResult {
	return &KernelResult{
		session:  session,
		reporter: supervisor,
		complete: complete,
	}
}

// NewFailed creates a new KernelResult with default values, which represents
// a failed traversal. This is used to create a KernelResult when a traversal
// fails, which allows the kernel to return a result that indicates the failure
// without providing any session or metrics information.
func NewFailed() *KernelResult {
	return &KernelResult{}
}

// Merge merges the given metrics into the KernelResult's reporter. This is used to
// merge the metrics collected during a traversal into the KernelResult, which allows
// the kernel to provide a complete set of metrics in the result when the traversal
// is completed.
func (r *KernelResult) Merge(other core.Metrics) {
	r.reporter.metrics.Merge(other)
}

// IsComplete returns true if the traversal is complete, false otherwise. This is
// used to check the completion status of the traversal in the KernelResult, which
// allows the caller to determine if the traversal was completed successfully or if
// it is still in progress or has failed.
func (r *KernelResult) IsComplete() bool {
	return r.complete
}

// Session returns the session associated with the KernelResult.
func (r *KernelResult) Session() core.Session {
	return r.session
}

// Metrics returns the metrics reporter associated with the KernelResult.
func (r *KernelResult) Metrics() core.Reporter {
	return r.reporter
}
