package opts

import (
	"io/fs"

	"github.com/snivilised/traverse/pref"
)

func Load(fS fs.FS, from string, settings ...pref.Option) (*LoadInfo, *Binder, error) {
	o := pref.DefaultOptions()
	binder := NewBinder()
	o.Events.Bind(&binder.Controls)

	file, err := fS.Open(from)
	if err != nil {
		return &LoadInfo{
			O:      o,
			WakeAt: "tbd",
		}, binder, err
	}
	defer file.Close()

	// TODO: save any active state on the binder, eg the wake point

	err = apply(o, settings...)

	return &LoadInfo{
		O:      o,
		WakeAt: "tbd",
	}, binder, err
}
