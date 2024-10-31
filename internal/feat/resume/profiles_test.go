package resume_test

import (
	"github.com/snivilised/traverse/enums"
)

const (
	Nothing                    = ""
	ResumeAtTeenageColor       = "RETRO-WAVE/College/Teenage Color"
	ResumeAtCanYouKissMeFirst  = "RETRO-WAVE/College/Teenage Color/A1 - Can You Kiss Me First.flac"
	StartAtElectricYouth       = "Electric Youth"
	StartAtBeforeLife          = "A1 - Before Life.flac"
	StartAtClientAlreadyActive = "this value doesn't matter"
)

type (
	strategyTheme struct {
		label string
	}

	strategyInvokeInfo struct {
		files   uint
		folders uint
	}

	strategyThemes      map[enums.ResumeStrategy]strategyTheme
	strategyInvocations map[enums.ResumeStrategy]strategyInvokeInfo

	resumeTestProfile struct {
		filtered   bool
		prohibited map[string]string
		mandatory  []string
	}

	profileThemes map[string]resumeTestProfile
)

var (
	prohibited = map[string]string{
		"RETRO-WAVE":                      Nothing,
		"Chromatics":                      Nothing,
		"Night Drive":                     Nothing,
		"A1 - The Telephone Call.flac":    Nothing,
		"A2 - Night Drive.flac":           Nothing,
		"cover.night-drive.jpg":           Nothing,
		"vinyl-info.night-drive.txt":      Nothing,
		"College":                         Nothing,
		"Northern Council":                Nothing,
		"A1 - Incident.flac":              Nothing,
		"A2 - The Zemlya Expedition.flac": Nothing,
		"cover.northern-council.jpg":      Nothing,
		"vinyl-info.northern-council.txt": Nothing,
	}
	filteredListenFlacs = []string{
		"A1 - Before Life.flac",
		"A2 - Runaway.flac",
	}
	filteredFlacs = []string{
		"A1 - Can You Kiss Me First.flac",
		"A2 - Teenage Color.flac",
		"A1 - Before Life.flac",
		"A2 - Runaway.flac",
	}
	textFiles = []string{
		"vinyl-info.teenage-color.txt",
		"vinyl-info.innerworld.txt",
	}
	strategies = []enums.ResumeStrategy{
		enums.ResumeStrategyFastward,
		enums.ResumeStrategySpawn,
	}

	profiles = profileThemes{
		// === Listening (uni/folder/file) (pend/active)

		"-> universal(pending): unfiltered": {
			filtered:   false,
			prohibited: prohibited,
			mandatory: append(append([]string{
				"Electric Youth",
				"Innerworld",
			}, filteredListenFlacs...), "vinyl-info.innerworld.txt"),
		},

		"-> universal(active): unfiltered": {
			filtered:   false,
			prohibited: prohibited,
			mandatory: append(append([]string{
				"Electric Youth",
				"Innerworld",
			}, filteredFlacs...), textFiles...),
		},

		"-> folders(pending): unfiltered": {
			filtered:   false,
			prohibited: prohibited,
			mandatory: []string{
				"Electric Youth",
				"Innerworld",
			},
		},

		"-> folders(active): unfiltered": {
			filtered:   false,
			prohibited: prohibited,
			mandatory: []string{
				"Teenage Color",
				"Electric Youth",
				"Innerworld",
			},
		},

		"-> files(pending): unfiltered": {
			filtered:   false,
			prohibited: prohibited,
			mandatory: []string{
				"A1 - Before Life.flac",
				"A2 - Runaway.flac",
				"vinyl-info.innerworld.txt",
			},
		},

		"-> files(active): unfiltered": {
			filtered:   false,
			prohibited: prohibited,
			mandatory:  append(filteredFlacs, textFiles...),
		},

		// === Filtering (uni/folder/file)

		"-> universal: filtered": {
			filtered:   true,
			prohibited: prohibited,
			mandatory: append([]string{
				"Electric Youth",
			}, filteredFlacs...),
		},

		"-> folders: filtered": {
			filtered:   true,
			prohibited: prohibited,
			mandatory: []string{
				"Electric Youth",
			},
		},

		"-> files: filtered": {
			filtered:   true,
			prohibited: prohibited,
			mandatory:  filteredFlacs,
		},

		// Listening and filtering (uni/folder/file)

		"-> universal: listen pending and filtered": {
			filtered:   true,
			prohibited: prohibited,
			mandatory: append([]string{
				"Electric Youth"}, filteredListenFlacs...),
		},

		"-> folders: listen pending and filtered": {
			filtered:   true,
			prohibited: prohibited,
			mandatory: []string{
				"Electric Youth",
			},
		},

		"-> files: listen pending and filtered": {
			filtered:   true,
			prohibited: prohibited,
			mandatory:  filteredListenFlacs,
		},
	}
)
