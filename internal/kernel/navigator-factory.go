package kernel

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

func PrimeArtefacts(creation *Creation,
	sealer enclave.GuardianSealer,
) *Artefacts {
	controller := New(creation, sealer)
	mediator := controller.Mediator()

	return &Artefacts{
		Kontroller: controller,
		Mediator:   mediator,
		Resources:  creation.Resources,
	}
}

func New(creation *Creation,
	sealer enclave.GuardianSealer,
) *NavigationController {
	o := creation.Harvest.Options()
	facade := creation.Facade
	resources := creation.Resources
	impl, _ := newImpl(o, creation)
	mediator := newMediator(&mediatorInfo{
		facade:       facade,
		subscription: creation.Subscription,
		o:            o,
		impl:         impl,
		sealer:       sealer,
		resources:    resources,
	})

	return newNavigationController(mediator)
}

func newImpl(o *pref.Options,
	creation *Creation,
) (impl NavigatorImpl, err error) {
	subscription := creation.Subscription

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
		resources: creation.Resources,
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
