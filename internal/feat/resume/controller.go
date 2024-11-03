package resume

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type Controller struct {
	kc         types.KernelController
	was        *pref.Was
	load       *opts.LoadInfo
	strategy   Strategy
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

func (c *Controller) Strategy() Strategy {
	return c.strategy
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
	sealer types.GuardianSealer,
	resources *types.Resources,
) (strategy Strategy) {
	load := harvest.Loaded()
	base := baseStrategy{
		o:        load.O,
		active:   load.State,
		was:      was,
		sealer:   sealer,
		kc:       kc,
		mediator: kc.Mediator(),
		forest:   resources.Forest,
	}

	switch was.Strategy {
	case enums.ResumeStrategyFastward:
		strategy = &fastwardStrategy{
			baseStrategy: base,
			role:         enums.RoleFastward,
		}

	case enums.ResumeStrategySpawn:
		strategy = &spawnStrategy{
			baseStrategy: base,
		}
	case enums.ResumeStrategyUndefined:
	}

	return strategy
}

func (c *Controller) Navigate(ctx context.Context) (*types.KernelResult, error) {
	if err := c.strategy.init(c.load); err != nil {
		return c.Result(ctx, err), err
	}

	return c.strategy.resume(ctx, c.was)
}
