package kernel

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type facilities struct {
}

func (f *facilities) Inject(pref.ActiveState) {}

func New(using *pref.Using, o *pref.Options,
	sealer types.GuardianSealer,
	resources *types.Resources,
) *Artefacts {
	impl := newImpl(using, o, resources)
	controller := newController(using, o, impl, sealer, resources)

	return &Artefacts{
		Kontroller: controller,
		Mediator:   controller.Mediator,
		Resources:  resources,
	}
}

func newController(using *pref.Using,
	o *pref.Options,
	impl NavigatorImpl,
	sealer types.GuardianSealer,
	resources *types.Resources,
) *NavigationController {
	return &NavigationController{
		Mediator: newMediator(using, o, impl, sealer, resources),
	}
}

func newImpl(using *pref.Using,
	o *pref.Options,
	resources *types.Resources,
) (impl NavigatorImpl) {
	base := navigator{
		using:     using,
		o:         o,
		resources: resources,
	}

	switch using.Subscription {
	case enums.SubscribeFiles:
		impl = &navigatorFiles{
			navigator: base,
		}

	case enums.SubscribeFolders, enums.SubscribeFoldersWithFiles:
		impl = &navigatorFolders{
			navigator: base,
		}

	case enums.SubscribeUniversal:
		impl = &navigatorUniversal{
			navigator: base,
		}

	case enums.SubscribeUndefined:
	}

	return
}
