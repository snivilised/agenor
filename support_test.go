package tv_test

import (
	"errors"

	tv "github.com/snivilised/traverse"
)

var (
	errBuildOptions = errors.New("options build error")
)

const (
	RootPath    = "traversal-root-path"
	RestorePath = "/from-restore-path"
	files       = 3
	folders     = 2
)

var noOpHandler = func(_ *tv.Node) error {
	return nil
}
