package pref

import (
	"github.com/snivilised/traverse/core"
)

type HibernateOptions struct {
	Wake  *core.FilterDef
	Sleep *core.FilterDef
}
