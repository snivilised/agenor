package tv

import (
	"io/fs"

	"github.com/snivilised/traverse/internal/feat/resume"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type extent interface {
	using() *pref.Using
	was() *pref.Was
	plugin(*kernel.Artefacts) types.Plugin
	options(...pref.Option) (*pref.Options, error)
	navFS() fs.ReadDirFS
	queryFS() fs.StatFS
	resFS() fs.FS
	complete() bool
}

type fileSystems struct {
	nas fs.ReadDirFS
	qus fs.StatFS
	res fs.FS
}

type baseExtent struct {
	fileSys fileSystems
}

func (ex *baseExtent) navFS() fs.ReadDirFS {
	return ex.fileSys.nas
}

func (ex *baseExtent) queryFS() fs.StatFS {
	return ex.fileSys.qus
}

func (ex *baseExtent) resFS() fs.FS {
	return ex.fileSys.nas
}

type primeExtent struct {
	baseExtent
	u *pref.Using
}

func (ex *primeExtent) using() *pref.Using {
	return ex.u
}

func (ex *primeExtent) was() *pref.Was {
	return nil
}

func (ex *primeExtent) plugin(*kernel.Artefacts) types.Plugin {
	return nil
}

func (ex *primeExtent) options(settings ...pref.Option) (*pref.Options, error) {
	return pref.Get(settings...)
}

func (ex *primeExtent) complete() bool {
	return true
}

type resumeExtent struct {
	baseExtent
	w      *pref.Was
	loaded *pref.LoadInfo
	rp     *resume.Plugin
}

func (ex *resumeExtent) using() *pref.Using {
	return &ex.w.Using
}

func (ex *resumeExtent) was() *pref.Was {
	return ex.w
}

func (ex *resumeExtent) plugin(artefacts *kernel.Artefacts) types.Plugin {
	ex.rp = &resume.Plugin{
		BasePlugin: kernel.BasePlugin{
			Mediator: artefacts.Mediator,
		},
		IfResult: artefacts.IfResult,
	}

	return ex.rp
}

func (ex *resumeExtent) options(settings ...pref.Option) (*pref.Options, error) {
	loaded, err := resume.Load(ex.fileSys.res, ex.w.From, settings...)
	ex.loaded = loaded

	// get the resume point from the resume persistence file
	// then set up hibernation with this defined as a hibernation
	// filter.
	//
	return loaded.O, err
}

func (ex *resumeExtent) complete() bool {
	return ex.rp.IfResult.IsComplete()
}
