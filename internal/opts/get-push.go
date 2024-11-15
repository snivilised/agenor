package opts

import (
	"github.com/snivilised/agenor/pref"
)

func Get(settings ...pref.Option) (o *pref.Options, b *Binder, err error) {
	o = pref.DefaultOptions()

	binder := NewBinder()
	o.Events.Bind(&binder.Controls)

	err = apply(o, settings...)

	return o, binder, err
}

func apply(o *pref.Options, settings ...pref.Option) (err error) {
	for _, option := range settings {
		if option != nil {
			err = option(o)

			if err != nil {
				return err
			}
		}
	}

	return err
}

func Push(o *pref.Options) *Binder {
	binder := NewBinder()
	o.Events.Bind(&binder.Controls)

	return binder
}
