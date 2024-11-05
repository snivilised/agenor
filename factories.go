package tv

import (
	"github.com/snivilised/pants"
	"github.com/snivilised/traverse/core"
)

// Walk requests a sequential traversal of a directory tree.
func Walk() NavigatorFactory {
	return &walkerFac{}
}

// Run requests a concurrent traversal of a directory tree.
func Run(wg pants.WaitGroup) NavigatorFactory {
	return &runnerFac{
		wg: wg,
	}
}

type factory struct {
}

type walkerFac struct {
	factory
}

func (f *walkerFac) Configure() Director {
	return director(func(bs *Builders) core.Navigator {
		// resume or prime? If resume, we need to access the hibernation
		// wake condition on the retrieved options. But how do we know what
		// the extent is, so we know if we need to make this query?
		//
		//
		artefacts, err := bs.buildAll()

		return &driver{
			session{
				sync: &sequential{
					trunk: trunk{ // TODO: why dont we just store the artefacts here?
						kc:  artefacts.kc,
						o:   artefacts.o,
						ext: artefacts.ext,
						err: err,
					},
				},
				plugins: artefacts.plugins,
			},
		}
	})
}

type runnerFac struct {
	factory
	wg pants.WaitGroup
}

func (f *runnerFac) Configure() Director {
	// Run: create the observable/worker-pool
	//
	return director(func(bs *Builders) core.Navigator {
		artefacts, err := bs.buildAll()

		return &driver{
			session{
				sync: &concurrent{
					trunk: trunk{ // TODO: why dont we just store the artefacts here?
						kc:  artefacts.kc,
						o:   artefacts.o,
						ext: artefacts.ext,
						err: err,
					},
					wg: f.wg,
				},
				plugins: artefacts.plugins,
			},
		}
	})
}
