package resume

import (
	"context"

	"github.com/pkg/errors"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/i18n"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type Controller struct {
	controller core.Navigator
	was        *pref.Was
	strategy   resumeStrategy
	facilities types.Facilities
}

func NewController(was *pref.Was, artefacts *kernel.Artefacts) *kernel.Artefacts {
	// The Navigator on the incoming artefacts is the core navigator. It is
	// decorated here for resume. The strategy only needs access to the core navigator.
	// The resume navigator delegates to the strategy.
	//
	var (
		strategy resumeStrategy
		err      error
	)

	if strategy, err = newStrategy(was, artefacts.Navigator); err != nil {
		return artefacts
	}

	return &kernel.Artefacts{
		Navigator: &Controller{
			controller: artefacts.Navigator,
			was:        was,
			strategy:   strategy,
			facilities: artefacts.Facilities,
		},
		Mediator: artefacts.Mediator,
	}
}

func newStrategy(was *pref.Was, nav core.Navigator) (strategy resumeStrategy, err error) {
	driver, ok := nav.(kernel.NavigatorDriver)

	if !ok {
		return nil, i18n.ErrInternalFailedToGetNavigatorDriver
	}

	base := baseStrategy{
		o:    was.O,
		nav:  nav,
		impl: driver.Impl(),
	}

	switch was.Strategy {
	case enums.ResumeStrategyFastward:
		strategy = &fastwardStrategy{
			baseStrategy: base,
		}
	case enums.ResumeStrategySpawn:
		strategy = &spawnStrategy{
			baseStrategy: base,
		}
	case enums.ResumeStrategyUndefined:
	}

	return strategy, nil
}

func (c *Controller) Navigate(_ context.Context) (core.TraverseResult, error) {
	return &types.KernelResult{
		Err: errors.Wrap(core.ErrNotImpl, "resume.Controller.Navigate"),
	}, nil
}
