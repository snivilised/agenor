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

func NewFailed() *KernelResult {
	return &KernelResult{}
}

func (r *KernelResult) Merge(other core.Metrics) {
	r.reporter.metrics.Merge(other)
}

func (r *KernelResult) IsComplete() bool {
	return r.complete
}

func (r *KernelResult) Session() core.Session {
	return r.session
}

func (r *KernelResult) Metrics() core.Reporter {
	return r.reporter
}
