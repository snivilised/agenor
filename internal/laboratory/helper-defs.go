package lab

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/pref"
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
		ByPassMetrics bool
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

	FilterErrataTE struct {
		NaviTE
		Filter *pref.FilterOptions
	}

	HybridFilterTE struct {
		NaviTE
		NodeDef  core.FilterDef
		ChildDef core.ChildFilterDef
	}

	PolyTE struct {
		NaviTE
		File      core.FilterDef
		Directory core.FilterDef
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

	CascadeTE struct {
		NaviTE
		NoRecurse bool
		Depth     uint
	}

	Quantities struct {
		Files       uint
		Directories uint
		Children    map[string]int
	}

	MatcherExpectation[T comparable] struct {
		Expected T
		Actual   T
	}

	Recall      map[string]int
	RecallScope map[string]enums.FilterScope
	RecallOrder map[string]int
)

func (x MatcherExpectation[T]) IsEqual() bool {
	return x.Actual == x.Expected
}

type Trigger struct {
	Metrics core.Metrics
}

func (t *Trigger) Times(m enums.Metric, n uint) *Trigger {
	t.Metrics[m].Times(n)

	return t
}
