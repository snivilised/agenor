package lab

import (
	"github.com/snivilised/agenor/core"
)

var (
	Perms = core.Permissions{
		File: 0o666, //nolint:mnd // ok (pedantic)
		Dir:  0o777, //nolint:mnd // ok (pedantic)
	}

	Static = struct {
		JSONFile        string
		JSONSubPath     string
		ElectricYouth   string
		NorthernCouncil string
		TeenageColor    string
		RetroWave       string
	}{
		JSONFile:        "resume-state.json",
		JSONSubPath:     "json/unmarshal/resume-state.json",
		ElectricYouth:   "Electric Youth",
		NorthernCouncil: "Northern Council",
		RetroWave:       "RETRO-WAVE",
		TeenageColor:    "Teenage Color",
	}
)
