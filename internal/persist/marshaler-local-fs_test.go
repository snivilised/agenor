package persist_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	nef "github.com/snivilised/nefilim"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/persist"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/test/hydra"
)

var _ = Describe("Marshaler", Ordered, func() {
	var (
		testPath string
		reporter *core.Supervisor
		trig     *lab.Trigger
	)

	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())

		testPath = hydra.Repo("test")
		testFile := filepath.Join(testPath, destination, tempFile)

		if _, err := os.Stat(testFile); err == nil {
			_ = os.Remove(testFile)
		}

		toPath := filepath.Join(testPath, destination)
		if err := os.MkdirAll(toPath, lab.Perms.Dir|os.ModeDir); err != nil {
			Fail(err.Error())
		}

		fromPath := filepath.Join(testPath, source)
		if err := os.MkdirAll(fromPath, lab.Perms.Dir|os.ModeDir); err != nil {
			Fail(err.Error())
		}

		reporter = core.NewSupervisor()
		trig = &lab.Trigger{
			Metrics: reporter.Many(
				enums.MetricNoFilesInvoked,
				enums.MetricNoFilesFilteredOut,
				enums.MetricNoDirectoriesInvoked,
				enums.MetricNoDirectoriesFilteredOut,
			),
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

					trig.Times(
						enums.MetricNoFilesInvoked, 1).Times(
						enums.MetricNoFilesFilteredOut, 2).Times(
						enums.MetricNoDirectoriesInvoked, 3).Times(
						enums.MetricNoDirectoriesFilteredOut, 4,
					)

					writerFS := nef.NewWriteFileFS(tv.Rel{
						Root:      testPath,
						Overwrite: NoOverwrite,
					})
					writePath := destination + "/" + tempFile
					jo, err := persist.Marshal(&persist.MarshalRequest{
						O: o,
						Active: &core.ActiveState{
							Tree:        destination,
							Hibernation: enums.HibernationPending,
							CurrentPath: "/top/a/b/c",
							Depth:       3,
							Metrics:     trig.Metrics,
						},
						Path: writePath,
						Perm: lab.Perms.File,
						FS:   writerFS,
					})

					Expect(err).To(Succeed())
					Expect(jo).NotTo(BeNil())
				})
			})
		})
	})

	When("given json.Options", func() {
		Context("unmarshal", func() {
			It("🧪 should: translate from json", func() {
				fS := nef.NewReadFileFS(nef.Rel{
					Root: testPath,
				})
				result, err := persist.Unmarshal(&persist.UnmarshalRequest{
					Restore: &enclave.RestoreState{
						Path:   lab.Static.JSONSubPath,
						FS:     fS,
						Resume: enums.ResumeStrategySpawn,
					},
				})

				Expect(err).To(Succeed())
				Expect(result).NotTo(BeNil())
			})
		})
	})
})
