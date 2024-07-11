package override

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/measure"
	"github.com/snivilised/traverse/tapable"
)

// override package provides a similar function to tapable except we
// use the name action to replace hook. The difference between the
// two are that hooks allow for the client to customise core internal
// behaviour, where as an action allows for internal behaviour to
// be customised by internal entities. One might wonder why this isn't
// implemented inside types as that package is for internal affairs only,
// but types does provide any functionality and types has dependencies
// that we should avoid in override; that is to say we need to avoid
// circular dependencies;...

type (
	Action[F any] interface {
		tapable.Invokable[F]
		// Intercept overrides the default tap-able core function
		Intercept(handler F)
	}

	Actions struct {
		HandleChildren Action[HandleChildrenInterceptor]
	}

	// ActionCtrl contains the handler function to be invoked. The control
	// is agnostic to the handler's signature and therefore can not invoke it.
	ActionCtrl[F any] struct {
		handler F
		def     F
	}

	HandleChildrenInterceptor func(
		inspection core.Inspection,
		mums measure.MutableMetrics,
	)
)

func NewActionCtrl[F any](handler F) *ActionCtrl[F] {
	return &ActionCtrl[F]{
		handler: handler,
	}
}

// add life-cycle style broadcast

func (c *ActionCtrl[F]) Intercept(handler F) {
	c.handler = handler
}

func (c *ActionCtrl[F]) Invoke() F {
	return c.handler
}
