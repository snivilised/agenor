package tv

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/hiber"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/refine"
	"github.com/snivilised/traverse/sampling"
)

type duffPrimeController struct {
	err     error
	root    string
	client  core.Client
	from    string
	options []pref.Option
}

type duffResult struct{}

func (r *duffResult) Error() error {
	return nil
}

func (c *duffPrimeController) Navigate() (core.TraverseResult, error) {
	return &duffResult{}, nil
}

type duffResumeController struct {
	err      error
	from     string
	strategy enums.ResumeStrategy
}

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

func (c *duffResumeController) Navigate() (core.TraverseResult, error) {
	return &duffResult{}, nil
}

// Prime extent
func Prime(using core.Using, settings ...pref.Option) *Builders {
	return &Builders{
		ob: optionals(func() (*pref.Options, error) {
			if err := using.Validate(); err != nil {
				return nil, err
			}

			// we probably need to mark something somehow to indicate
			// Prime
			//
			return pref.Get(settings...)
		}),
		nb: factory(func(o *pref.Options) (core.Navigator, error) {
			controller, err := kernel.Prime(using, o)

			if err != nil {
				controller = &duffPrimeController{
					err: err,
				}
			}

			return controller, err
		}),
		pb: features(activated),
	}
}

// Resume extent
func Resume(as As, settings ...pref.Option) *Builders {
	return &Builders{
		// we need state; record the hibernation wake point, so
		// using a func here is probably not optimal.
		//
		ob: optionals(func() (*pref.Options, error) {
			if err := as.Validate(); err != nil {
				return nil, err
			}

			// we probably need to mark something somehow to indicate
			// Resume so we can query the hibernation condition and
			// apply.
			//
			o, err := pref.Load(as.From, settings...)

			// get the resume point from the resume persistence file
			// then set up hibernation with this defined as a hibernation
			// filter.
			//
			return o, err
		}),
		nb: factory(func(o *pref.Options) (core.Navigator, error) {
			controller, err := kernel.Resume(as, o,
				kernel.DecorateController(func(n core.Navigator) core.Navigator {
					// TODO: create the resume controller
					//
					return n
				}),
			)

			if err != nil {
				controller = &duffResumeController{
					err: err,
				}
			}

			// at this point, the resume controller does not know
			// the wake point as would be loaded by the options
			// builder.
			//
			return controller, err
		}),
		pb: features(activated),
	}
}

// Director
type Director interface {
	// Extent represents the magnitude of the traversal; ie we can
	// perform a full Prime run, or Resume from a previously
	// cancelled run.
	//
	Extent(bs *Builders) core.Navigator
}

// NavigatorFactory
type NavigatorFactory interface {
	// Configure is a factory function that creates a navigator.
	// We don't return an error here as that would make using the factory
	// awkward. Instead, if there is an error during the build process,
	// we return a fake navigator that when invoked immediately returns
	// a traverse error indicating the build issue.
	//
	Configure() Director
}

type walker struct { // NavigatorFactory
}

func (f *walker) Configure() Director {
	// Walk
	//
	return director(func(bs *Builders) core.Navigator {
		// resume or prime? If resume, we need to access the hibernation
		// wake condition on the retrieved options. But how do we know what
		// the extent is, so we know if we need to make this query?
		//
		//
		artefacts, _ := bs.buildAll() // TODO: check error

		// Announce the availability of the navigator via UsePlugin interface
		ctx, _ := artefacts.o.Acceleration.Cancellation()
		_ = services.Broker.Emit(ctx, services.TopicInterceptNavigator, artefacts.nav)

		return &driver{
			session{
				// do we store a context/cancel on the session? (and pass in via Configure)
				//
				o:       artefacts.o,
				nav:     artefacts.nav,
				plugins: artefacts.plugins,
			},
		}
	})
}

type runner struct { // NavigatorFactory
}

func (f *runner) Configure() Director {
	// Run: create the observable/worker-pool
	//
	return director(func(bs *Builders) core.Navigator {
		artefacts, _ := bs.buildAll() // TODO: check error

		return &driver{
			session{
				o:   artefacts.o,
				nav: artefacts.nav,
			},
		}
	})
}

func Walk() NavigatorFactory {
	return &walker{}
}

func Run() NavigatorFactory {
	return &runner{}
}
