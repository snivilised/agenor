package kernel

import (
	"fmt"
	"strings"

	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/persist"
	"github.com/snivilised/traverse/pref"
)

type author struct {
	o     *pref.Options
	perms core.Permissions
}

func (a *author) write(v vexation) (string, error) {
	vapour := v.vapour()
	s := vapour.static()
	forest := s.mediator.resources.Forest
	fS := forest.R
	calc := fS.Calc()
	directory, file := a.destination(v, calc)

	if err := fS.MakeDirAll(directory, a.perms.Dir); err != nil {
		return "", err
	}

	path := calc.Join(directory, file)
	active := vapour.active(s.tree, forest,
		s.mediator.periscope.Depth(),
		s.mediator.metrics,
	)

	request := &persist.MarshalRequest{
		Active: active,
		O:      a.o,
		Path:   path,
		Perm:   a.perms.File,
		FS:     fS,
	}

	_, err := persist.Marshal(request)

	return path, err
}

func (a *author) destination(v vexation, calc nef.PathCalc) (directory, file string) {
	extent := v.extent()
	cause := v.cause()
	now := core.Now()

	directory = nef.ResolvePath(a.o.Monitor.Admin.Path)
	if !strings.HasSuffix(directory, core.ResumeTail) {
		directory = calc.Join(directory, core.ResumeTail)
	}

	file = fmt.Sprintf("%v.%v.%v.json",
		extent, cause, now.Format(core.FileSystemTimeFormat),
	)

	return directory, file
}
