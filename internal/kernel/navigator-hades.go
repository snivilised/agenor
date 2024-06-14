package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
)

func HadesNav(err error) core.Navigator {
	return &navigatorHades{
		err: err,
	}
}

type hadesResult struct {
	err error
}

func (r *hadesResult) Error() error {
	return r.err
}

type navigatorHades struct {
	err error
}

func (n *navigatorHades) Navigate(_ context.Context) (core.TraverseResult, error) {
	return &hadesResult{
		err: n.err,
	}, n.err
}
