package tv

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
)

type buildArtefacts struct {
	o       *pref.Options
	kc      enclave.KernelController
	plugins []enclave.Plugin
	ext     extent
}

// Builders performs build orchestration via its buildAll method. Builders
// is instructed by the factories (via Configure) of which there are 2; one
// for Walk and one for Run. The Prime/Resume extents create the Builders
// instance.
type Builders struct {
	facade    pref.Facade
	scaffold  scaffoldBuilder
	navigator kernel.NavigatorBuilder
	plugins   pluginsBuilder
	extent    extentBuilder
}

func (bs *Builders) buildAll(addons ...Addon) (*buildArtefacts, error) {
	// BUILD SCAFFOLD
	//
	scaffold, err := bs.scaffold.build(addons...)
	harvest := scaffold.harvest()
	ext := scaffold.extent()
	o := harvest.Options()

	if err != nil {
		return &buildArtefacts{
			o:   o,
			kc:  kernel.HadesNav(o, err),
			ext: ext,
		}, err
	}

	// BUILD NAVIGATOR
	//
	subscription := ext.subscription()
	artefacts := bs.navigator.Build(&kernel.Creation{
		Facade:       ext.facade(),
		Subscription: subscription,
		Harvest:      harvest,
		Resources: &enclave.Resources{
			Forest:     ext.forest(),
			Supervisor: core.NewSupervisor(),
			Binder:     harvest.Binder(),
		},
	})

	// BUILD PLUGINS
	//
	plugins, err := bs.plugins.build(o,
		ext,
		artefacts,
		ext.plugin(artefacts),
	)

	if err != nil {
		return &buildArtefacts{
			o:   o,
			kc:  kernel.HadesNav(o, err),
			ext: ext,
		}, err
	}

	// INIT PLUGINS
	//
	activeRoles := lo.Map(plugins,
		func(plugin enclave.Plugin, _ int) enums.Role {
			return plugin.Role()
		},
	)
	order := manifest(activeRoles)
	artefacts.Mediator.Arrange(activeRoles, order)

	pi := &enclave.PluginInit{
		O:          o,
		Kontroller: artefacts.Kontroller,
		Controls:   &harvest.Binder().Controls,
		Resources:  artefacts.Resources,
	}

	for _, p := range plugins {
		if err := p.Init(pi); err != nil {
			return &buildArtefacts{
				o:       o,
				kc:      artefacts.Kontroller,
				plugins: plugins,
				ext:     ext,
			}, err
		}
	}

	return &buildArtefacts{
		o:       o,
		kc:      artefacts.Kontroller,
		plugins: plugins,
		ext:     ext,
	}, nil
}
