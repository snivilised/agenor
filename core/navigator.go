package core

import (
	"context"
)

type Navigator interface {
	Navigate(ctx context.Context) (TraverseResult, error)
}
