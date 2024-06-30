package cycle_test

import "github.com/snivilised/traverse/measure"

const (
	traversalRoot = "/traversal-root"
	anotherRoot   = "/another-root"
)

type testResult struct {
	err error
}

func (r *testResult) Metrics() measure.Reporter {
	return nil
}

func (r *testResult) Error() error {
	return r.err
}
