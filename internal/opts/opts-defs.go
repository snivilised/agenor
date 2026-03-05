package opts

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/pref"
)

// LoadInfo is the information required to load options.
type LoadInfo struct {
	O     *pref.Options
	State *core.ActiveState
}
