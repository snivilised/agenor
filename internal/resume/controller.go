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
	was      *pref.Was
	nav      core.Navigator
	strategy resumeStrategy
}

func NewController(was *pref.Was, nav core.Navigator) (ctrl *Controller, err error) {
	var (
		strategy resumeStrategy
	)

	if strategy, err = newStrategy(was, nav); err != nil {
		return nil, err
	}

	return &Controller{
		was:      was,
		nav:      nav,
		strategy: strategy,
	}, nil
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
	return &types.NavigateResult{
		Err: errors.Wrap(core.ErrNotImpl, "resume.Controller.Navigate"),
	}, nil
}
