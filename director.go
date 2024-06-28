package tv

import (
	"io/fs"
	"os"

	"github.com/snivilised/traverse/internal/hiber"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/refine"
	"github.com/snivilised/traverse/internal/resume"
	"github.com/snivilised/traverse/internal/sampling"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type ifActive func(o *pref.Options, mediator types.Mediator) types.Plugin

// activated interrogates options and invokes requests on behalf of the user
// to activate features according to option selections. other plugins will
// be initialised after primary plugins
func activated(o *pref.Options, mediator types.Mediator,
	others ...types.Plugin,
) (plugins []types.Plugin, err error) {
	var (
		all = []ifActive{
			hiber.IfActive, refine.IfActive, sampling.IfActive,
		}
	)

	plugins = []types.Plugin{}

	for _, active := range all {
		if plugin := active(o, mediator); plugin != nil {
			plugins = append(plugins, plugin)
		}
	}

	for _, plugin := range others {
		if plugin != nil {
			plugins = append(plugins, plugin)
		}
	}

	for _, plugin := range plugins {
		err = plugin.Register()

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
		filesystem: pref.FileSystem(func() fs.FS {
			if using.GetFS != nil {
				return using.GetFS()
			}
			return os.DirFS(using.Root)
		}),
		extent: extension(func(fsys fs.FS) extent {
			return &primeExtent{
				baseExtent: baseExtent{
					fsys: fileSystems{
						nas: fsys,
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
		navigator: kernel.Builder(func(o *pref.Options, res *types.Resources) (*kernel.Artefacts, error) {
			return kernel.New(using, o, &kernel.Benign{}, res), nil
		}),
		plugins: features(activated), // swap over features & activated
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
		filesystem: pref.FileSystem(func() fs.FS {
			if was.Using.GetFS != nil {
				return was.Using.GetFS()
			}
			return os.DirFS(was.Using.Root)
		}),
		extent: extension(func(fsys fs.FS) extent {
			return &resumeExtent{
				baseExtent: baseExtent{
					fsys: fileSystems{
						nas: fsys,
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
		navigator: kernel.Builder(func(o *pref.Options, res *types.Resources) (*kernel.Artefacts, error) {
			artefacts := kernel.New(&was.Using, o, resume.GetSealer(was), res)

			return resume.NewController(was, artefacts), nil
		}),
		plugins: features(activated),
	}
}
