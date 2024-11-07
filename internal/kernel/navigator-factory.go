package kernel

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

func PrimeArtefacts(using *pref.Using,
	harvest enclave.OptionHarvest,
	resources *enclave.Resources,
	sealer enclave.GuardianSealer,
) *Artefacts {
	ci := &enclave.ControllerInfo{
		Facade:    using,
		Harvest:   harvest,
		Resources: resources,
		Sealer:    sealer,
	}
	controller := New(ci)
	mediator := controller.Mediator()

	return &Artefacts{
		Kontroller: controller,
		Mediator:   mediator,
		Resources:  resources,
	}
}

func New(ci *enclave.ControllerInfo) *NavigationController {
	o := ci.Harvest.Options()
	facade := ci.Facade
	resources := ci.Resources
	impl, _ := newImpl(facade, o, resources)
	mediator := newMediator(&mediatorInfo{
		facade:    facade,
		o:         o,
		impl:      impl,
		sealer:    ci.Sealer,
		resources: resources,
	})

	return newNavigationController(mediator)
}

func newImpl(facade pref.Facade,
	o *pref.Options,
	resources *enclave.Resources,
) (impl NavigatorImpl, err error) {
	subscription := facade.Sub()

	agent := navigatorAgent{
		ao: &agentOptions{
			hooks:   &o.Hooks,
			defects: &o.Defects,
		},
		ro: &readOptions{
			hooks: readHooks{
				read: o.Hooks.ReadDirectory,
				sort: o.Hooks.Sort,
			},
			behaviour: &o.Behaviours.Sort,
		},
		resources: resources,
	}

	switch subscription {
	case enums.SubscribeFiles:
		impl = &navigatorFiles{
			navigatorAgent: agent,
		}

	case enums.SubscribeDirectories, enums.SubscribeDirectoriesWithFiles:
		impl = &navigatorDirectories{
			navigatorAgent: agent,
		}

	case enums.SubscribeUniversal:
		impl = &navigatorUniversal{
			navigatorAgent: agent,
		}

	case enums.SubscribeUndefined:
		err = locale.ErrUsageMissingSubscription
	}

	return impl, err
}
