package tv_test

import (
	"errors"
)

var (
	errBuildOptions = errors.New("options build error")
)

const (
	RootPath    = "/traversal-root-path"
	RestorePath = "/from-restore-path"
	files       = 3
	folders     = 2
)
