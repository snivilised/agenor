package traverse

import (
	"errors"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type buildArtefacts struct {
	o       *pref.Options
	nav     core.Navigator
	plugins []types.Plugin
}

type Builders struct {
	ob optionsBuilder
	nb navigatorBuilder
	pb pluginsBuilder
}

func (bs *Builders) buildAll() (*buildArtefacts, error) {
	o, optionsErr := bs.ob.build()
	if optionsErr != nil {
		return nil, optionsErr
	}

	nav, navErr := bs.nb.build(o)
	if navErr != nil {
		return nil, navErr
	}

	plugins, pluginsErr := bs.pb.build(o)
	if pluginsErr != nil {
		return nil, pluginsErr
	}

	if host, ok := nav.(types.UsePlugin); ok {
		es := []error{}
		for _, p := range plugins {
			registrationErr := host.Register(p)
			es = append(es, registrationErr)
		}

		if pluginErr := errors.Join(es...); pluginErr != nil {
			return nil, pluginErr
		}
	}

	return &buildArtefacts{
		o:       o,
		nav:     nav,
		plugins: plugins,
	}, nil
}
