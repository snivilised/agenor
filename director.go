package age

import (
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/feat/filter"
	"github.com/snivilised/agenor/internal/feat/hiber"
	"github.com/snivilised/agenor/internal/feat/nanny"
	"github.com/snivilised/agenor/internal/feat/resume"
	"github.com/snivilised/agenor/internal/feat/sampling"
	"github.com/snivilised/agenor/internal/kernel"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/pref"
)

const (
	noOverwrite = true
)

type (
	ifActive func(o *pref.Options,
		sub enums.Subscription, mediator enclave.Mediator,
	) enclave.Plugin
)

func features(o *pref.Options,
	ext extent,
	artefacts *kernel.Artefacts,
	others ...enclave.Plugin,
) (plugins []enclave.Plugin, err error) {
	var (
		all = []ifActive{
			// filtering must happen before sampling so that
			// ReadDirectory hooks are applied to incorrect
			// order. How can we decouple ourselves from this
			// requirement? => the cure is worse than the disease
			//
			hiber.IfActive, nanny.IfActive, filter.IfActive, sampling.IfActive,
		}
	)

	// double reduce, the inner reduce 'all' creates list of active plugins
	// and the outer, adds other plugins to the activated list.
	plugins = lo.Reduce(others,
		func(acc []enclave.Plugin, plugin enclave.Plugin, _ int) []enclave.Plugin {
			if plugin != nil {
				acc = append(acc, plugin)
			}
			return acc
		},
		lo.Reduce(all,
			func(acc []enclave.Plugin, query ifActive, _ int) []enclave.Plugin {
				if plugin := query(o, ext.subscription(), artefacts.Mediator); plugin != nil {
					acc = append(acc, plugin)
				}
				return acc
			},
			[]enclave.Plugin{},
		),
	)

	for _, plugin := range plugins {
		err = plugin.Register(artefacts.Kontroller)

		if err != nil {
			return nil, err
		}
	}

	return plugins, nil
}

// Prime extent requests that the navigator performs a full
// traversal from the root path specified.
func Prime(facade pref.Facade, settings ...pref.Option) *Builders {
	return &Builders{
		facade: facade,
		scaffold: scaffolding(func(addons ...Addon) (scaffold, error) {
			return newPrimaryPlatform(facade, addons, settings...)
		}),
		navigator: kernel.Builder(func(inception *kernel.Inception) *kernel.Artefacts {
			return kernel.PrimeArtefacts(
				inception,
				&kernel.Benign{},
			)
		}),
		plugins: activated(features),
	}
}

// Resume extent requests that the navigator performs a resume
// traversal, loading state from a previously saved session
// as a result of it being terminated prematurely via a ctrl-c
// interrupt.
func Resume(facade pref.Facade, settings ...pref.Option) *Builders {
	return &Builders{
		facade: facade,
		scaffold: scaffolding(func(addons ...Addon) (scaffold, error) {
			return newResumePlatform(facade, addons, settings...)
		}),
		navigator: kernel.Builder(resume.Artefacts),
		plugins:   activated(features),
	}
}
