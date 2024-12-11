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
	ActiveTE struct {
		Depth          int // must correspond to the correct depth of resumeAt
		ResumeAt       string
		HibernateState enums.Hibernation
	}

	GeneralTE struct {
		DescribedTE
		NaviTE
	}
	NaviTE struct {
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
		DescribedTE
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
		DescribedTE
		NaviTE
		Filter *pref.FilterOptions
	}

	HibernateTE struct {
		DescribedTE
		NaviTE
		Hibernate *core.HibernateOptions
		Mute      bool
	}

	HybridFilterTE struct {
		DescribedTE
		NaviTE
		NodeDef  core.FilterDef
		ChildDef core.ChildFilterDef
	}

	PolyTE struct {
		DescribedTE
		NaviTE
		File      core.FilterDef
		Directory core.FilterDef
	}

	ResumeTE struct {
		DescribedTE
		NaviTE
		Active         ActiveTE
		ClientListenAt string
		Profile        string
	}

	SampleTE struct {
		DescribedTE
		NaviTE
		SampleType enums.SampleType
		Reverse    bool
		Filter     *FilterTE
		NoOf       pref.EntryQuantities
		Each       pref.EachDirectoryEntryPredicate
		While      pref.WhileDirectoryPredicate
	}

	CascadeTE struct {
		DescribedTE
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
		Consume bool
	}

	AsyncPostage struct {
		Entry *AsyncOkTE
		Path  string
		FS    *luna.MemFS
		On    core.OutputFunc
	}

	BuilderFunc func(post *AsyncPostage) *age.Builders

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

func WithTestContext(specCtx SpecContext, fn func(context.Context, context.CancelFunc)) {
	defer leaktest.Check(GinkgoT())()

	ctx, cancel := context.WithCancel(specCtx)
	defer cancel()

	fn(ctx, cancel)
}

func FormatCascadeTestDescription(entry *CascadeTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.DescribedTE.Given, entry.DescribedTE.Should)
}

func FormatGeneralTestDescription(entry *GeneralTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.Given, entry.Should)
}

func FormatHibernateTestDescription(entry *HibernateTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.DescribedTE.Given, entry.DescribedTE.Should)
}

func FormatHybridFilterTestDescription(entry *HybridFilterTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.DescribedTE.Given, entry.DescribedTE.Should)
}

func FormatFilterTestDescription(entry *FilterTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: apply filter", entry.DescribedTE.Given)
}

func FormatFilterErrataTestDescription(entry *FilterErrataTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.DescribedTE.Given, entry.DescribedTE.Should)
}

func FormatPolyFilterTestDescription(entry *PolyTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.DescribedTE.Given, entry.DescribedTE.Should)
}

func FormatResumeTestDescription(entry *ResumeTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.DescribedTE.Given, entry.DescribedTE.Should)
}

func FormatSampleTestDescription(entry *SampleTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.DescribedTE.Given, entry.DescribedTE.Should)
}

type DescribedTestItem interface {
	WhenGiven() string
	ItShould() string
}

func FormatTestDescription(entry *DescribedTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.WhenGiven(), entry.ItShould())
}

type DescribedTE struct {
	Given  string
	Should string
}

func (e *DescribedTE) WhenGiven() string {
	return e.Given
}

func (e *DescribedTE) ItShould() string {
	return e.Should
}
