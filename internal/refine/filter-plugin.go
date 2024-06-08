package refine

import (
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options) types.Plugin {
	active := o.Core.Filter.Node != nil

	if active {
		return &Plugin{}
	}

	return nil
}

type Plugin struct {
}

func (p *Plugin) Name() string {
	return "filtering"
}

func (p *Plugin) Init() error {
	return nil
}
