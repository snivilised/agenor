package kernel

import (
	"github.com/snivilised/traverse/core"
)

type Navigator interface {
	Navigate() (core.TraverseResult, error)
}
