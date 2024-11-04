package resume

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/pref"
)

// 📦 pkg: resume - depends on hiber, filter and persist.
// filter should also contain persistence concerns (actually
// these may be internal modules, eg internal/serial/JSON).

const (
	badge             = "badge: resume"
	followingSiblings = true
)

type (
	Strategy interface {
		init(load *opts.LoadInfo) error
		resume(context.Context, *pref.Was) (*enclave.KernelResult, error)
		ifResult() bool
	}

	baseStrategy struct {
		o        *pref.Options
		active   *core.ActiveState
		was      *pref.Was
		sealer   enclave.GuardianSealer
		kc       enclave.KernelController
		mediator enclave.Mediator
		forest   *core.Forest
	}

	conclusion struct {
		active    *core.ActiveState
		tree      string
		current   string
		inclusive bool
	}

	shard struct {
		siblings *kernel.Contents
	}
)
