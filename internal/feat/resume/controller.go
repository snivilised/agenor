package resume

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

type Controller struct {
	kc         types.KernelController
	was        *pref.Was
	strategy   resumeStrategy
	facilities types.Facilities
}

func (c *Controller) Ignite(ignition *types.Ignition) {
	c.kc.Ignite(ignition)
}

func (c *Controller) Result(ctx context.Context, err error) *types.KernelResult {
	return c.kc.Result(ctx, err)
}

func (c *Controller) Mediator() types.Mediator {
	return c.kc.Mediator()
}

func (c *Controller) Conclude(result core.TraverseResult) {
	c.kc.Conclude(result)
}

func NewController(was *pref.Was, artefacts *kernel.Artefacts) *kernel.Artefacts {
	// The Controller on the incoming artefacts is the core navigator. It is
	// decorated here for resume. The strategy only needs access to the core navigator.
	// The resume navigator delegates to the strategy.
	//
	var (
		strategy resumeStrategy
		err      error
	)

	if strategy, err = newStrategy(was, artefacts.Kontroller); err != nil {
		return artefacts
	}

	return &kernel.Artefacts{
		Kontroller: &Controller{
			kc:         artefacts.Kontroller,
			was:        was,
			strategy:   strategy,
			facilities: artefacts.Facilities,
		},
		Mediator:  artefacts.Mediator,
		Resources: artefacts.Resources,
		IfResult:  strategy.ifResult,
	}
}

func newStrategy(was *pref.Was, kc types.KernelController) (strategy resumeStrategy, err error) {
	driver, ok := kc.(kernel.NavigatorDriver)

	if !ok {
		return nil, locale.ErrInternalFailedToGetNavigatorDriver
	}

	base := baseStrategy{
		o:    was.O,
		kc:   kc,
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

func (c *Controller) Navigate(ctx context.Context) (core.TraverseResult, error) {
	return c.strategy.resume(ctx)
}
