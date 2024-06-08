package tv

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/hiber"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/refine"
	"github.com/snivilised/traverse/internal/sampling"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type ifActive func(o *pref.Options) types.Plugin

// activated interrogates options and invokes requests on behalf of the user
// to activate features according to option selections
func activated(o *pref.Options) ([]types.Plugin, error) {
	var (
		all = []ifActive{
			hiber.IfActive, refine.IfActive, sampling.IfActive,
		}
		plugins = []types.Plugin{}
		err     error
	)

	for _, active := range all {
		if plugin := active(o); plugin != nil {
			plugins = append(plugins, plugin)
			err = plugin.Init()
		}
	}

	return plugins, err
}

// Prime extent requests that the navigator performs a full
// traversal from the root path specified.
func Prime(using pref.Using, settings ...pref.Option) *Builders {
	return &Builders{
		extent: &primeExtent{
			u: &using,
		},
		options: optionals(func() (*pref.Options, error) {
			if err := using.Validate(); err != nil {
				return nil, err
			}

			if using.O != nil {
				return using.O, nil
			}

			// we probably need to mark something somehow to indicate
			// Prime
			//
			return pref.Get(settings...)
		}),
		navigator: builder(func(o *pref.Options) (core.Navigator, error) {
			return kernel.PrimeNav(using, o)
		}),
		plugins: features(activated),
	}
}

// Resume extent requests that the navigator performs a resume
// traversal, loading state from a previously saved session
// as a result of it being terminated prematurely via a ctrl-c
// interrupt.
func Resume(was Was, settings ...pref.Option) *Builders {
	return &Builders{
		extent: &resumeExtent{
			w: &was,
		},
		// we need state; record the hibernation wake point, so
		// using a func here is probably not optimal.
		//
		options: optionals(func() (*pref.Options, error) {
			if err := was.Validate(); err != nil {
				return nil, err
			}

			// TODO: we probably need to mark something somehow to indicate
			// Resume so we can query the hibernation condition and
			// apply; this has been done by extent, so that querying
			// for hibernation condition just needs to be added to
			// the extent interface.
			//
			o, err := pref.Load(was.From, settings...)

			// get the resume point from the resume persistence file
			// then set up hibernation with this defined as a hibernation
			// filter.
			//
			return o, err
		}),
		navigator: builder(func(o *pref.Options) (core.Navigator, error) {
			// at this point, the resume controller does not know
			// the wake point as would be loaded by the options
			// builder.
			//
			return kernel.ResumeNav(was, o,
				kernel.DecorateController(func(n core.Navigator) core.Navigator {
					// TODO: create the resume controller
					//
					return n
				}),
			)
		}),
		plugins: features(activated),
	}
}
