package kernel

import (
	"github.com/snivilised/traverse/core"
)

func HadesNav(err error) (core.Navigator, error) {
	return &navigatorHades{
		err: err,
	}, err
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

func (n *navigatorHades) Navigate() (core.TraverseResult, error) {
	return &hadesResult{
		err: n.err,
	}, n.err
}