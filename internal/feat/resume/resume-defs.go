package resume

import (
	"context"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/kernel"
	"github.com/snivilised/agenor/internal/opts"
	"github.com/snivilised/agenor/pref"
)

const (
	followingSiblings = true
)

type (
	// Strategy is the strategy of the resume controller.
	Strategy interface {
		init(load *opts.LoadInfo) error
		ignite()
		resume(context.Context) (*enclave.KernelResult, error)
		ifResult() bool
	}

	baseStrategy struct {
		o         *pref.Options
		active    *core.ActiveState
		relic     *pref.Relic
		sealer    enclave.GuardianSealer
		kc        enclave.KernelController
		mediator  enclave.Mediator
		forest    *core.Forest
		resources *enclave.Resources
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
