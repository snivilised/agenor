package lab

import (
	"context"
	"fmt"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:stylecheck,revive // ok

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/nefilim/test/luna"
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

	AsyncResumeTE struct {
		At               string
		Strategy         enums.ResumeStrategy
		HibernationState enums.Hibernation
	}

	AsyncTE struct {
		Given        string
		Should       string
		Path         func() string
		Subscription enums.Subscription
		Callback     core.Client
		Builder      BuilderFunc
		Resume       AsyncResumeTE
		NoWorkers    uint
		CPU          bool
	}

	AsyncOkTE struct {
		AsyncTE
	}

	BuilderFunc func(entry *AsyncOkTE, path string, fS *luna.MemFS) *age.Builders

	CompositeTE struct {
		AsyncTE
		IsWalk  bool
		IsPrime bool
		Facade  pref.Facade
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

func WithTestContext(specCtx SpecContext, fn func(context.Context)) {
	defer leaktest.Check(GinkgoT())()

	ctx, cancel := context.WithCancel(specCtx)
	defer cancel()

	fn(ctx)
}

func FormatTestDescription(entry *NaviTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.Given, entry.Should)
}
