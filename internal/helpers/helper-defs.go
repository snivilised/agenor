package helpers

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/pref"
)

type (
	NaviTE struct {
		Given         string
		Should        string
		Relative      string
		Once          bool
		Visit         bool
		CaseSensitive bool
		Subscription  enums.Subscription
		Callback      core.Client
		Mandatory     []string
		Prohibited    []string
		ExpectedNoOf  Quantities
		ExpectedErr   error
	}

	FilterTE struct {
		NaviTE
		Description     string
		Pattern         string
		Scope           enums.FilterScope
		Negate          bool
		ErrorContains   string
		IfNotApplicable enums.TriStateBool
		Custom          core.TraverseFilter
		Type            enums.FilterType
		Sample          core.SampleTraverseFilter
	}

	PolyTE struct {
		NaviTE
		File   core.FilterDef
		Folder core.FilterDef
	}

	SampleTE struct {
		NaviTE
		SampleType enums.SampleType
		Reverse    bool
		Filter     *FilterTE
		NoOf       pref.EntryQuantities
		Each       pref.EachDirectoryEntryPredicate
		While      pref.WhileDirectoryPredicate
	}

	Quantities struct {
		Files    uint
		Folders  uint
		Children map[string]int
	}

	RecordingMap      map[string]int
	RecordingScopeMap map[string]enums.FilterScope
	RecordingOrderMap map[string]int
)
