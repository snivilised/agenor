package resume

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

// 📦 pkg: resume - depends on hiber, filter and persist.
// filter should also contain persistence concerns (actually
// these may be internal modules, eg internal/serial/JSON).

const (
	badge = "badge: resume"
)

type resumeStrategy interface {
	types.Link
	init(load *opts.LoadInfo) error
	attach()
	detach()
	resume(context.Context, *pref.Was) (*types.KernelResult, error)
	ifResult() bool
	finish() error
}

type baseStrategy struct {
	active *core.ActiveState
	was    *pref.Was
	kc     types.KernelController
	impl   kernel.NavigatorImpl
}
