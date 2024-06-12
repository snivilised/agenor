package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type NavigationController struct {
	impl NavigatorImpl
	o    *pref.Options
}

func (nc *NavigationController) Navigate(context.Context) (core.TraverseResult, error) {
	return &navigationResult{}, nil
}

func (nc *NavigationController) Impl() NavigatorImpl {
	return nc.impl
}

func (nc *NavigationController) Register(plugin types.Plugin) error {
	_ = plugin

	return nil
}

func (nc *NavigationController) Interceptor() types.Interception {
	return nil
}

func (nc *NavigationController) Facilitate() types.Facilities {
	return nil
}
