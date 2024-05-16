package pref

import (
	"github.com/snivilised/traverse/core"
)

type FilterOptions struct {
	Node *core.FilterDef
}

func WithFilter(filter *core.FilterDef) Option {
	return func(o *Options) error {
		o.Core.Filter.Node = filter

		return nil
	}
}
