package resume

import (
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type Plugin struct {
	kernel.BasePlugin
}

func (p *Plugin) Name() string {
	return "resume"
}

func (p *Plugin) Register() error {
	return nil
}

func (p *Plugin) Next(node *core.Node) (bool, error) {
	_ = node
	// apply the wake filter

	return true, nil
}

func (p *Plugin) Role() enums.Role {
	return enums.RoleFastward
}

func (p *Plugin) Init() error {
	return p.Mediator.Decorate(p)
}

func GetSealer(was *pref.Was) types.GuardianSealer {
	if was.Strategy == enums.ResumeStrategyFastward {
		return &fastwardGuardianSealer{}
	}

	return &kernel.Benign{}
}

func Load(res fs.FS, from string, settings ...pref.Option) (*pref.LoadInfo, error) {
	return pref.Load(res, from, settings...)
}
