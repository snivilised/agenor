package resume

import (
	"context"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/kernel"
	"github.com/snivilised/agenor/internal/opts"
	"github.com/snivilised/agenor/pref"
)

// Controller controls the traversal.
type Controller struct {
	med      enclave.Mediator
	relic    *pref.Relic
	load     *opts.LoadInfo
	strategy Strategy
}

// Ignite ignites the controller.
func (c *Controller) Ignite(ignition *enclave.Ignition) {
	c.strategy.ignite()
	c.med.Ignite(ignition)
}

// Result returns the result of the traversal.
func (c *Controller) Result(ctx context.Context) *enclave.KernelResult {
	return c.med.Result(ctx)
}

// Strategy returns the strategy of the controller.
func (c *Controller) Strategy() Strategy {
	return c.strategy
}

// Snooze snoozes the controller.
func (c *Controller) Snooze(ctx context.Context,
	_ *core.ActiveState,
) (*enclave.KernelResult, error) {
	return c.Result(ctx), nil
}

// Bye is called when the traversal is finished.
func (c *Controller) Bye(result core.TraverseResult) {
	c.med.Bye(result)
}

func newStrategy(inception *kernel.Inception,
	sealer enclave.GuardianSealer,
	mediator enclave.Mediator,
) (strategy Strategy) {
	load := inception.Harvest.Loaded()
	relic, _ := inception.Facade.(*pref.Relic)
	base := baseStrategy{
		o:         load.O,
		active:    load.State,
		relic:     relic,
		sealer:    sealer,
		kc:        mediator,
		mediator:  mediator,
		forest:    inception.Resources.Forest,
		resources: inception.Resources,
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

// Navigate navigates the tree within the context of resume.
func (c *Controller) Navigate(ctx context.Context) (*enclave.KernelResult, error) {
	if err := c.strategy.init(c.load); err != nil {
		return c.Result(ctx), err
	}

	return c.strategy.resume(ctx)
}
