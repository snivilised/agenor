package nfs_test

import (
	"fmt"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/traverse/nfs"
)

var _ = Describe("ResolvePath", func() {
	DescribeTable("Overrides provided",
		func(entry *RPEntry) {
			mocks := nfs.ResolveMocks{
				HomeFunc: fakeHomeResolver,
				AbsFunc:  fakeAbsResolver,
			}

			if filepath.Separator == '/' {
				actual := nfs.ResolvePath(entry.path, mocks)
				Expect(actual).To(Equal(entry.expect))
			} else {
				normalisedPath := strings.ReplaceAll(entry.path, "/", string(filepath.Separator))
				normalisedExpect := strings.ReplaceAll(entry.expect, "/", string(filepath.Separator))

				actual := nfs.ResolvePath(normalisedPath, mocks)
				Expect(actual).To(Equal(normalisedExpect))
			}
		},
		func(entry *RPEntry) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
		},

		Entry(nil, &RPEntry{
			given:  "path is a valid absolute path",
			should: "return path unmodified",
			path:   "/home/rabbitweed/foo",
			expect: "/home/rabbitweed/foo",
		}),
		Entry(nil, &RPEntry{
			given:  "path contains leading ~",
			should: "replace ~ with home path",
			path:   "~/foo",
			expect: "/home/rabbitweed/foo",
		}),
		Entry(nil, &RPEntry{
			given:  "path is relative to cwd",
			should: "replace ~ with home path",
			path:   "./foo",
			expect: "/home/rabbitweed/music/xpander/foo",
		}),
		Entry(nil, &RPEntry{
			given:  "path is relative to parent",
			should: "replace ~ with home path",
			path:   "../foo",
			expect: "/home/rabbitweed/music/foo",
		}),
		Entry(nil, &RPEntry{
			given:  "path is relative to grand parent",
			should: "replace ~ with home path",
			path:   "../../foo",
			expect: "/home/rabbitweed/foo",
		}),
	)

	When("No overrides provided", func() {
		Context("and: home", func() {
			It("🧪 should: not fail", func() {
				nfs.ResolvePath("~/")
			})
		})

		Context("and: abs cwd", func() {
			It("🧪 should: not fail", func() {
				nfs.ResolvePath("./")
			})
		})

		Context("and: abs parent", func() {
			It("🧪 should: not fail", func() {
				nfs.ResolvePath("../")
			})
		})

		Context("and: abs grand parent", func() {
			It("🧪 should: not fail", func() {
				nfs.ResolvePath("../..")
			})
		})
	})
})
