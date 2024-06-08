package tv

import (
	"github.com/snivilised/traverse/pref"
)

type extent interface {
	using() *pref.Using
	was() *pref.Was
}

type baseExtent struct {
}

type primeExtent struct {
	baseExtent
	u *pref.Using
}

func (pe *primeExtent) using() *pref.Using {
	return pe.u
}

func (pe *primeExtent) was() *pref.Was {
	return nil
}

type resumeExtent struct {
	baseExtent
	w *pref.Was
}

func (re *resumeExtent) using() *pref.Using {
	return &re.w.Using
}

func (re *resumeExtent) was() *pref.Was {
	return re.w
}
