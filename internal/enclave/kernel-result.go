package enclave

import (
	"github.com/snivilised/agenor/core"
)

// KernelResult is the internal representation of core.TraverseResult
type KernelResult struct {
	session  core.Session
	reporter core.Reporter
	complete bool
	err      error
}

func NewResult(session core.Session,
	supervisor *core.Supervisor,
	err error,
	complete bool,
) *KernelResult {
	return &KernelResult{
		session:  session,
		reporter: supervisor,
		err:      err,
		complete: complete,
	}
}

func NewFailed(err error) *KernelResult {
	return &KernelResult{
		err: err,
	}
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

func (r *KernelResult) Error() error {
	return r.err
}
