package resume_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/life"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/test/hanno"
	"github.com/snivilised/agenor/tfs"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"
)

const (
	verbose = false
)

var noOp = func(string) {}

var _ = Describe("Resume", Ordered, func() {
	var (
		from string
		fS   *luna.MemFS
	)

	BeforeAll(func() {
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
				}
			},
		)).To(Succeed())

		fS = hanno.Nuxx(verbose, lab.Static.RetroWave)
		from = lab.GetJSONPath()
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("walk",
		func(ctx SpecContext, entry *lab.ResumeTE) {
			invocations := strategyInvocations{}

			for _, strategy := range []enums.ResumeStrategy{
				enums.ResumeStrategyFastward,
				enums.ResumeStrategySpawn,
			} {
				recall := make(lab.Recall)
				profile, ok := profiles[entry.Profile]

				if !ok {
					Fail(fmt.Sprintf("bad test, missing profile for '%v'", entry.Profile))
				}

				once := func(node *age.Node) error { //nolint:unparam // return nil error ok
					_, found := recall[node.Extension.Name]
					Expect(found).To(BeFalse())
					recall[node.Extension.Name] = len(node.Children)

					return nil
				}

				callback := func(servant age.Servant) error {
					node := servant.Node()
					depth := node.Extension.Depth
					GinkgoWriter.Printf(
						"---> ⏩ %v: (depth:%v) '%v'\n", strategy, depth, node.Path,
					)
					msg := fmt.Sprintf("%v, was invoked, but does not satisfy sample criteria",
						lab.Reason(node.Extension.Name),
					)
					Expect(entry.Prohibited).ToNot(ContainElement(node.Extension.Name), msg)

					if strategy == enums.ResumeStrategyFastward {
						segments := strings.Split(node.Path, "/")
						last := segments[len(segments)-1]

						if _, found := prohibited[last]; found {
							Fail(fmt.Sprintf("item: '%v' should have been fast forwarded over", node.Path))
						}
					}

					return once(node)
				}

				result, err := age.Walk().Configure(enclave.Loader(func(active *core.ActiveState) {
					GinkgoWriter.Printf("===> 🐚 restoring state: resume at=%v, subscription=%v\n",
						entry.Active.ResumeAt, entry.Subscription,
					)
					active.Tree = entry.Relative
					active.Depth = lo.Ternary(entry.Active.Depth == 0, 2, entry.Active.Depth)
					active.TraverseDescription.IsRelative = true
					active.ResumeDescription.IsRelative = false
					active.Subscription = entry.Subscription
					active.CurrentPath = entry.Active.ResumeAt
					active.Hibernation = entry.Active.HibernateState
				})).Extent(age.Resume(
					&pref.Relic{
						Head: pref.Head{
							Handler: callback,
							GetForest: func(_ string) *core.Forest {
								return &core.Forest{
									T: fS,
									R: tfs.New(),
								}
							},
						},
						From:     from,
						Strategy: strategy,
					},
					age.IfElseOptionF(strategy == enums.ResumeStrategyFastward,
						func() pref.Option {
							return age.WithOnBegin(func(_ *life.BeginState) {
								Fail("begin handler should not be invoked because begin notification muted")
							})
						},
						func() pref.Option {
							return age.WithOnBegin(lab.Begin("🛡️"))
						},
					),
					age.WithOnEnd(lab.End("🏁")),
				)).Navigate(ctx)

				if profile.mandatory != nil {
					for _, name := range profile.mandatory {
						_, found := recall[name]
						Expect(found).To(BeTrue(),
							fmt.Sprintf("mandatory item failure -> %v", lab.Reasons.Node(name)),
						)
					}
				}

				invocations[strategy] = strategyInvokeInfo{
					result.Metrics().Count(enums.MetricNoDirectoriesInvoked),
					result.Metrics().Count(enums.MetricNoFilesInvoked),
				}

				lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
					FS:            fS,
					Recording:     recall,
					Path:          entry.Relative,
					Result:        result,
					Err:           err,
					ExpectedErr:   entry.ExpectedErr,
					ByPassMetrics: true,
				})
			}
		},
		lab.FormatResumeTestDescription,

		// === Listening (uni/folder/file) (pend/active)
		//
		// for the active cases, it doesn't really matter what the resumeAt is set
		// to, because the listener is already in the active listening state. But resumeAt
		// still has to be set because that is what would happen in the real world.
		//
		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal: listen pending",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtTeenageColor,
				HibernateState: enums.HibernationPending,
			},
			ClientListenAt: StartAtElectricYouth,
			Profile:        "-> universal(pending): unfiltered",
		}),

		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal: listen active",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtTeenageColor,
				HibernateState: enums.HibernationActive,
			},
			// For these scenarios (START_AT_CLIENT_ALREADY_ACTIVE), since
			// listening is already active, the value of resumeAt is irrelevant,
			// because the client was already listening in the previous session,
			// which is reflected by the state being active. So in essence, the client
			// listen value is a historical event, so the value defined here is a moot
			// point.
			//
			ClientListenAt: StartAtClientAlreadyActive,
			Profile:        "-> universal(active): unfiltered",
		}),

		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "folders: listen pending",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtTeenageColor,
				HibernateState: enums.HibernationPending,
			},
			ClientListenAt: StartAtElectricYouth,
			Profile:        "-> folders(pending): unfiltered",
		}),

		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "folders: listen active",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtTeenageColor,
				HibernateState: enums.HibernationActive,
			},
			ClientListenAt: StartAtClientAlreadyActive,
			Profile:        "-> folders(active): unfiltered",
		}),

		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "files: listen pending",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtCanYouKissMeFirst,
				HibernateState: enums.HibernationPending,
			},
			ClientListenAt: StartAtBeforeLife,
			Profile:        "-> files(pending): unfiltered",
		}),

		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "files: listen active",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtCanYouKissMeFirst,
				HibernateState: enums.HibernationActive,
			},
			ClientListenAt: StartAtClientAlreadyActive,
			Profile:        "-> files(active): unfiltered",
		}),

		// === Filtering (uni/folder/file)

		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal: listen not active/deaf",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtTeenageColor,
				HibernateState: enums.HibernationRetired, // TODO(listen not active):check Retired is correct enum!!!
			},
			Profile: "-> universal: filtered",
		}),

		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "folders: listen not active/deaf",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtTeenageColor,
				HibernateState: enums.HibernationRetired, // TODO:check Retired
			},
			Profile: "-> folders: filtered",
		}),

		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "files: listen not active/deaf",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtCanYouKissMeFirst,
				HibernateState: enums.HibernationRetired, // TODO:check Retired
			},
			Profile: "-> files: filtered",
		}),

		// === Listening and filtering (uni/folder/file)

		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal: listen pending and filtered",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtTeenageColor,
				HibernateState: enums.HibernationPending,
			},
			ClientListenAt: StartAtElectricYouth,
			Profile:        "-> universal: listen pending and filtered",
		}),

		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal: listen active and filtered",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtTeenageColor,
				HibernateState: enums.HibernationActive,
			},
			ClientListenAt: StartAtElectricYouth,
			Profile:        "-> universal: filtered",
		}),

		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "folders: listen pending and filtered",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtTeenageColor,
				HibernateState: enums.HibernationPending,
			},
			ClientListenAt: StartAtElectricYouth,
			Profile:        "-> folders: listen pending and filtered",
		}),

		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "folders: listen active and filtered",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtTeenageColor,
				HibernateState: enums.HibernationActive,
			},
			ClientListenAt: StartAtElectricYouth,
			Profile:        "-> folders: filtered",
		}),

		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "files: listen pending and filtered",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtCanYouKissMeFirst,
				HibernateState: enums.HibernationPending,
			},
			ClientListenAt: StartAtBeforeLife,
			Profile:        "-> files: listen pending and filtered",
		}),

		Entry(nil, &lab.ResumeTE{
			DescribedTE: lab.DescribedTE{
				Given: "files: listen active and filtered",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			Active: lab.ActiveTE{
				ResumeAt:       ResumeAtCanYouKissMeFirst,
				HibernateState: enums.HibernationActive,
			},
			ClientListenAt: StartAtBeforeLife,
			Profile:        "-> files: filtered",
		}),
	)
})
