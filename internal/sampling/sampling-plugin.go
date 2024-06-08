package sampling

import (
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options) types.Plugin {
	active := (o.Core.Sampling.NoOf.Files > 0) || (o.Core.Sampling.NoOf.Folders > 0)

	if active {
		return &Plugin{}
	}

	return nil
}

type Plugin struct {
}

func (p *Plugin) Name() string {
	return "sampling"
}

func (p *Plugin) Init() error {
	return nil
}
