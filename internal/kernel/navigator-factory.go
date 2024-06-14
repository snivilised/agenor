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
) *Artefacts {
	impl := newImpl(using, o)
	controller := newController(using, o, impl, sealer)

	return &Artefacts{
		Navigator: controller,
		Mediator:  controller.mediator,
	}
}

func newController(using *pref.Using,
	o *pref.Options,
	impl NavigatorImpl,
	sealer types.GuardianSealer,
) *NavigationController {
	return &NavigationController{
		mediator: &mediator{
			root:   using.Root,
			impl:   impl,
			client: newGuardian(using.Handler, sealer),
			frame: &navigationFrame{
				periscope: level.New(),
			},
			pad: newScratch(o),
			o:   o,
		},
	}
}

func newImpl(using *pref.Using,
	o *pref.Options,
) (navigator NavigatorImpl) {
	base := navigatorBase{
		using: using,
		o:     o,
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
