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
				// [KEEP-FILTER-IN-SYNC] keep this in sync with filter-plugin/childScheme.init
				files := inspection.Sort(enums.EntryTypeFile)
				inspection.AssignChildren(files)

				// The behaviour of this default child handler is to assign
				// the children and tick the metrics. However, when filtering is
				// active, then this handler should be overridden by the filter
				// to only tick the child metric, if their parent is filtered in.
				//
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
	active := lo.Map(plugins,
		func(plugin types.Plugin, _ int) enums.Role {
			return plugin.Role()
		},
	)
	order := manifest(active)
	artefacts.Mediator.Arrange(active, order)
	pi := &types.PluginInit{
		Actions: actions,
		O:       o,
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
