package opts

import (
	"github.com/snivilised/traverse/pref"
)

// 📦 pkg: opts - internal options handling

type ActiveState struct {
}

type LoadInfo struct {
	O      *pref.Options
	State  *ActiveState
	WakeAt string
}
