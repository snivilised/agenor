package tv

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type buildArtefacts struct {
	o       *pref.Options
	nav     core.Navigator
	plugins []types.Plugin
	ext     extent
}

type Builders struct {
	fs        pref.FsBuilder
	options   optionsBuilder
	navigator kernel.NavigatorBuilder
	plugins   pluginsBuilder
	extent    extentBuilder
}

func (bs *Builders) buildAll() (*buildArtefacts, error) {
	ext := bs.extent.build(bs.fs.Build())

	o, optionsErr := bs.options.build(ext)
	if optionsErr != nil {
		had := kernel.HadesNav(optionsErr)

		return &buildArtefacts{
			o:   o,
			nav: had,
			ext: ext,
		}, optionsErr
	}

	artefacts, navErr := bs.navigator.Build(o)
	if navErr != nil {
		had := kernel.HadesNav(navErr)

		return &buildArtefacts{
			o:   o,
			nav: had,
			ext: ext,
		}, navErr
	}

	plugins, pluginsErr := bs.plugins.build(o,
		artefacts.Mediator,
		ext.plugin(artefacts.Mediator),
	)

	if pluginsErr != nil {
		had := kernel.HadesNav(pluginsErr)

		return &buildArtefacts{
			o:   o,
			nav: had,
			ext: ext,
		}, pluginsErr
	}

	for _, p := range plugins {
		if bindErr := p.Init(); bindErr != nil {
			return &buildArtefacts{
				o:       o,
				nav:     artefacts.Navigator,
				plugins: plugins,
				ext:     ext,
			}, bindErr
		}
	}

	return &buildArtefacts{
		o:       o,
		nav:     artefacts.Navigator,
		plugins: plugins,
		ext:     ext,
	}, nil
}
