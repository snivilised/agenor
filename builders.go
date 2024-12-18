package age

import (
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/kernel"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/pref"
)

type buildArtefacts struct {
	o         *pref.Options
	kc        enclave.KernelController
	plugins   []enclave.Plugin
	ext       extent
	swappable enclave.Swapper
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
		o.Monitor.Log.Error(err.Error())

		return &buildArtefacts{
			o:   o,
			kc:  kernel.HadesNav(o, err),
			ext: ext,
		}, err
	}

	// BUILD NAVIGATOR
	//
	artefacts := bs.navigator.Build(&kernel.Inception{
		Facade:       ext.facade(),
		Subscription: ext.subscription(),
		Harvest:      harvest,
		Resources: &enclave.Resources{
			Forest:     ext.forest(),
			Supervisor: enclave.NewSupervisor(),
			Binder:     harvest.Binder(),
		},
	})

	if artefacts.Error != nil {
		return &buildArtefacts{
			o:         o,
			kc:        kernel.HadesNav(o, err),
			ext:       ext,
			swappable: artefacts.Mediator,
		}, artefacts.Error
	}

	// BUILD PLUGINS
	//
	plugins, err := bs.plugins.build(o,
		ext,
		artefacts,
		ext.plugin(artefacts),
	)

	if err != nil {
		o.Monitor.Log.Error(err.Error())

		return &buildArtefacts{
			o:         o,
			kc:        kernel.HadesNav(o, err),
			ext:       ext,
			swappable: artefacts.Mediator,
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
		Controls:   harvest.Binder().Controls,
		Resources:  artefacts.Resources,
	}

	for _, p := range plugins {
		if err := p.Init(pi); err != nil {
			o.Monitor.Log.Error(err.Error())

			return &buildArtefacts{
				o:         o,
				kc:        artefacts.Kontroller,
				plugins:   plugins,
				ext:       ext,
				swappable: artefacts.Mediator,
			}, err
		}
	}

	return &buildArtefacts{
		o:         o,
		kc:        artefacts.Kontroller,
		plugins:   plugins,
		ext:       ext,
		swappable: artefacts.Mediator,
	}, nil
}
