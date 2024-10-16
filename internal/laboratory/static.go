package lab

import (
	"io/fs"
)

var (
	Perms = struct {
		File fs.FileMode
		Dir  fs.FileMode
	}{
		File: 0o666, //nolint:mnd // ok (pedantic)
		Dir:  0o777, //nolint:mnd // ok (pedantic)
	}

	Static = struct {
		RetroWave string
	}{
		RetroWave: "RETRO-WAVE",
	}
)
