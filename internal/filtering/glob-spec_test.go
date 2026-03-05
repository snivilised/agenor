package filtering

import (
	"fmt"
	"regexp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/li18ngo"
)

type (
	// GlobSpecTE is a test entry for glob spec tests.
	GlobSpecTE struct {
		// Given is the description of the test case.
		Given string
		// Should is the description of what the test case should do.
		Should string
		// Pattern is the glob pattern to test.
		Pattern string
		// Match is the expected match result.
		Match string
		// Not is true if the pattern should not match.
		Not bool
		// Base is the expected base part of the pattern
		Base string
		// Ext is the expected extension.
		Ext string
	}
)

// FormatGlobSpecTestDescription formats the glob spec test description.
func FormatGlobSpecTestDescription(entry *GlobSpecTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.Given, entry.Should)
}

var _ = Describe("GlobSpec", Ordered, func() {
	var re *regexp.Regexp

	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())

		re = regexp.MustCompile(GlobExPattern)
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("GlobExPattern",
		func(entry *GlobSpecTE) {
			spec, err := parse(entry.Pattern, re)
			Expect(err).To(Succeed(), "failed: parse")
			Expect(spec.base).To(Equal(entry.Base), "failed: base")
			Expect(spec.ext).To(Equal(entry.Ext), "failed: ext")
			Expect(spec.matcher).To(Equal(entry.Match), "failed: matcher")
			Expect(spec.excludeExt).To(Equal(entry.Not), "failed: excludeExt")
		},
		FormatGlobSpecTestDescription,
		Entry(nil, &GlobSpecTE{
			Given:   "single: any jpg",
			Should:  "match",
			Pattern: "*.jpg",
			Base:    "*",
			Ext:     "jpg",
			Match:   "*.jpg",
		}),

		Entry(nil, &GlobSpecTE{
			Given:   "single: any non jpg",
			Should:  "match",
			Pattern: "*.!jpg",
			Base:    "*",
			Ext:     "jpg",
			Match:   "*.jpg",
			Not:     true,
		}),

		Entry(nil, &GlobSpecTE{
			Given:   "single: any jpg starting with a",
			Should:  "match",
			Pattern: "a*.jpg",
			Base:    "a*",
			Ext:     "jpg",
			Match:   "a*.jpg",
		}),

		Entry(nil, &GlobSpecTE{
			Given:   "single: base with multiple wildcards",
			Should:  "match",
			Pattern: "g*m*a.jpg",
			Base:    "g*m*a",
			Ext:     "jpg",
			Match:   "g*m*a.jpg",
		}),

		Entry(nil, &GlobSpecTE{
			Given:   "single: base with dot",
			Should:  "match",
			Pattern: "d*l.a.jpg",
			Base:    "d*l.a",
			Ext:     "jpg",
			Match:   "d*l.a.jpg",
		}),

		Entry(nil, &GlobSpecTE{
			Given:   "single: empty base, any jpg",
			Should:  "match",
			Pattern: ".jpg",
			Base:    "",
			Ext:     "jpg",
			Match:   "*.jpg",
		}),
	)
})
