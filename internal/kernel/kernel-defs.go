package kernel

type navigationResult struct {
	err error
}

func (r *navigationResult) Error() error {
	return r.err
}

type navigatorImpl interface {
}
