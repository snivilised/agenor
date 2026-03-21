package lab

import (
	"context"
	"fmt"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:staticcheck // ok

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/nefilim/test/luna"
)

type (
	// ActiveTE is used to capture the state of an active traversal at the point of resumption.
	ActiveTE struct {
		Depth          int               // must correspond to the correct depth of resumeAt
		ResumeAt       string            // the path at which to resume traversal
		HibernateState enums.Hibernation // the expected hibernation state at the point of resumption
	}

	// GeneralTE is a base struct for test entries that include a description of the test scenario.
	GeneralTE struct {
		DescribedTE
		NaviTE
	}

	// NaviTE captures common navigation-related parameters for test entries.
	NaviTE struct {
		// Relative is the path to be used in the test, relative to the root
		// of the traversal.
		Relative string

		// Once indicates whether the callback should be executed only once.
		Once bool

		// Visit indicates whether the traversal should actually visit the nodes or
		// just simulate the traversal.
		Visit bool

		// CaseSensitive indicates whether the path matching should be case-sensitive.
		CaseSensitive bool

		// Subscription is the type of subscription to be used in the test.
		Subscription enums.Subscription

		// Callback is the client callback to be invoked during the traversal.
		Callback core.Client

		// Mandatory lists the paths that must be included in the traversal for
		// the test to pass.
		Mandatory []string

		// Prohibited lists the paths that must not be included in the traversal
		// for the test to pass.
		Prohibited []string

		// ByPassMetrics indicates whether to bypass metrics collection during the traversal.
		ByPassMetrics bool

		// ExpectedNoOf captures the expected number of files, directories, and children
		// for the test scenario.
		ExpectedNoOf Quantities

		// ExpectedErr captures any expected error that should occur during the traversal.
		ExpectedErr error
	}

	// FilterTE captures parameters for tests related to traversal filters.
	FilterTE struct {
		DescribedTE
		NaviTE

		// Description provides a human-readable description of the filter being tested.
		Description string

		// Pattern is the filter pattern to be applied during the traversal.
		Pattern string

		// Scope defines the scope of the filter (e.g., files, directories, or both).
		Scope enums.FilterScope

		// Negate indicates whether the filter should be negated (i.e., exclude
		// matching entries).
		Negate bool

		// ErrorContains captures a substring that should be present in any error message
		// resulting from the filter application, if an error is expected.
		ErrorContains string

		// IfNotApplicable indicates whether the test should be skipped if the
		// filter is not applicable to the given path (e.g., if testing a file filter
		// on a directory path).
		IfNotApplicable enums.TriStateBool

		// Custom allows for any custom filter logic to be included in the test entry.
		Custom core.TraverseFilter

		// Type indicates the type of filter being tested (e.g., glob, regex, custom).
		Type enums.FilterType

		// Sample captures a sample traversal filter that can be used for testing purposes.
		Sample core.SampleTraverseFilter
	}

	// FilterErrataTE captures parameters for tests related to filter errata, which are
	// specific edge cases or exceptions in filter behavior.
	FilterErrataTE struct {
		DescribedTE
		NaviTE
		// Filter captures the filter configuration that is expected to produce errata
		// (i.e., unexpected results or errors) during traversal.
		Filter *pref.FilterOptions
	}

	// HibernateTE captures parameters for tests related to traversal hibernation, which is the
	// ability to pause and resume a traversal at a specific point.
	HibernateTE struct {
		DescribedTE
		NaviTE
		// Hibernate captures the hibernation options to be applied during the traversal, such as
		// the expected hibernation state at the point of resumption.
		Hibernate *core.HibernateOptions

		// Mute indicates whether to suppress output during the traversal, which can be useful for
		// tests that focus on hibernation behavior without concern for output.
		Mute bool
	}

	// HybridFilterTE captures parameters for tests that involve both file and
	// directory filters, allowing for complex filtering scenarios that may
	// involve interactions between file and directory matching.
	HybridFilterTE struct {
		DescribedTE
		NaviTE

		// NodeDef captures the filter definition for the current node being tested, which
		// may be either a file or a directory depending on the context of the test.
		NodeDef core.FilterDef

		// ChildDef captures the filter definition for the children of the current node,
		// which may also be either files or directories depending on the context of the
		// test. This allows for testing scenarios where the filtering behavior of a
		// node may depend on the characteristics of its children (e.g., a directory
		// filter that only matches if it contains certain files).
		ChildDef core.ChildFilterDef
	}

	// PolyTE captures parameters for tests that involve multiple filter definitions,
	// allowing for scenarios where both file and directory filters are applied
	// simultaneously or in combination.
	PolyTE struct {
		DescribedTE
		NaviTE
		// File captures the filter definition for file entries in the traversal,
		// which may be used to test scenarios where file filtering behavior is a
		// key aspect of the test.
		File core.FilterDef

		// Directory captures the filter definition for directory entries in
		// the traversal, which may be used to test scenarios where directory
		// filtering behavior is a key aspect of the test, or where the
		// interaction between file and directory filters is being evaluated.
		Directory core.FilterDef
	}

	// ResumeTE captures parameters for tests that involve resuming a traversal
	// from a specific point, allowing for scenarios where the traversal is
	// paused and then resumed, and the behavior at the point of resumption is a
	// key aspect of the test.
	ResumeTE struct {
		DescribedTE
		NaviTE
		// Active captures the state of the active traversal at the point of resumption,
		// including the depth, the path to resume at, and the expected hibernation state.
		Active ActiveTE

		// ClientListenAt captures the path at which the client should listen for
		// resumption, which may be used to test scenarios where the resumption point
		// is a specific path in the traversal, and the behavior of the traversal at
		// that point is being evaluated.
		ClientListenAt string

		// Profile captures the profile to be used during the traversal, which may be
		// relevant for tests that involve specific traversal configurations or
		// optimizations that are associated with certain profiles.
		Profile string
	}

	// SampleTE captures parameters for tests that involve sampling entries during
	// traversal, allowing for scenarios where only a subset of entries are
	// processed based on certain criteria, and the behavior of the sampling
	// mechanism is a key aspect of the test.
	SampleTE struct {
		DescribedTE
		NaviTE
		// SampleType indicates the type of sampling being tested (e.g., random,
		// systematic, stratified), which may be relevant for tests that evaluate
		// the behavior of different sampling strategies during traversal.
		SampleType enums.SampleType

		// SampleSize captures the size of the sample to be taken during the
		// traversal, which may be relevant for tests that evaluate the behavior
		// of sampling mechanisms with different sample sizes, and how that affects
		// the traversal results.
		Reverse bool

		// Filter captures the filter to be applied during sampling, which may be
		// relevant for tests that evaluate how sampling interacts with filtering
		// criteria during traversal.
		Filter *FilterTE

		// NoOf captures the expected number of entries that should be included in
		// the sample, which may be relevant for tests that evaluate the behavior of
		// sampling mechanisms in terms of how many entries are selected based on
		// the sampling criteria and the size of the sample.
		NoOf pref.EntryQuantities

		// Each captures any predicates that should be applied to each entry during the
		// sampling process, which may be relevant for tests that evaluate how specific
		// conditions or criteria affect the selection of entries during sampling.
		Each pref.EachDirectoryEntryPredicate

		// While captures any predicates that should be applied to the traversal while
		// sampling is active, which may be relevant for tests that evaluate how certain
		// conditions or criteria affect the traversal behavior during the sampling
		// process.
		While pref.WhileDirectoryPredicate
	}

	// CascadeTE captures parameters for tests that involve cascading filters, which
	// are filters that apply to both files and directories, and may have interactions
	// between the two types of entries. This allows for testing scenarios where the
	// filtering behavior of a node may depend on the characteristics of its children
	// (e.g., a directory filter that only matches if it contains certain files), and
	// where the interaction between file and directory filters is being evaluated.
	CascadeTE struct {
		DescribedTE
		NaviTE
		// NoRecurse indicates whether the traversal should avoid recursing into directories
		// that match the filter criteria, which may be relevant for tests that evaluate
		// the behavior of cascading filters in terms of how they affect the traversal
		// when certain entries are matched.
		NoRecurse bool

		// File captures the filter definition for file entries in the traversal, which
		// may be used to test scenarios where file filtering behavior is a key aspect
		// of the test.
		File core.FilterDef

		// Depth captures the expected depth of the traversal at which certain behaviors
		// should occur, which may be relevant for tests that evaluate how cascading
		// filters affect the traversal at different levels of the directory hierarchy.
		Depth uint
	}

	// AsyncResumeTE captures parameters for tests that involve asynchronous
	// traversal operations, allowing for scenarios where the traversal is
	// performed in a non-blocking manner, and the behavior of the traversal
	// under asynchronous conditions is a key aspect of the test.
	AsyncResumeTE struct {
		// At captures the path at which the traversal should be resumed, which may
		// be relevant for tests that evaluate the behavior of traversal
		// resumption at specific points in the directory hierarchy.
		At string

		// Strategy captures the resume strategy to be used during traversal resumption,
		// which may be relevant for tests that evaluate how different resume
		// strategies affect the behavior of the traversal when it is resumed.
		Strategy enums.ResumeStrategy

		// HibernationState captures the expected hibernation state at the
		// point of resumption, which may be relevant for tests that evaluate
		// how the traversal behaves when it is resumed from a hibernated state, and
		// how the hibernation state affects the traversal behavior at the point of
		// resumption.
		HibernationState enums.Hibernation
	}

	// AsyncTE captures parameters for tests that involve asynchronous traversal
	// operations, allowing for scenarios where the traversal is performed in a
	// non-blocking manner, and the behavior of the traversal under asynchronous
	// conditions is a key aspect of the test.
	AsyncTE struct {
		// Given sets out the preconditions or context for the
		// test, providing a description of the scenario being tested.
		Given string

		// Should sets out the expected outcome or behavior that should occur
		// during the test, providing a description of what the test is verifying.
		Should string

		// Path captures the path to be used in the test
		Path func() string

		// Subscription captures the type of subscription to be used in the test
		Subscription enums.Subscription

		// Callback captures the client callback to be invoked during the traversal
		Callback core.Client

		// Builder captures the function that builds the traversal configuration
		// or the test
		Builder BuilderFunc

		// File captures the filter definition for file entries in the traversal
		Resume AsyncResumeTE

		// NoWorkers captures the number of workers to be used during
		// the asynchronous traversal
		NoWorkers uint

		// CPU defines wether the number of workers used in a worker pool
		// should match the number of CPU cores available.
		CPU bool
	}

	// AsyncOkTE captures parameters for tests that involve asynchronous
	// traversal operations where the expectation is that the traversal
	// will complete successfully without errors,
	// allowing for scenarios where the focus is on verifying correct behavior
	// under asynchronous conditions.
	AsyncOkTE struct {
		AsyncTE
		// Consume indicates whether the test should consume the traversal results
		// from the output channel.
		Consume bool
	}

	// AsyncPostage used for async tests
	AsyncPostage struct {
		// Entry async test inputs
		Entry *AsyncOkTE

		// Path path spec
		Path string

		// FS filesystem
		FS *luna.MemFS

		// On output func
		On core.OutputFunc
	}

	// BuilderFunc captures the function that builds the traversal configuration
	// or the test
	BuilderFunc func(post *AsyncPostage) *agenor.Builders

	// CompositeTE captures parameters for tests that involve composite
	// traversal operations, allowing for scenarios where the traversal is
	// performed in a non-blocking manner, and the behavior of the traversal
	// under asynchronous conditions is a key aspect of the test.
	CompositeTE struct {
		AsyncTE
		// IsWalk indicates whether the traversal should be performed as a walk.
		IsWalk bool
		// IsPrime indicates whether the traversal should be performed as a prime.
		IsPrime bool
		// Facade captures the facade to be used in the traversal.
		Facade pref.Facade
	}

	// Quantities captures the quantities of entries in the traversal.
	Quantities struct {
		// Files captures the number of files in the traversal.
		Files uint
		// Directories captures the number of directories in the traversal.
		Directories uint
		// Children captures the number of children in the traversal.
		Children map[string]int
	}

	// MatcherExpectation captures the expected and actual values for a matcher.
	MatcherExpectation[T comparable] struct {
		Expected T
		Actual   T
	}

	// Recall captures the recall of entries in the traversal.
	Recall map[string]int
	// RecallScope captures the scope of the recall in the traversal.
	RecallScope map[string]enums.FilterScope
	// RecallOrder captures the order of the recall in the traversal.
	RecallOrder map[string]int
)

// IsEqual returns true if the expected and actual values are equal.
func (x MatcherExpectation[T]) IsEqual() bool {
	return x.Actual == x.Expected
}

// Trigger captures the metrics for a test.
type Trigger struct {
	Metrics core.Metrics
}

// Times sets the number of times a metric should be triggered.
func (t *Trigger) Times(m enums.Metric, n uint) *Trigger {
	t.Metrics[m].Times(n)

	return t
}

// WithTestContext creates a test context and runs the given function.
func WithTestContext(specCtx SpecContext, fn func(context.Context, context.CancelFunc)) {
	defer leaktest.Check(GinkgoT())()

	ctx, cancel := context.WithCancel(specCtx)
	defer cancel()

	fn(ctx, cancel)
}

// FormatCascadeTestDescription formats the test description for a cascade test.
func FormatCascadeTestDescription(entry *CascadeTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.Given, entry.Should)
}

// FormatGeneralTestDescription formats the test description for a general test.
func FormatGeneralTestDescription(entry *GeneralTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.Given, entry.Should)
}

// FormatHibernateTestDescription formats the test description for a hibernate test.
func FormatHibernateTestDescription(entry *HibernateTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.Given, entry.Should)
}

// FormatHybridFilterTestDescription formats the test description for a hybrid filter test.
func FormatHybridFilterTestDescription(entry *HybridFilterTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.Given, entry.Should)
}

// FormatFilterTestDescription formats the test description for a filter test.
func FormatFilterTestDescription(entry *FilterTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: apply filter", entry.Given)
}

// FormatFilterErrataTestDescription formats the test description for a filter errata test.
func FormatFilterErrataTestDescription(entry *FilterErrataTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.Given, entry.Should)
}

// FormatPolyFilterTestDescription formats the test description for a poly filter test.
func FormatPolyFilterTestDescription(entry *PolyTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.Given, entry.Should)
}

// FormatResumeTestDescription formats the test description for a resume test.
func FormatResumeTestDescription(entry *ResumeTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.Given, entry.Should)
}

// FormatSampleTestDescription formats the test description for a sample test.
func FormatSampleTestDescription(entry *SampleTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.Given, entry.Should)
}

// DescribedTestItem is an interface for test items that can be described.
type DescribedTestItem interface {
	// WhenGiven returns the description of the test case.
	WhenGiven() string
	// ItShould returns the expected outcome of the test case.
	ItShould() string
}

// FormatTestDescription formats the test description for a test item.
func FormatTestDescription(entry *DescribedTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.WhenGiven(), entry.ItShould())
}

// DescribedTE captures parameters for tests that can be described.
type DescribedTE struct {
	// Given captures the description of the test case.
	Given string
	// Should captures the expected outcome of the test case.
	Should string
}

// WhenGiven returns the description of the test case.
func (e *DescribedTE) WhenGiven() string {
	return e.Given
}

// ItShould returns the expected outcome of the test case.
func (e *DescribedTE) ItShould() string {
	return e.Should
}
