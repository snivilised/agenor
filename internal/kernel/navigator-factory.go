package kernel

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func New(using *pref.Using, o *pref.Options,
	sealer types.GuardianSealer,
	resources *types.Resources,
) *Artefacts {
	impl := newImpl(using, o, resources)
	mediator := newMediator(&mediatorInfo{
		using:     using,
		o:         o,
		impl:      impl,
		sealer:    sealer,
		resources: resources,
	})
	controller := newNavigationController(mediator)

	return &Artefacts{
		Kontroller: controller,
		Mediator:   mediator,
		Resources:  resources,
	}
}

func newImpl(using *pref.Using,
	o *pref.Options,
	resources *types.Resources,
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

	case enums.SubscribeFolders, enums.SubscribeFoldersWithFiles:
		impl = &navigatorFolders{
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
