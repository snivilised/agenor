package sampling

import (
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options) types.Plugin {
	active := (o.Core.Sampling.NoOf.Files > 0) || (o.Core.Sampling.NoOf.Folders > 0)
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
	return "sampling"
}

func (p *Plugin) Init() error {
	return nil
}
