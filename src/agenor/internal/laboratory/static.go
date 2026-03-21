package lab

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
)

var (
	// Perms defines the file and directory permissions used in the laboratory
	// tests. The file permissions are set to 0o666 (read and write for everyone),
	// and the directory permissions are set to 0o777 (read, write, and execute
	// for everyone). These permissions are used to ensure that the tests can
	// create and modify files and directories as needed without encountering
	// permission issues.
	Perms = core.Permissions{
		File: 0o666, //nolint:mnd // ok (pedantic)
		Dir:  0o777, //nolint:mnd // ok (pedantic)
	}

	// Static defines the static values used in the laboratory tests.
	Static = struct {
		// JSONFile is the name of the JSON file used in the tests.
		JSONFile string

		// JSONSubPath is the subpath to the JSON file used in the tests.
		JSONSubPath string

		// ElectricYouth is a string representing the name of the "Electric Youth" band.
		ElectricYouth string

		// NorthernCouncil is a string representing the name of the "Northern Council" band.
		NorthernCouncil string

		// TeenageColor is a string representing the name of the "Teenage Color" band.
		TeenageColor string

		// RetroWave is a string representing the name of the "RETRO-WAVE" band.
		RetroWave string
	}{
		JSONFile:        "resume-state.json",
		JSONSubPath:     "json/unmarshal/resume-state.json",
		ElectricYouth:   "Electric Youth",
		NorthernCouncil: "Northern Council",
		RetroWave:       "RETRO-WAVE",
		TeenageColor:    "Teenage Color",
	}
)
