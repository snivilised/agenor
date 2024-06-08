package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type navigationController struct {
	impl navigatorImpl
	o    *pref.Options
}

func (nc *navigationController) Navigate(context.Context) (core.TraverseResult, error) {
	return &navigationResult{}, nil
}

func (nc *navigationController) Register(plugin types.Plugin) error {
	_ = plugin

	return nil
}

func (nc *navigationController) Interceptor() types.Interception {
	return nil
}

func (nc *navigationController) Facilitate() types.Facilities {
	return nil
}
