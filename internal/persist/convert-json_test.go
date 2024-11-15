package persist_test

import (
	"os"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	age "github.com/snivilised/agenor"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/opts"
	"github.com/snivilised/agenor/internal/persist"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("Convert Options via JSON", Ordered, func() {
	var (
		fS age.TraversalFS
	)

	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())
	})

	BeforeEach(func() {
		fS = &luna.MemFS{
			MapFS: fstest.MapFS{
				home: &fstest.MapFile{
					Mode: os.ModeDir,
				},
			},
		}

		_ = fS.MakeDirAll(destination, lab.Perms.Dir|os.ModeDir)
	})

	Context("ToJSON", func() {
		Context("given: source Options instance", func() {
			It("should: convert to JSON", func() {
				o, _, err := opts.Get(
					pref.WithDepth(4),
				)
				Expect(err).To(Succeed())
				Expect(persist.ToJSON(o)).To(HaveMarshaledEqual(o))
			})
		})
	})
})
