package opts

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/pref"
)

// Bind creates a Binder and binds the Events to the Controls
func Bind(o *pref.Options, active *core.ActiveState,
	settings ...pref.Option,
) (*LoadInfo, *Binder, error) {
	binder := NewBinder()
	o.Events.Bind(binder.Controls)

	err := apply(o, settings...)

	return &LoadInfo{
		O:     o,
		State: active,
	}, binder, err
}
