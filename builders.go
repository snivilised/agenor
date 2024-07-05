package tv

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/measure"
	"github.com/snivilised/traverse/pref"
)

type buildArtefacts struct {
	o           *pref.Options
	kc          types.KernelController
	plugins     []types.Plugin
	activeRoles []enums.Role
	ext         extent
}

type Builders struct {
	filesystem pref.FileSystemBuilder
	options    optionsBuilder
	navigator  kernel.NavigatorBuilder
	plugins    pluginsBuilder
	extent     extentBuilder
}

func (bs *Builders) buildAll() (*buildArtefacts, error) {
	// BUILD FILE SYSTEM & EXTENT
	//
	ext := bs.extent.build(bs.filesystem.Build())

	// BUILD OPTIONS
	//
	o, optionsErr := bs.options.build(ext)
	if optionsErr != nil {
		return &buildArtefacts{
			o:   o,
			kc:  kernel.HadesNav(optionsErr),
			ext: ext,
		}, optionsErr
	}

	// BUILD NAVIGATOR
	//
	artefacts, navErr := bs.navigator.Build(o, &types.Resources{
		FS: types.FileSystems{
			N: ext.navFS(),
			R: ext.resFS(),
		},
		Supervisor: measure.New(),
	})
	if navErr != nil {
		return &buildArtefacts{
			o:   o,
			kc:  kernel.HadesNav(navErr),
			ext: ext,
		}, navErr
	}

	// BUILD PLUGINS
	//
	plugins, pluginsErr := bs.plugins.build(o,
		artefacts.Mediator,
		artefacts.Kontroller,
		ext.plugin(artefacts),
	)

	if pluginsErr != nil {
		return &buildArtefacts{
			o:   o,
			kc:  kernel.HadesNav(pluginsErr),
			ext: ext,
		}, pluginsErr
	}

	// INIT PLUGINS
	//
	roles := lo.Map(plugins, func(plugin types.Plugin, _ int) enums.Role {
		return plugin.Role()
	})

	artefacts.Mediator.Arrange(roles)

	for _, p := range plugins {
		if bindErr := p.Init(); bindErr != nil {
			return &buildArtefacts{
				o:       o,
				kc:      artefacts.Kontroller,
				plugins: plugins,
				ext:     ext,
			}, bindErr
		}
	}

	return &buildArtefacts{
		o:           o,
		kc:          artefacts.Kontroller,
		plugins:     plugins,
		activeRoles: roles,
		ext:         ext,
	}, nil
}
