package opts

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/pref"
)

// LoadInfo is the information required to load options.
type LoadInfo struct {
	O     *pref.Options
	State *core.ActiveState
}
