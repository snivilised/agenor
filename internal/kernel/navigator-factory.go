package kernel

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/pref"
)

func PrimeNav(using pref.Using, o *pref.Options) (core.Navigator, error) {
	return newController(&using, o)
}

func ResumeNav(with pref.As, o *pref.Options,
	resumption Resumption,
) (controller core.Navigator, err error) {
	controller, err = newController(&with.Using, o)

	if err != nil {
		return HadesNav(err)
	}

	return resumption.Decorate(controller), err
}

type Resumption interface {
	Decorate(core.Navigator) core.Navigator
}

type DecorateController func(core.Navigator) core.Navigator

func (f DecorateController) Decorate(source core.Navigator) core.Navigator {
	return f(source)
}

func newController(using *pref.Using,
	o *pref.Options,
) (navigator core.Navigator, err error) {
	if err = using.Validate(); err != nil {
		return
	}

	impl := newImpl(using, o)

	navigator = &navigationController{
		impl: impl,
		o:    o,
	}

	return
}

func newImpl(using *pref.Using,
	o *pref.Options,
) (navigator navigatorImpl) {
	base := navigatorBase{
		o: o,
	}

	switch using.Subscription { //nolint:exhaustive // already validated by using
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
	}

	return
}
