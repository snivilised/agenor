package refine

import (
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options) types.Plugin {
	active := o.Core.Filter.Node != nil
	plugin := lo.TernaryF(active,
		func() types.Plugin {
			return &Plugin{}
		},
		func() types.Plugin {
			return nil
		},
	)

	return plugin
}

type Plugin struct {
}

func (p *Plugin) Name() string {
	return "filtering"
}

func (p *Plugin) Init() error {
	return nil
}
