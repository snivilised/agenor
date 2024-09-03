package filtering

import (
	"github.com/snivilised/traverse/internal/third/lo"
)

// ChildFilter ================================================================

// Child filter used when subscription is FoldersWithFiles
type Child struct {
	Name    string
	Pattern string
	Negate  bool
}

func (f *Child) Description() string {
	return f.Name
}

func (f *Child) Validate() error {
	return nil
}

func (f *Child) Source() string {
	return f.Pattern
}

func (f *Child) invert(result bool) bool {
	return lo.Ternary(f.Negate, !result, result)
}
