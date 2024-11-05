package kernel

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

func WithArtefacts(using *pref.Using,
	o *pref.Options,
	resources *enclave.Resources,
	sealer enclave.GuardianSealer,
) *Artefacts {
	controller := New(using, o, resources, sealer)
	mediator := controller.Mediator()

	return &Artefacts{
		Kontroller: controller,
		Mediator:   mediator,
		Resources:  resources,
	}
}

func New(facade pref.Facade, o *pref.Options,
	resources *enclave.Resources,
	sealer enclave.GuardianSealer,
) *NavigationController {
	impl, _ := newImpl(facade, o, resources)
	mediator := newMediator(&mediatorInfo{
		facade:    facade,
		o:         o,
		impl:      impl,
		sealer:    sealer,
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
