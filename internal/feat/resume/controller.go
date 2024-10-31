package resume

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

type Controller struct {
	kc         types.KernelController
	was        *pref.Was
	load       *opts.LoadInfo
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

func (c *Controller) Resume(context.Context,
	*core.ActiveState,
) (*types.KernelResult, error) {
	return &types.KernelResult{}, nil
}

func (c *Controller) Conclude(result core.TraverseResult) {
	c.kc.Conclude(result)
}

func newStrategy(was *pref.Was,
	harvest types.OptionHarvest,
	kc types.KernelController,
) (strategy resumeStrategy, err error) {
	driver, ok := kc.(kernel.NavigatorDriver)
	active := harvest.Loaded().State

	if !ok {
		return nil, locale.ErrInternalFailedToGetNavigatorDriver
	}

	switch was.Strategy {
	case enums.ResumeStrategyFastward:
		strategy = &fastwardStrategy{
			baseStrategy: baseStrategy{
				active: active,
				was:    was,
				kc:     kc,
				impl:   driver.Impl(),
			},
			role: enums.RoleFastward,
		}

	case enums.ResumeStrategySpawn:
		strategy = &spawnStrategy{
			baseStrategy: baseStrategy{
				active: active,
				was:    was,
				kc:     kc,
				impl:   driver.Impl(),
			},
		}
	case enums.ResumeStrategyUndefined:
	}

	// we can't do this:
	// kc.Mediator().Decorate(strategy)
	//
	// ... here, because its too early; the Mediator
	// is not fully constructed
	//
	return strategy, nil
}

func (c *Controller) Navigate(ctx context.Context) (*types.KernelResult, error) {
	if err := c.kc.Mediator().Decorate(c.strategy); err != nil {
		return c.Result(ctx, err), err
	}

	if err := c.strategy.init(c.load); err != nil {
		return c.Result(ctx, err), err
	}

	return c.strategy.resume(ctx, c.was)
}
