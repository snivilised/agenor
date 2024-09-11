package opts

import (
	"io/fs"

	"github.com/snivilised/traverse/pref"
)

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
