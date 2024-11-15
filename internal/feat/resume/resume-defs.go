package resume

import (
	"context"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/kernel"
	"github.com/snivilised/agenor/internal/opts"
	"github.com/snivilised/agenor/pref"
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
		resume(context.Context) (*enclave.KernelResult, error)
		ifResult() bool
	}

	baseStrategy struct {
		o        *pref.Options
		active   *core.ActiveState
		relic    *pref.Relic
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
