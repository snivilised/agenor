package cycle_test

const (
	traversalRoot = "/traversal-root"
	anotherRoot   = "/another-root"
)

type testResult struct {
	err error
}

func (r *testResult) Error() error {
	return r.err
}
