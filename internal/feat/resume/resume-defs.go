package resume

import (
	"context"

	"github.com/snivilised/traverse/internal/kernel"
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
	init()
	attach()
	detach()
	resume(context.Context) (*types.KernelResult, error)
	ifResult() bool
	finish() error
}

type baseStrategy struct {
	o    *pref.Options
	kc   types.KernelController
	impl kernel.NavigatorImpl
}
