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
}

type Builders struct {
	options   optionsBuilder
	navigator kernel.NavigatorBuilder
	plugins   pluginsBuilder
	ext       extent
}

func (bs *Builders) buildAll() (*buildArtefacts, error) {
	o, optionsErr := bs.options.build(bs.ext)
	if optionsErr != nil {
		had := kernel.HadesNav(optionsErr)

		return &buildArtefacts{
			o:   o,
			nav: had,
		}, optionsErr
	}

	artefacts, navErr := bs.navigator.Build(o)
	if navErr != nil {
		had := kernel.HadesNav(navErr)

		return &buildArtefacts{
			o:   o,
			nav: had,
		}, navErr
	}

	plugins, pluginsErr := bs.plugins.build(o,
		artefacts.Mediator,
		bs.ext.plugin(artefacts.Mediator),
	)

	if pluginsErr != nil {
		had := kernel.HadesNav(pluginsErr)

		return &buildArtefacts{
			o:   o,
			nav: had,
		}, pluginsErr
	}

	for _, p := range plugins {
		if bindErr := p.Init(); bindErr != nil {
			return &buildArtefacts{
				o:       o,
				nav:     artefacts.Navigator,
				plugins: plugins,
			}, bindErr
		}
	}

	return &buildArtefacts{
		o:       o,
		nav:     artefacts.Navigator,
		plugins: plugins,
	}, nil
}
