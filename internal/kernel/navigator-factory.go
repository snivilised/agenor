package kernel

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/pref"
)

func WithArtefacts(using *pref.Using, o *pref.Options,
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

func New(using *pref.Using, o *pref.Options,
	resources *enclave.Resources,
	sealer enclave.GuardianSealer,
) *NavigationController {
	impl := newImpl(using, o, resources)
	mediator := newMediator(&mediatorInfo{
		using:     using,
		o:         o,
		impl:      impl,
		sealer:    sealer,
		resources: resources,
	})

	return newNavigationController(mediator)
}

func newImpl(using *pref.Using,
	o *pref.Options,
	resources *enclave.Resources,
) (impl NavigatorImpl) {
	agent := navigatorAgent{
		using: using,
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

	switch using.Subscription {
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
	}

	return impl
}
