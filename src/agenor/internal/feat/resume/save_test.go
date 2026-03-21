package resume_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/internal/enclave"
	lab "github.com/snivilised/jaywalk/src/agenor/internal/laboratory"
	"github.com/snivilised/jaywalk/src/internal/services"
	"github.com/snivilised/jaywalk/locale"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/agenor/test/hanno"
	"github.com/snivilised/jaywalk/src/agenor/tfs"
	"github.com/snivilised/li18ngo"
	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/nefilim/test/luna"
)

const home = "home/prodigy"

type arrangeSave struct {
	jdir, name string
	rS         agenor.TraversalFS
}

func (s arrangeSave) arrange() *saveAsserter {
	s.jdir = lab.GetJSONDir()
	calc := s.rS.Calc()
	mocks := &nef.ResolveMocks{
		HomeFunc: func() (string, error) {
			return calc.Join(s.jdir, "marshal", home), nil
		},
	}

	directory, _ := mocks.HomeFunc()
	_, err := s.rS.Ensure(nef.PathAs{
		Name:    directory,
		Default: s.name,
		Perm:    lab.Perms.Dir,
	})

	Expect(err).To(Succeed())
	Expect(luna.AsDirectory(directory)).To(luna.ExistInFS(s.rS))

	return &saveAsserter{
		directory: directory,
	}
}

type saveAsserter struct {
	directory string
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
		rS         agenor.TraversalFS
	)

	BeforeAll(func() {
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
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

		fS = hanno.Nuxx(verbose, lab.Static.RetroWave)
		rS = tfs.New()
	})

	Context("Walk", func() {
		When("given: panic", func() {
			Context("prime", func() {
				It("🧪 should: save", func(specCtx SpecContext) {
					lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
						save := (arrangeSave{
							name: "prime.walk.panic-save.json",
							rS:   rS,
							jdir: jdir,
						}).arrange()

						result, err := agenor.Walk().Configure().Extent(agenor.Prime(
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
							agenor.WithOnBegin(lab.Begin("🛡️")),
							agenor.WithOnEnd(lab.End("🏁")),
							pref.WithAdminPath(save.directory),
						)).Navigate(ctx)

						save.assert(result, err, locale.ErrCorePanicOccurred)
					})
				})
			})

			Context("resume", func() {
				It("🧪 should: save", func(specCtx SpecContext) {
					lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
						save := (arrangeSave{
							name: "resume.walk.panic-save.json",
							rS:   rS,
						}).arrange()

						result, err := agenor.Walk().Configure(enclave.Loader(func(active *core.ActiveState) {
							active.Tree = lab.Static.RetroWave
							active.TraverseDescription = core.FsDescription{
								IsRelative: true,
							}
							active.CurrentPath = lab.Static.NorthernCouncil
							active.Subscription = enums.SubscribeUniversal
						})).Extent(agenor.Resume(
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
							agenor.WithOnBegin(lab.Begin("🛡️")),
							agenor.WithOnEnd(lab.End("🏁")),
							pref.WithAdminPath(save.directory),
						)).Navigate(ctx)

						save.assert(result, err, locale.ErrCorePanicOccurred)
					})
				})
			})
		})
	})
	// TODO: repeat for concurrent sync (agenor.Run).
})
