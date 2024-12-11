package filtering

import (
	"fmt"
	"regexp"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/li18ngo"
)

type (
	GlobSpecTE struct {
		Given   string
		Should  string
		Pattern string
		Match   string
		Not     bool
		Base    string
		Ext     string
	}
)

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
