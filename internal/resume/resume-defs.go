package resume

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

// refine should also contain persistence concerns (actually
// these may be internal modules, eg internal/serial/JSON). Depends on hiber, refine
// and persist.

type resumeStrategy interface {
	init()
	attach()
	detach()
	resume() (*types.NavigateResult, error)
	finish() error
}

type baseStrategy struct {
	o    *pref.Options
	nav  core.Navigator
	impl kernel.NavigatorImpl
}
