package kernel

import (
	"errors"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/pref"
)

func Prime(using core.Using, o *pref.Options) (core.Navigator, error) {
	return newController(&using, o)
}

func Resume(with core.As, o *pref.Options, resumption Resumption) (core.Navigator, error) {
	controller, err := newController(&with.Using, o)

	return resumption.Decorate(controller), err
}

type Resumption interface {
	Decorate(core.Navigator) core.Navigator
}

type DecorateController func(core.Navigator) core.Navigator

func (f DecorateController) Decorate(source core.Navigator) core.Navigator {
	return f(source)
}

func newController(using *core.Using, o *pref.Options) (core.Navigator, error) {
	if err := using.Validate(); err != nil {
		return nil, err
	}

	var (
		impl core.Navigator
		err  error
	)

	impl, err = newImpl(using, o)

	navigator := &navigationController{
		impl: impl,
		o:    o,
	}

	return navigator, err
}

func newImpl(using *core.Using, o *pref.Options) (core.Navigator, error) {
	var (
		navigator    core.Navigator
		err          error
		subscription = using.Subscription
	)

	switch subscription {
	case enums.SubscribeFiles:
		navigator = &navigationController{
			o: o,
		} // just temporary (create the impl's)
	case enums.SubscribeFolders:
		navigator = &navigationController{}
	case enums.SubscribeFoldersWithFiles:
		navigator = &navigationController{}
	case enums.SubscribeUniversal:
		navigator = &navigationController{}
	case enums.SubscribeUndefined:
		err = errors.New("invalid subscription")
	}

	return navigator, err
}
