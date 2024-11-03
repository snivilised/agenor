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

// Builders performs build orchestration via its buildAll method. Builders
// is instructed by the factories (via Configure) of which there are 2; one
// for Walk and one for Run. The Prime/Resume extents create the Builders
// instance.
type Builders struct {
	using     *pref.Using
	forest    pref.ForestBuilder
	harvest   optionsBuilder
	navigator kernel.NavigatorBuilder
	plugins   pluginsBuilder
	extent    extentBuilder
}

func (bs *Builders) buildAll() (*buildArtefacts, error) {
	// BUILD FILE SYSTEM & EXTENT
	//
	ext := bs.extent.build(
		bs.forest.Build(bs.using.Tree),
	)

	// BUILD OPTIONS
	//
	harvest, optionsErr := bs.harvest.build(ext)
	if optionsErr != nil {
		return &buildArtefacts{
			o:   harvest.Options(),
			kc:  kernel.HadesNav(harvest.Options(), optionsErr),
			ext: ext,
		}, optionsErr
	}

	// BUILD NAVIGATOR
	//
	artefacts := bs.navigator.Build(harvest, &types.Resources{
		Forest:     ext.forest(),
		Supervisor: measure.New(),
		Binder:     harvest.Binder(),
	})

	// BUILD PLUGINS
	//
	plugins, pluginsErr := bs.plugins.build(harvest.Options(),
		bs.using,
		artefacts.Mediator,
		artefacts.Kontroller,
		ext.plugin(artefacts),
	)

	if pluginsErr != nil {
		return &buildArtefacts{
			o:   harvest.Options(),
			kc:  kernel.HadesNav(harvest.Options(), pluginsErr),
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
		O:          harvest.Options(),
		Kontroller: artefacts.Kontroller,
		Controls:   &harvest.Binder().Controls,
		Resources:  artefacts.Resources,
	}

	for _, p := range plugins {
		if bindErr := p.Init(pi); bindErr != nil {
			return &buildArtefacts{
				o:       harvest.Options(),
				kc:      artefacts.Kontroller,
				plugins: plugins,
				ext:     ext,
			}, bindErr
		}
	}

	return &buildArtefacts{
		o:       harvest.Options(),
		kc:      artefacts.Kontroller,
		plugins: plugins,
		ext:     ext,
	}, nil
}
