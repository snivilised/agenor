package opts

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/pref"
)

func Bind(o *pref.Options, active *core.ActiveState,
	settings ...pref.Option,
) (*LoadInfo, *Binder, error) {
	binder := NewBinder()
	o.Events.Bind(&binder.Controls)

	err := apply(o, settings...)

	return &LoadInfo{
		O:     o,
		State: active,
	}, binder, err
}
