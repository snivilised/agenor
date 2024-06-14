package kernel

import (
	"github.com/snivilised/traverse/internal/level"
)

// navigationFrame represents info that relates to the navigator as
// a whole. Does not contain any information that relates to the current
// node as that is transient in nature (and already represented by
// navigationVapour)
type navigationFrame struct {
	periscope *level.Periscope
}
