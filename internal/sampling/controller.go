package sampling

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/enums"
)

type handlers struct {
	descend cycle.NodeHandler
	ascend  cycle.NodeHandler
}

type controller struct {
	o  *samplingOptions
	on handlers
}

func (p *controller) Role() enums.Role {
	return enums.RoleSampler
}

func (p *controller) Next(_ *core.Node) (bool, error) {
	return false, nil
}
