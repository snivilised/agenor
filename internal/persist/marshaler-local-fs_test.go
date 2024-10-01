package persist_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/persist"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/lfs"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("Marshaler", Ordered, func() {
	var testPath string

	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())

		testPath = lab.Repo("test")
		testFile := filepath.Join(testPath, to, tempFile)

		if _, err := os.Stat(testFile); err == nil {
			_ = os.Remove(testFile)
		}

		toPath := filepath.Join(testPath, to)
		if err := os.MkdirAll(toPath, permDir|os.ModeDir); err != nil {
			Fail(err.Error())
		}

		fromPath := filepath.Join(testPath, from)
		if err := os.MkdirAll(fromPath, permDir|os.ModeDir); err != nil {
			Fail(err.Error())
		}
	})

	Context("local-fs", func() {
		When("given pref.Options", func() {
			Context("marshall", func() {
				It("🧪 should: translate to json", func() {
					o, _, err := opts.Get(
						pref.WithDepth(4),
					)
					Expect(err).To(Succeed())

					writerFS := lfs.NewWriteFileFS(testPath, NoOverwrite)
					writePath := to + "/" + tempFile
					jo, err := persist.Marshal(&persist.MarshalState{
						O: o,
						Active: &types.ActiveState{
							Root:        to,
							Hibernation: enums.HibernationPending,
							NodePath:    "/root/a/b/c",
							Depth:       3,
						},
					},
						writePath, permFile, writerFS,
					)

					Expect(err).To(Succeed())
					Expect(jo).NotTo(BeNil())
				})
			})
		})
	})

	When("given json.Options", func() {
		Context("unmarshal", func() {
			XIt("🧪 should: translate from json", func() {
				/*
					o, _, _ := opts.Get()
					marshaller = persist.NewReader(o, &types.ActiveState{
						Root:        "some-root-path",
						Hibernation: enums.HibernationPending,
						NodePath:    "/root/a/b/c",
						Depth:       3,
					})
				*/
				readerFS := lfs.NewReadFileFS("/some-path")
				state, err := persist.Unmarshal(&types.RestoreState{
					Path:   "some-restore-path",
					Resume: enums.ResumeStrategySpawn,
				}, "/some-path", readerFS)
				_ = state

				Expect(err).To(Succeed())
			})
		})
	})
})
