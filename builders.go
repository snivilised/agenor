package tv

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/internal/override"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/measure"
	"github.com/snivilised/traverse/pref"
)

type buildArtefacts struct {
	o       *pref.Options
	kc      types.KernelController
	plugins []types.Plugin
	ext     extent
}

type Builders struct {
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
			kc:  kernel.HadesNav(optionsErr),
			ext: ext,
		}, optionsErr
	}

	// BUILD NAVIGATOR
	//
	actions := &override.Actions{
		HandleChildren: override.NewActionCtrl[override.HandleChildrenInterceptor](
			func(inspection core.Inspection, mums measure.MutableMetrics) {
				// [KEEP-FILTER-IN-SYNC] keep this in sync with filter plugin.Init
				files := inspection.Sort(enums.EntryTypeFile)
				inspection.AssignChildren(files)
				mums[enums.MetricNoChildFilesFound].Times(uint(len(files)))
			},
		),
	}

	artefacts, navErr := bs.navigator.Build(o, &types.Resources{
		FS: FileSystems{
			N: ext.navFS(),
			Q: ext.queryFS(),
			R: ext.resFS(),
		},
		Supervisor: measure.New(),
		Actions:    actions,
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
	artefacts.Mediator.Arrange(lo.Map(plugins,
		func(plugin types.Plugin, _ int) enums.Role {
			return plugin.Role()
		},
	))

	for _, p := range plugins {
		if bindErr := p.Init(&types.PluginInit{
			Actions: actions,
			O:       o,
		}); bindErr != nil {
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
