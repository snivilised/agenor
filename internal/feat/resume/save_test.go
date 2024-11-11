package resume_test

import (
	"context"
	"time"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/nefilim/test/luna"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/test/hydra"
	"github.com/snivilised/traverse/tfs"
)

const home = "home/prodigy"

type arrangeSave struct {
	jdir, name string
	rS         tv.TraversalFS
}

func (s arrangeSave) arrange() *saveAsserter {
	s.jdir = lab.GetJSONDir()
	calc := s.rS.Calc()
	mocks := &nef.ResolveMocks{
		HomeFunc: func() (string, error) {
			return calc.Join(s.jdir, "marshal", home), nil
		},
	}

	full, _ := mocks.HomeFunc()
	file, err := s.rS.Ensure(nef.PathAs{
		Name:    full,
		Default: s.name,
		Perm:    lab.Perms.Dir,
	})
	directory, _ := calc.Split(file)

	Expect(err).To(Succeed())
	Expect(luna.AsDirectory(directory)).To(luna.ExistInFS(s.rS))

	return &saveAsserter{
		directory: directory,
	}
}

type saveAsserter struct {
	expectedErr error
	directory   string
}

func (a *saveAsserter) assert(result core.TraverseResult, actual, expected error) {
	Expect(actual).To(MatchError(expected))
	Expect(result).NotTo(BeNil())

	if err, ok := actual.(*locale.TraversalSavedError); ok {
		Expect(err.Destination).NotTo(BeEmpty())
	}
}

var _ = Describe("Save", Ordered, func() {
	var (
		from, jdir string
		fS         *luna.MemFS
		rS         tv.TraversalFS
	)

	BeforeAll(func() {
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "traverse"},
				}
			},
		)).To(Succeed())
		from = lab.GetJSONPath()
		jdir = lab.GetJSONDir()
		core.Now = func() time.Time {
			t, _ := time.Parse(time.DateTime, "2024-11-14 15:04:05")
			return t
		}
	})

	BeforeEach(func() {
		services.Reset()
		fS = hydra.Nuxx(verbose, lab.Static.RetroWave)
		rS = tfs.New()
	})

	Context("Walk", func() {
		When("given: panic", func() {
			Context("prime", func() {
				It("🧪 should: save", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					save := (arrangeSave{
						name: "prime.walk.panic-save.json",
						rS:   rS,
						jdir: jdir,
					}).arrange()

					result, err := tv.Walk().Configure().Extent(tv.Prime(
						&pref.Using{
							Subscription: enums.SubscribeDirectories,
							Head: pref.Head{
								Handler: lab.PanicAt(lab.Static.TeenageColor),
								GetForest: func(_ string) *core.Forest {
									return &core.Forest{
										T: fS,
										R: rS,
									}
								},
							},
							Tree: lab.Static.RetroWave,
						},
						tv.WithOnBegin(lab.Begin("🛡️")),
						tv.WithOnEnd(lab.End("🏁")),
						pref.WithAdminPath(save.directory),
					)).Navigate(ctx)

					save.assert(result, err, locale.ErrCorePanicOccurred)
				})
			})

			Context("resume", func() {
				It("🧪 should: save", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					save := (arrangeSave{
						name: "resume.walk.panic-save.json",
						rS:   rS,
					}).arrange()

					result, err := tv.Walk().Configure(enclave.Loader(func(active *core.ActiveState) {
						active.Tree = lab.Static.RetroWave
						active.TraverseDescription = core.FsDescription{
							IsRelative: true,
						}
						active.CurrentPath = lab.Static.NorthernCouncil
						active.Subscription = enums.SubscribeUniversal
					})).Extent(tv.Resume(
						&pref.Relic{
							Head: pref.Head{
								Handler: lab.PanicAt(lab.Static.ElectricYouth),
								GetForest: func(_ string) *core.Forest {
									return &core.Forest{
										T: fS,
										R: rS,
									}
								},
							},
							From:     from,
							Strategy: enums.ResumeStrategyFastward,
						},
						tv.WithOnBegin(lab.Begin("🛡️")),
						tv.WithOnEnd(lab.End("🏁")),
						pref.WithAdminPath(save.directory),
					)).Navigate(ctx)

					save.assert(result, err, locale.ErrCorePanicOccurred)
				})
			})
		})
	})
	// TODO: repeat for concurrent sync (tv.Run).
})
