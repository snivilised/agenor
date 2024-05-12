package traverse

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/pref"
)

type duffNavigatorController struct {
	root    string
	client  core.Client
	from    string
	options []pref.Option
}

type duffResult struct{}

func (r *duffResult) Error() error {
	return nil
}

func (n *duffNavigatorController) Navigate() (core.TraverseResult, error) {
	return &duffResult{}, nil
}

// Prime extent
func Prime(opts ...pref.Option) OptionsBuilder {
	return optionals(func() *pref.Options {
		binder := pref.NewBinder()

		// we probably need to mark something somehow to indicate
		// Prime
		//
		return pref.Request(binder, opts...)
	})
}

// Resume extent
func Resume(from string, _ enums.ResumeStrategy, opts ...pref.Option) OptionsBuilder {
	// we need state; record the hibernation wake point, so
	// using a func here is probably not optimal.
	//
	return optionals(func() *pref.Options {
		binder := pref.NewBinder()

		// we probably need to mark something somehow to indicate
		// Resume so we can query the hibernation condition and
		// apply.
		//
		load, _ := pref.Load(binder, from, opts...)

		// get the resume point from the resume persistence file
		// then set up hibernation with this defined as a hibernation
		// filter.
		//
		return load.O
	})
}

// Director
type Director interface {
	// Extent represents the magnitude of the traversal; ie we can
	// perform a full Prime run, or Resume from a previously
	// cancelled run.
	//
	Extent(ob OptionsBuilder) core.Navigator
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

type baseFactory struct {
	factory syncBuilder // the sync factory function
}

type linear struct { // NavigatorFactory
	baseFactory
}

func (f *linear) Configure() Director {
	// Walk
	//
	return direct(func(ob OptionsBuilder) core.Navigator {
		// resume or prime? If resume, we need to access the hibernation
		// wake condition on the retrieved options. But how do we know what
		// the extent is, so we know if we need to make this query?
		//
		//
		return &driver{
			session{
				o:   ob.get(),
				nav: &duffNavigatorController{},
			},
		}
	})
}

type accelerator struct { // NavigatorFactory
	baseFactory
}

func (f *accelerator) Configure() Director {
	// Run: create the observable/worker-pool
	//
	return direct(func(ob OptionsBuilder) core.Navigator {
		return &driver{
			session{
				o:   ob.get(),
				nav: &duffNavigatorController{},
			},
		}
	})
}

func Walk() NavigatorFactory {
	// this could just be a function, because linear doesn't carry any
	// state and implements a single method interface
	//
	return &linear{
		baseFactory{
			// TODO: where to invoke this from??? (pass in extent?)
			factory: sync(func(at string) error {
				_ = at

				// TODO: set up hibernation filter on navigator/options
				//
				return nil
			}),
		},
		// extent builder (primary or resume)
		// sync builder (sequential) ---> depends on extent (resume: query hibernate condition)
	}
}

func Run() NavigatorFactory {
	return &accelerator{
		baseFactory{
			factory: sync(func(at string) error {
				_ = at

				// TODO: set up hibernation filter on observable
				//
				return nil
			}),
		},
		// extent builder (primary or resume)
		// sync builder (reactive) ---> depends on extent (resume: query hibernate condition)
	}
}
