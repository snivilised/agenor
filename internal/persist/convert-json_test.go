package persist_test

import (
	"os"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"
	tv "github.com/snivilised/traverse"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/persist"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("Convert Options via JSON", Ordered, func() {
	var (
		fS tv.TraversalFS
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
