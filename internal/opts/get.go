package opts

import (
	"io/fs"

	"github.com/snivilised/traverse/pref"
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

func GetWith(o *pref.Options) *Binder {
	binder := NewBinder()
	o.Events.Bind(&binder.Controls)

	return binder
}

type ActiveState struct {
}

type LoadInfo struct {
	O      *pref.Options
	State  *ActiveState
	WakeAt string
}

// put this into file load
func Load(_ fs.FS, from string, settings ...pref.Option) (*LoadInfo, *Binder, error) {
	o := pref.DefaultOptions()
	// do load
	_ = from
	binder := NewBinder()
	o.Events.Bind(&binder.Controls)

	// TODO: save any active state on the binder, eg the wake point

	err := apply(o, settings...)

	return &LoadInfo{
		O:      o,
		WakeAt: "tbd",
	}, binder, err
}
