package lfs_test

import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/lfs"
	"github.com/snivilised/traverse/locale"
)

// Note [clash/no-clash]: when an item is moved to the destination, clash
// refers to the scenario where the source already exists in the destination,
// whereas no-clash is the opposite.
//
// Eg, if moving src/foo.txt to dest/foo.txt, a clash scenario means foo.txt
// already exists in dest/ and no clash where it doesn't. The success of the
// operation depends on wether the overwrite flag has been specified in the
// filesystem.

var _ = Describe("op: move", Ordered, func() {
	var (
		root string
		fS   lfs.UniversalFS
	)

	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())

		root = lab.Repo("test")
	})

	DescribeTable("fs: UniversalFS",
		func(entry fsTE[lfs.UniversalFS]) {
			for _, overwrite := range []bool{false, true} {
				scratch(root)

				fS = lfs.NewUniversalFS(lfs.At{
					Root:      root,
					Overwrite: overwrite,
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

		// 🔆
		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "[from] file exists, [to] directory exists, [no-clash]",
			should:  "succeed",
			note:    "filename not included in the destination path (from/file.txt => to)",
			op:      "Move",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Move.From.File,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				Expect(fS.Move(entry.from, lab.Static.FS.Move.Destination)).To(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
				Expect(AsFile(lab.Static.FS.Move.To.File)).To(ExistInFS(fS))
			},
		}),

		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "[from] file exists, [to] directory exists, [clash]",
			should:  "succeed, only if override",
			note:    "filename not included in the destination path (from/file.txt => to)",
			op:      "Move",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Move.From.File,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				Expect(require(root,
					entry.require,
					entry.from,
				)).To(Succeed())
				Expect(require(root,
					lab.Static.FS.Move.Destination,
					lab.Static.FS.Move.To.File,
				)).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				if entry.overwrite {
					Expect(fS.Move(entry.from, lab.Static.FS.Move.Destination)).To(Succeed(),
						fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
					)
					Expect(AsFile(lab.Static.FS.Move.To.File)).To(ExistInFS(fS))
					return
				}
				Expect(fS.Move(entry.from, lab.Static.FS.Move.Destination)).NotTo(Succeed())
			},
		}),

		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "[from] file exists, [to] directory exists, [no-clash]",
			should:  "succeed",
			note:    "filename IS included in the destination path (from/file.txt => to/file.txt)",
			op:      "Move",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Move.From.File,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				destination := lab.Static.FS.Move.To.File
				Expect(fS.Move(entry.from, destination)).To(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
				Expect(AsFile(destination)).To(ExistInFS(fS))
			},
		}),

		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "[from] file exists, [to] directory exists, [clash]",
			should:  "succeed, only if override",
			note:    "filename IS included in the destination path (from/file.txt => to/file.txt)",
			op:      "Move",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Move.From.File,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())

				if entry.overwrite {
					Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
					return
				}
				Expect(require(root,
					lab.Static.FS.Move.Destination,
					lab.Static.FS.Move.To.File,
				)).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				destination := lab.Static.FS.Move.To.File

				if entry.overwrite {
					Expect(fS.Move(entry.from, destination)).To(Succeed(),
						fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
					)
					Expect(AsFile(destination)).To(ExistInFS(fS))
					return
				}
				Expect(fS.Move(entry.from, destination)).NotTo(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
			},
		}),

		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "[from] directory exists, [to] directory exists, [no clash]",
			should:  "succeed",
			note:    "directory not included in the destination path (from/dir => to)",
			op:      "Move",
			require: lab.Static.FS.Move.From.Directory,
			from:    lab.Static.FS.Move.From.Directory,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				Expect(fS.Move(entry.from, lab.Static.FS.Move.Destination)).To(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
				Expect(AsDirectory(lab.Static.FS.Move.To.Directory)).To(ExistInFS(fS))
			},
		}),

		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "[from] directory exists, [to] directory exists, [clash]",
			should:  "fail",
			note:    "directory not included in the destination path (from/dir => to)",
			op:      "Move",
			require: lab.Static.FS.Move.From.Directory,
			from:    lab.Static.FS.Move.From.Directory,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.To.Directory)).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				if entry.overwrite {
					Expect(fS.Move(entry.from, lab.Static.FS.Move.Destination)).NotTo(Succeed(),
						fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
					)
					return
				}

				Expect(fS.Move(entry.from, lab.Static.FS.Move.Destination)).NotTo(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
			},
		}),

		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "[from] directory exists, [to] directory exists, [no clash]",
			should:  "succeed",
			note:    "directory IS included in the destination path (from/dir => to/dir)",
			op:      "Move",
			require: lab.Static.FS.Move.From.Directory,
			from:    lab.Static.FS.Move.From.Directory,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				destination := lab.Static.FS.Move.To.Directory
				Expect(fS.Move(entry.from, destination)).To(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
				Expect(AsDirectory(destination)).To(ExistInFS(fS))
			},
		}),

		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "[from] directory exists, [to] directory exists, [clash]",
			should:  "fail",
			note:    "directory IS included in the destination path (from/dir => to/dir)",
			op:      "Move",
			require: lab.Static.FS.Move.From.Directory,
			from:    lab.Static.FS.Move.From.Directory,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.To.Directory)).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				destination := lab.Static.FS.Move.To.Directory
				Expect(fS.Move(entry.from, destination)).NotTo(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
			},
		}),

		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "[from] file missing",
			should:  "fail",
			op:      "Move",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.Foo,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				err := fS.Move(entry.from, entry.to)
				Expect(err).NotTo(Succeed(), fmt.Sprintf("OVERWRITE: %v", entry.overwrite))
				Expect(locale.IsBinaryFsOpError(err)).To(BeTrue())
			},
		}),

		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "[from] file exists, [to] directory missing",
			should:  "fail",
			op:      "Move",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Move.From.File,
			to:      lab.Static.FS.Move.Destination,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				destination := filepath.Join(entry.to, lab.Static.Foo)
				Expect(fS.Move(entry.from, destination)).NotTo(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
			},
		}),

		// The following tests are a duplicate of those defined for the rename
		// operation 💠, where the target item is being renamed into the same directory,
		// but these should be rejected with an error, because this amounts to a
		// rename which is not the intended use of Move. Instead a custom error
		// is returned that denotes same directory so that the client can detect this
		// and invoke rename, with the same parameters. The exception to this is
		// if the from/to names refer to the same file, in which case the Move
		// is ignored as a no op.
		//
		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "[from] file exists, [to] name does not exist, [no-clash]",
			should:  "fail, same directory move, use rename instead",
			op:      "Move",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Rename.From.File,
			to:      lab.Static.FS.Rename.To.File,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				IsSameDirMoveRejectionError(fS.Move(entry.from, entry.to), entry.should)
			},
		}),

		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "[from] file exists, [to] equal to [from], [clash]",
			should:  "succeed, ignored",
			op:      "Move",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Rename.From.File,
			to:      lab.Static.FS.Rename.From.File,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				Expect(fS.Move(entry.from, entry.to)).To(Succeed())
			},
		}),

		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "[from] directory exists, [to] name does not exist, [no-clash]",
			should:  "fail, same directory move, use rename instead",
			op:      "Move",
			require: lab.Static.FS.Rename.From.Directory,
			from:    lab.Static.FS.Rename.From.Directory,
			to:      lab.Static.FS.Rename.To.Directory,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				IsSameDirMoveRejectionError(fS.Move(entry.from, entry.to), entry.should)
			},
		}),

		Entry(nil, fsTE[lfs.UniversalFS]{
			given:   "[from] directory exists, [to] equal to [from], [clash]",
			should:  "fail, directory names can't be same",
			op:      "Move",
			require: lab.Static.FS.Rename.From.Directory,
			from:    lab.Static.FS.Rename.From.Directory,
			to:      lab.Static.FS.Rename.From.Directory,
			arrange: func(entry fsTE[lfs.UniversalFS], _ lfs.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
			},
			action: func(entry fsTE[lfs.UniversalFS], fS lfs.UniversalFS) {
				IsSameDirMoveRejectionError(fS.Move(entry.from, entry.to), entry.should)
			},
		}),
	)
})
