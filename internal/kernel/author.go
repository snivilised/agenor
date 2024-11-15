package kernel

import (
	"fmt"
	"strings"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/internal/persist"
	"github.com/snivilised/agenor/pref"
	nef "github.com/snivilised/nefilim"
)

type author struct {
	o     *pref.Options
	perms core.Permissions
}

func (a *author) write(vex vexation) (string, error) {
	vapour := vex.vapour()
	static := vapour.static()
	forest := static.mediator.resources.Forest
	fS := forest.R
	calc := fS.Calc()
	directory, file := a.destination(vex, calc)

	if err := fS.MakeDirAll(directory, a.perms.Dir); err != nil {
		return "", err
	}

	path := calc.Join(directory, file)
	active := vapour.active(vex.ancestor(), forest,
		static.mediator.periscope.Depth(),
		static.mediator.metrics,
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
	mag := v.magnitude()
	cause := v.cause()
	now := core.Now()

	directory = nef.ResolvePath(a.o.Monitor.Admin.Path)
	if !strings.HasSuffix(directory, core.ResumeTail) {
		directory = calc.Join(directory, core.ResumeTail)
	}

	file = fmt.Sprintf("%v.%v.%v.json",
		mag, cause, now.Format(core.FileSystemTimeFormat),
	)

	return directory, file
}
