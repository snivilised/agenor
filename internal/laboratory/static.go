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
		JSONFile    string
		JSONSubPath string
		RetroWave   string
	}{
		JSONFile:    "resume-state.json",
		JSONSubPath: "json/unmarshal/resume-state.json",
		RetroWave:   "RETRO-WAVE",
	}
)
