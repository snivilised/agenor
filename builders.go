package tv

import (
	"errors"

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
	navigator navigatorBuilder
	plugins   pluginsBuilder
	extent    extent
}

func (bs *Builders) buildAll() (*buildArtefacts, error) {
	o, optionsErr := bs.options.build()
	if optionsErr != nil {
		had, _ := kernel.HadesNav(optionsErr)

		return &buildArtefacts{
			o:   o,
			nav: had,
		}, optionsErr
	}

	nav, navErr := bs.navigator.build(o)
	if navErr != nil {
		had, _ := kernel.HadesNav(navErr)

		return &buildArtefacts{
			o:   o,
			nav: had,
		}, navErr
	}

	plugins, pluginsErr := bs.plugins.build(o)
	if pluginsErr != nil {
		had, _ := kernel.HadesNav(pluginsErr)

		return &buildArtefacts{
			o:   o,
			nav: had,
		}, pluginsErr
	}

	if host, ok := nav.(types.UsePlugin); ok {
		es := []error{}
		for _, p := range plugins {
			registrationErr := host.Register(p)
			es = append(es, registrationErr)
		}

		if pluginErr := errors.Join(es...); pluginErr != nil {
			had, _ := kernel.HadesNav(pluginErr)

			return &buildArtefacts{
				o:   o,
				nav: had,
			}, pluginErr
		}
	}

	return &buildArtefacts{
		o:       o,
		nav:     nav,
		plugins: plugins,
	}, nil
}
