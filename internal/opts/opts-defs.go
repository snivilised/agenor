package opts

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/pref"
)

// 📦 pkg: opts - internal options handling; can't use persist
//
// TODO: The traverse-api table might be wrong as far as opts and persist
// is concerned. The table shows opts is above persist, yet we can't
// use persist from opts, so something is wrong and needs clarification.
//

type LoadInfo struct {
	O     *pref.Options
	State *core.ActiveState
}
