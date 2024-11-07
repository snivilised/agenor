package resume

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/pref"
)

type Controller struct {
	kc         enclave.KernelController
	relic      *pref.Relic
	load       *opts.LoadInfo
	strategy   Strategy
	facilities enclave.Facilities
}

func (c *Controller) Ignite(ignition *enclave.Ignition) {
	c.kc.Ignite(ignition)
}

func (c *Controller) Result(ctx context.Context, err error) *enclave.KernelResult {
	return c.kc.Result(ctx, err)
}

func (c *Controller) Mediator() enclave.Mediator {
	return c.kc.Mediator()
}

func (c *Controller) Strategy() Strategy {
	return c.strategy
}

func (c *Controller) Resume(context.Context,
	*core.ActiveState,
) (*enclave.KernelResult, error) {
	return &enclave.KernelResult{}, nil
}

func (c *Controller) Conclude(result core.TraverseResult) {
	c.kc.Conclude(result)
}

func newStrategy(ci *enclave.ControllerInfo,
	kc enclave.KernelController,
) (strategy Strategy) {
	load := ci.Harvest.Loaded()
	relic, _ := ci.Facade.(*pref.Relic)
	base := baseStrategy{
		o:        load.O,
		active:   load.State,
		relic:    relic,
		sealer:   ci.Sealer,
		kc:       kc,
		mediator: kc.Mediator(),
		forest:   ci.Resources.Forest,
	}

	switch relic.Strategy {
	case enums.ResumeStrategyFastward, enums.ResumeStrategyUndefined:
		strategy = &fastwardStrategy{
			baseStrategy: base,
			role:         enums.RoleFastward,
		}

	case enums.ResumeStrategySpawn:
		strategy = &spawnStrategy{
			baseStrategy: base,
		}
	}

	return strategy
}

func (c *Controller) Navigate(ctx context.Context) (*enclave.KernelResult, error) {
	if err := c.strategy.init(c.load); err != nil {
		return c.Result(ctx, err), err
	}

	return c.strategy.resume(ctx)
}
