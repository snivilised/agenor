package tv

import (
	"io/fs"

	"github.com/snivilised/traverse/internal/feat/hiber"
	"github.com/snivilised/traverse/internal/feat/refine"
	"github.com/snivilised/traverse/internal/feat/resume"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/sampling"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type ifActive func(o *pref.Options, mediator types.Mediator) types.Plugin

// features interrogates options and invokes requests on behalf of the user
// to activate features according to option selections. other plugins will
// be initialised after primary plugins
func features(o *pref.Options, mediator types.Mediator,
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
			hiber.IfActive, refine.IfActive, sampling.IfActive,
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
				if plugin := query(o, mediator); plugin != nil {
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
	// TODO: we need to create an aux file system, which is bound
	// to a pre-defined location, that will be called upon if
	// the navigation session is terminated either by a ctrl-c or
	// by a panic.
	//
	return &Builders{
		readerFS: pref.CreateReadDirFS(func() fs.ReadDirFS {
			if using.GetReadDirFS != nil {
				return using.GetReadDirFS()
			}

			return NewNativeFS(using.Root)
		}),
		queryFS: pref.CreateQueryStatusFS(func(qsys fs.FS) fs.StatFS {
			if using.GetQueryStatusFS != nil {
				return using.GetQueryStatusFS(qsys)
			}

			return NewQueryStatusFS(qsys)
		}),
		extent: extension(func(rsys fs.ReadDirFS, qsys fs.StatFS) extent {
			return &primeExtent{
				baseExtent: baseExtent{
					fileSys: fileSystems{
						nas: rsys,
						qus: qsys,
					},
				},
				u: using,
			}
		}),
		options: optionals(func(ext extent) (*pref.Options, error) {
			if err := using.Validate(); err != nil {
				return nil, err
			}

			if using.O != nil {
				return using.O, nil
			}

			return ext.options(settings...)
		}),
		navigator: kernel.Builder(func(o *pref.Options,
			resources *types.Resources,
		) (*kernel.Artefacts, error) {
			return kernel.New(using, o, &kernel.Benign{}, resources), nil
		}),
		plugins: activated(features),
	}
}

// Resume extent requests that the navigator performs a resume
// traversal, loading state from a previously saved session
// as a result of it being terminated prematurely via a ctrl-c
// interrupt.
func Resume(was *Was, settings ...pref.Option) *Builders {
	// TODO: the navigation file system, baseExtent.sys, will be set for
	// resume, only once the resume file has been loaded, as
	// its only at this point, we know where the original root
	// path was.
	//
	return &Builders{
		readerFS: pref.CreateReadDirFS(func() fs.ReadDirFS {
			if was.Using.GetReadDirFS != nil {
				return was.Using.GetReadDirFS()
			}
			return NewNativeFS(was.Root)
		}),
		queryFS: pref.CreateQueryStatusFS(func(fsys fs.FS) fs.StatFS {
			if was.Using.GetQueryStatusFS != nil {
				return was.Using.GetQueryStatusFS(fsys)
			}

			return NewQueryStatusFS(fsys)
		}),
		extent: extension(func(rsys fs.ReadDirFS, qsys fs.StatFS) extent {
			return &resumeExtent{
				baseExtent: baseExtent{
					fileSys: fileSystems{
						nas: rsys,
						qus: qsys,
					},
				},
				w: was,
			}
		}),
		// we need state; record the hibernation wake point, so
		// using a func here is probably not optimal.
		//
		options: optionals(func(ext extent) (*pref.Options, error) {
			if err := was.Validate(); err != nil {
				return nil, err
			}

			return ext.options(settings...)
		}),
		navigator: kernel.Builder(func(o *pref.Options,
			resources *types.Resources,
		) (*kernel.Artefacts, error) {
			return resume.NewController(was,
				kernel.New(&was.Using, o, resume.GetSealer(was), resources),
			), nil
		}),
		plugins: activated(features),
	}
}
