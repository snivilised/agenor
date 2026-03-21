package hiber

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/internal/enclave"
	"github.com/snivilised/jaywalk/src/agenor/life"
	"github.com/snivilised/jaywalk/src/agenor/pref"
)

type (
	nextFn func(servant core.Servant, node *core.Node,
		inspection enclave.Inspection,
	) (bool, error)

	state struct {
		next nextFn
	}
	hibernateStates map[enums.Hibernation]state

	triggers struct {
		wake  core.TraverseFilter
		sleep core.TraverseFilter
	}

	profile interface {
		init(controls *life.Controls) error
		next(servant core.Servant, node *core.Node,
			inspection enclave.Inspection,
		) (bool, error)
	}

	common struct {
		ho       *core.HibernateOptions
		fo       *pref.FilterOptions
		triggers triggers
		controls *life.Controls
	}
)

func launch(ho *core.HibernateOptions) enums.Hibernation {
	if ho.WakeAt != nil {
		return enums.HibernationPending
	}

	return enums.HibernationActive
}
