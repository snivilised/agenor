package core

type Navigator interface {
	Navigate() (TraverseResult, error)
}

type NavigatorFunc func() (TraverseResult, error)

func (fn NavigatorFunc) Navigate() (TraverseResult, error) {
	return fn()
}
