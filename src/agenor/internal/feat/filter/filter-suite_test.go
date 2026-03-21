package filter_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	lab "github.com/snivilised/jaywalk/src/agenor/internal/laboratory"
)

func TestFilter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filter Suite")
}

// FilterTE represents a test case for filter testing.
type FilterTE struct {
	lab.NaviTE
	// Description provides a human-readable explanation of the test case
	Description string

	// Pattern is the filter pattern to be applied during the test
	Pattern string

	// Scope defines the scope of the filter (e.g., file, directory, or both)
	Scope enums.FilterScope

	// Negate indicates whether the filter should be negated (i.e.,
	// exclude matches instead of including them)
	Negate bool

	// ErrorContains specifies a substring that should be present in the error
	// message if the filter fails
	ErrorContains string

	// IfNotApplicable indicates how the filter should behave if it is not applicable
	// to the current node (e.g., if the filter is designed for files but is
	// applied to a directory)
	IfNotApplicable enums.TriStateBool

	// Custom allows for a custom filter function to be provided for more complex
	// filtering logic that cannot be captured by the other fields
	Custom core.TraverseFilter

	// Type specifies the type of filter being tested (e.g., glob, regex, etc.)
	Type enums.FilterType

	// Sample provides a sample servant that can be used for testing the filter
	Sample core.SampleTraverseFilter
}

type PolyTE struct {
	lab.NaviTE
	File      core.FilterDef
	Directory core.FilterDef
}
