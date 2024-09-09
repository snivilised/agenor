package tv

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/measure"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type buildArtefacts struct {
	o       *pref.Options
	kc      types.KernelController
	plugins []types.Plugin
	ext     extent
}

type Builders struct {
	using     *pref.Using
	readerFS  pref.ReadDirFileSystemBuilder
	queryFS   pref.QueryStatusFileSystemBuilder
	options   optionsBuilder
	navigator kernel.NavigatorBuilder
	plugins   pluginsBuilder
	extent    extentBuilder
}

func (bs *Builders) buildAll() (*buildArtefacts, error) {
	// BUILD FILE SYSTEM & EXTENT
	//
	reader := bs.readerFS.Build()
	ext := bs.extent.build(
		reader,
		bs.queryFS.Build(reader),
	)

	// BUILD OPTIONS
	//
	o, optionsErr := bs.options.build(ext)
	if optionsErr != nil {
		return &buildArtefacts{
			o:   o,
			kc:  kernel.HadesNav(o, optionsErr),
			ext: ext,
		}, optionsErr
	}

	// BUILD NAVIGATOR
	//
	artefacts, navErr := bs.navigator.Build(o, &types.Resources{
		FS: FileSystems{
			N: ext.navFS(),
			Q: ext.queryFS(),
			R: ext.resFS(),
		},
		Supervisor: measure.New(),
	})

	if navErr != nil {
		return &buildArtefacts{
			o:   o,
			kc:  kernel.HadesNav(o, navErr),
			ext: ext,
		}, navErr
	}

	// BUILD PLUGINS
	//
	plugins, pluginsErr := bs.plugins.build(o,
		bs.using,
		artefacts.Mediator,
		artefacts.Kontroller,
		ext.plugin(artefacts),
	)

	if pluginsErr != nil {
		return &buildArtefacts{
			o:   o,
			kc:  kernel.HadesNav(o, pluginsErr),
			ext: ext,
		}, pluginsErr
	}

	// INIT PLUGINS
	//
	active := lo.Map(plugins,
		func(plugin types.Plugin, _ int) enums.Role {
			return plugin.Role()
		},
	)
	order := manifest(active)
	artefacts.Mediator.Arrange(active, order)
	pi := &types.PluginInit{
		O:        o,
		Controls: artefacts.Mediator.Controls(),
	}

	for _, p := range plugins {
		if bindErr := p.Init(pi); bindErr != nil {
			return &buildArtefacts{
				o:       o,
				kc:      artefacts.Kontroller,
				plugins: plugins,
				ext:     ext,
			}, bindErr
		}
	}

	return &buildArtefacts{
		o:       o,
		kc:      artefacts.Kontroller,
		plugins: plugins,
		ext:     ext,
	}, nil
}
