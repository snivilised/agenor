package tv

import (
	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/feat/filter"
	"github.com/snivilised/traverse/internal/feat/hiber"
	"github.com/snivilised/traverse/internal/feat/nanny"
	"github.com/snivilised/traverse/internal/feat/resume"
	"github.com/snivilised/traverse/internal/feat/sampling"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

const (
	noOverwrite = true
)

type (
	ifActive func(o *pref.Options,
		using *pref.Using, mediator types.Mediator,
	) types.Plugin
)

// features interrogates options and invokes requests on behalf of the user
// to activate features according to option selections. other plugins will
// be initialised after primary plugins
func features(o *pref.Options, using *pref.Using, mediator types.Mediator,
	kc types.KernelController,
	others ...types.Plugin,
) (plugins []types.Plugin, err error) {
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
		func(acc []types.Plugin, plugin types.Plugin, _ int) []types.Plugin {
			if plugin != nil {
				acc = append(acc, plugin)
			}
			return acc
		},
		lo.Reduce(all,
			func(acc []types.Plugin, query ifActive, _ int) []types.Plugin {
				if plugin := query(o, using, mediator); plugin != nil {
					acc = append(acc, plugin)
				}
				return acc
			},
			[]types.Plugin{},
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
func Prime(using *pref.Using, settings ...pref.Option) *Builders {
	return &Builders{
		using: using,
		forest: pref.CreateForest(func(root string) *core.Forest {
			if using.GetForest != nil {
				return using.GetForest(root)
			}
			return &core.Forest{
				T: nef.NewTraverseFS(Rel{
					Root:      root,
					Overwrite: noOverwrite,
				}),
				R: nef.NewTraverseABS(),
			}
		}),
		extent: extension(func(forest *core.Forest) extent {
			return &primeExtent{
				baseExtent: baseExtent{
					trees: forest,
				},
				u: using,
			}
		}),
		harvest: optionBuilder(func(ext extent) (types.OptionHarvest, error) {
			type baggage struct {
				harvest types.OptionHarvest
				err     error
			}

			b := func(ve error) *baggage {
				return lo.TernaryF(using.O != nil,
					func() *baggage {
						return &baggage{
							harvest: &optionHarvest{
								o:      using.O,
								binder: opts.Push(using.O),
							},
							err: ve,
						}
					},
					func() *baggage {
						oa, err := ext.options(settings...)

						return &baggage{
							harvest: oa,
							err:     lo.Ternary(ve != nil, ve, err),
						}
					},
				)
			}(using.Validate())

			return b.harvest, b.err
		}),
		navigator: kernel.Builder(func(harvest types.OptionHarvest, // pass in controls here, or put on resources
			resources *types.Resources,
		) (*kernel.Artefacts, error) {
			return kernel.New(using, harvest.Options(), &kernel.Benign{}, resources), nil
		}),
		plugins: activated(features),
	}
}

// Resume extent requests that the navigator performs a resume
// traversal, loading state from a previously saved session
// as a result of it being terminated prematurely via a ctrl-c
// interrupt.
func Resume(was *Was, settings ...pref.Option) *Builders {
	return &Builders{
		using: &was.Using,
		forest: pref.CreateForest(func(root string) *core.Forest {
			if was.Using.GetForest != nil {
				return was.Using.GetForest(root)
			}
			return &core.Forest{
				T: nef.NewTraverseFS(Rel{
					Root:      root,
					Overwrite: noOverwrite,
				}),
				R: nef.NewTraverseABS(),
			}
		}),
		extent: extension(func(forest *core.Forest) extent {
			return &resumeExtent{
				baseExtent: baseExtent{
					trees: forest,
				},
				w: was,
			}
		}),
		// we need state; record the hibernation wake point, so
		// using a func here is probably not optimal.
		//
		harvest: optionBuilder(func(ext extent) (harvest types.OptionHarvest, err error) {
			harvest, err = ext.options(settings...)

			if err != nil {
				return harvest, err
			}

			err = was.Validate()

			return harvest, err
		}),
		navigator: kernel.Builder(func(harvest types.OptionHarvest,
			resources *types.Resources,
		) (*kernel.Artefacts, error) {
			return resume.NewController(
				was,
				harvest,
				kernel.New(&was.Using, harvest.Options(), resume.GetSealer(was), resources),
			), nil
		}),
		plugins: activated(features),
	}
}
