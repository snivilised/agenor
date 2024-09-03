package sampling_test

import (
	"io/fs"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/filtering"
)

func TestSampling(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sampling Suite")
}

// customSamplingFilter is a custom sampling filter that just happens
// to use a glob as part of its implementation. The client can of course
// define their own custom implementation using filter.SampleFilter.
type customSamplingFilter struct {
	filtering.Sample
	description string
	pattern     string
}

func (f *customSamplingFilter) Description() string {
	return f.description
}

func (f *customSamplingFilter) Scope() enums.FilterScope {
	return f.Sample.Scope()
}

func (f *customSamplingFilter) Matching(children []fs.DirEntry) []fs.DirEntry {
	return f.Sample.Matching(children,
		func(entry fs.DirEntry, _ int) bool {
			matched, _ := filepath.Match(f.pattern, entry.Name())
			return matched
		},
	)
}
