package core

import (
	"context"
)

type Navigator interface {
	Navigate(ctx context.Context) (TraverseResult, error)
}

type Navigate func() (TraverseResult, error)

func (fn Navigate) Navigate() (TraverseResult, error) {
	return fn()
}
