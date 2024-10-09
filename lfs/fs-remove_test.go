package lfs_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/lfs"
)

var _ = Describe("op: remove", Ordered, func() {
	var (
		root string
		fS   lfs.UniversalFS
	)

	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())

		root = lab.Repo("test")
	})

	BeforeEach(func() {
		scratchPath := filepath.Join(root, lab.Static.FS.Scratch)

		if _, err := os.Stat(scratchPath); err == nil {
			Expect(os.RemoveAll(scratchPath)).To(Succeed(),
				fmt.Sprintf("failed to delete existing directory %q", scratchPath),
			)
		}
	})

	DescribeTable("removal",
		func(entry fsTE[lfs.UniversalFS]) {
			for _, overwrite := range []bool{false, true} {
				fS = lfs.NewUniversalFS(lfs.At{
					Root:      root,
					Overwrite: entry.overwrite,
				})
				entry.overwrite = overwrite

				if entry.arrange != nil {
					entry.arrange(entry, fS)
				}
				entry.action(entry, fS)
			}
		},
		func(entry fsTE[lfs.UniversalFS]) string {
			return fmt.Sprintf("🧪 ===> given: target is '%v', %v should: '%v'",
				entry.given, entry.op, entry.should,
			)
		},
		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "file and exists",
			should:  "succeed",
			op:      "Remove",
			require: lab.Static.FS.Scratch,
			target:  lab.Static.FS.Remove.File,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				err := require(root, entry.require, entry.target)
				Expect(err).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				if entry.overwrite {
					// tbd
					return
				}
				Expect(fS.Remove(entry.target)).To(Succeed())
			},
		}),
		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "path does not exist",
			should:  "fail",
			op:      "Remove",
			require: lab.Static.FS.Scratch,
			target:  lab.Static.Foo,
			action: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				if entry.overwrite {
					// tbd
					return
				}
				Expect(fS.Remove(entry.target)).To(MatchError(os.ErrNotExist))
			},
		}),
		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "directory exists and not empty",
			should:  "fail",
			op:      "Remove",
			require: lab.Static.FS.Scratch,
			target:  lab.Static.FS.Scratch,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				err := require(root, entry.require, lab.Static.FS.Remove.File)
				Expect(err).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				if entry.overwrite {
					// tbd
					return
				}
				Expect(errors.Unwrap(fS.Remove(entry.target))).To(
					MatchError("directory not empty"),
				)
			},
		}),
		//
		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "path does not exist",
			should:  "succeed",
			op:      "RemoveAll",
			require: lab.Static.FS.Scratch,
			target:  lab.Static.Foo,
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				if entry.overwrite {
					// tbd
					return
				}
				Expect(fS.RemoveAll(entry.target)).To(Succeed())
			},
		}),
		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "directory exists and not empty",
			should:  "succeed",
			op:      "RemoveAll",
			require: lab.Static.FS.Scratch,
			target:  lab.Static.FS.Scratch,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				err := require(root, entry.require, lab.Static.FS.Remove.File)
				Expect(err).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				if entry.overwrite {
					// tbd
					return
				}
				Expect(fS.RemoveAll(entry.target)).To(Succeed())
			},
		}),
	)
})
