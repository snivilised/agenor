package helpers

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
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

	Quantities struct {
		Files    uint
		Folders  uint
		Children map[string]int
	}

	RecordingMap      map[string]int
	RecordingScopeMap map[string]enums.FilterScope
	RecordingOrderMap map[string]int
)
