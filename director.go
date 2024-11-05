package tv

import (
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/feat/filter"
	"github.com/snivilised/traverse/internal/feat/hiber"
	"github.com/snivilised/traverse/internal/feat/nanny"
	"github.com/snivilised/traverse/internal/feat/resume"
	"github.com/snivilised/traverse/internal/feat/sampling"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
)

const (
	noOverwrite = true
)

type (
	ifActive func(o *pref.Options,
		facade pref.Facade, mediator enclave.Mediator,
	) enclave.Plugin
)

// features interrogates options and invokes requests on behalf of the user
// to activate features according to option selections. other plugins will
// be initialised after primary plugins
func features(o *pref.Options, facade pref.Facade, mediator enclave.Mediator,
	kc enclave.KernelController,
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

	// double reduce, the first reduce 'all' creates list of active plugins
	// and the second, adds other plugins to the activated list.
	plugins = lo.Reduce(others,
		func(acc []enclave.Plugin, plugin enclave.Plugin, _ int) []enclave.Plugin {
			if plugin != nil {
				acc = append(acc, plugin)
			}
			return acc
		},
		lo.Reduce(all,
			func(acc []enclave.Plugin, query ifActive, _ int) []enclave.Plugin {
				if plugin := query(o, facade, mediator); plugin != nil {
					acc = append(acc, plugin)
				}
				return acc
			},
			[]enclave.Plugin{},
		),
	)

	for _, plugin := range plugins {
		err = plugin.Register(kc)

		if err != nil {
			return nil, err
		}
	}

	return plugins, nil
}

// Prime extent requests that the navigator performs a full
// traversal from the root path specified.
func Prime(facade pref.Facade, settings ...pref.Option) *Builders {
	using, _ := facade.(*pref.Using) // TODO: Create a test that accidentally sets relic facade

	return &Builders{
		facade: facade,
		scaffold: scaffolding(func(facade pref.Facade) (scaffold, error) {
			return newPrimaryPlatform(facade, settings...)
		}),
		navigator: kernel.Builder(func(harvest enclave.OptionHarvest,
			resources *enclave.Resources,
		) *kernel.Artefacts {
			return kernel.WithArtefacts(
				using,
				harvest.Options(),
				resources,
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
	relic, _ := facade.(*pref.Relic) // TODO: Create a test that accidentally sets using facade

	return &Builders{
		facade: facade,
		scaffold: scaffolding(func(pref.Facade) (scaffold, error) {
			return newResumePlatform(facade, settings...)
		}),
		navigator: kernel.Builder(func(harvest enclave.OptionHarvest,
			resources *enclave.Resources,
		) *kernel.Artefacts {
			return resume.Artefacts(
				relic,
				harvest,
				resources,
			)
		}),
		plugins: activated(features),
	}
}
