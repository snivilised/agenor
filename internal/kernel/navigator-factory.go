package kernel

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/level"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type facilities struct {
}

func (f *facilities) Inject(pref.ActiveState) {}

func New(using *pref.Using, o *pref.Options,
	sealer types.GuardianSealer,
	res *types.Resources,
) *Artefacts {
	impl := newImpl(using, o, res)
	controller := newController(using, o, impl, sealer, res)

	return &Artefacts{
		Navigator: controller,
		Mediator:  controller.mediator,
		Resources: res,
	}
}

func newController(using *pref.Using,
	o *pref.Options,
	impl NavigatorImpl,
	sealer types.GuardianSealer,
	res *types.Resources,
) *NavigationController {
	return &NavigationController{
		mediator: &mediator{
			root:   using.Root,
			impl:   impl,
			client: newGuardian(using.Handler, sealer),
			frame: &navigationFrame{
				periscope: level.New(),
			},
			pad:       newScratch(o),
			o:         o,
			resources: res,
		},
	}
}

func newImpl(using *pref.Using,
	o *pref.Options,
	res *types.Resources,
) (navigator NavigatorImpl) {
	base := navigatorBase{
		using: using,
		o:     o,
		res:   res,
	}

	switch using.Subscription {
	case enums.SubscribeFiles:
		navigator = &navigatorFiles{
			navigatorBase: base,
		}

	case enums.SubscribeFolders, enums.SubscribeFoldersWithFiles:
		navigator = &navigatorFolders{
			navigatorBase: base,
		}

	case enums.SubscribeUniversal:
		navigator = &navigatorUniversal{
			navigatorBase: base,
		}

	case enums.SubscribeUndefined:
	}

	return
}
