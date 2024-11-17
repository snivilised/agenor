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
	"github.com/snivilised/agenor/test/hydra"
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

		fS = hydra.Nuxx(verbose, lab.Static.RetroWave)
		from = lab.GetJSONPath()
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("walk",
		func(ctx SpecContext, entry *resumeTE) {
			invocations := strategyInvocations{}

			for _, strategy := range []enums.ResumeStrategy{
				enums.ResumeStrategyFastward,
				enums.ResumeStrategySpawn,
			} {
				recall := make(lab.Recall)
				profile, ok := profiles[entry.profile]

				if !ok {
					Fail(fmt.Sprintf("bad test, missing profile for '%v'", entry.profile))
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
						entry.active.resumeAt, entry.Subscription,
					)
					active.Tree = entry.Relative
					active.Depth = lo.Ternary(entry.active.depth == 0, 2, entry.active.depth)
					active.TraverseDescription.IsRelative = true
					active.ResumeDescription.IsRelative = false
					active.Subscription = entry.Subscription
					active.CurrentPath = entry.active.resumeAt
					active.Hibernation = entry.active.hibernateState
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
							return age.WithOnBegin(func(state *life.BeginState) {
								lab.Begin("🛡️")(state)
								//
								// don't enforce this yet, we need to disable notifications
								//
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
		func(entry *resumeTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
		},

		// === Listening (uni/folder/file) (pend/active)
		//
		// for the active cases, it doesn't really matter what the resumeAt is set
		// to, because the listener is already in the active listening state. But resumeAt
		// still has to be set because that is what would happen in the real world.
		//
		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: listen pending",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			active: activeTE{
				resumeAt:       ResumeAtTeenageColor,
				hibernateState: enums.HibernationPending,
			},
			clientListenAt: StartAtElectricYouth,
			profile:        "-> universal(pending): unfiltered",
		}),

		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: listen active",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			active: activeTE{
				resumeAt:       ResumeAtTeenageColor,
				hibernateState: enums.HibernationActive,
			},
			// For these scenarios (START_AT_CLIENT_ALREADY_ACTIVE), since
			// listening is already active, the value of resumeAt is irrelevant,
			// because the client was already listening in the previous session,
			// which is reflected by the state being active. So in essence, the client
			// listen value is a historical event, so the value defined here is a moot
			// point.
			//
			clientListenAt: StartAtClientAlreadyActive,
			profile:        "-> universal(active): unfiltered",
		}),

		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "folders: listen pending",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			active: activeTE{
				resumeAt:       ResumeAtTeenageColor,
				hibernateState: enums.HibernationPending,
			},
			clientListenAt: StartAtElectricYouth,
			profile:        "-> folders(pending): unfiltered",
		}),

		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "folders: listen active",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			active: activeTE{
				resumeAt:       ResumeAtTeenageColor,
				hibernateState: enums.HibernationActive,
			},
			clientListenAt: StartAtClientAlreadyActive,
			profile:        "-> folders(active): unfiltered",
		}),

		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "files: listen pending",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:       ResumeAtCanYouKissMeFirst,
				hibernateState: enums.HibernationPending,
			},
			clientListenAt: StartAtBeforeLife,
			profile:        "-> files(pending): unfiltered",
		}),

		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "files: listen active",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:       ResumeAtCanYouKissMeFirst,
				hibernateState: enums.HibernationActive,
			},
			clientListenAt: StartAtClientAlreadyActive,
			profile:        "-> files(active): unfiltered",
		}),

		// === Filtering (uni/folder/file)

		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: listen not active/deaf",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			active: activeTE{
				resumeAt:       ResumeAtTeenageColor,
				hibernateState: enums.HibernationRetired, // TODO(listen not active):check Retired is correct enum!!!
			},
			profile: "-> universal: filtered",
		}),

		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "folders: listen not active/deaf",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			active: activeTE{
				resumeAt:       ResumeAtTeenageColor,
				hibernateState: enums.HibernationRetired, // TODO:check Retired
			},
			profile: "-> folders: filtered",
		}),

		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "files: listen not active/deaf",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:       ResumeAtCanYouKissMeFirst,
				hibernateState: enums.HibernationRetired, // TODO:check Retired
			},
			profile: "-> files: filtered",
		}),

		// === Listening and filtering (uni/folder/file)

		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: listen pending and filtered",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			active: activeTE{
				resumeAt:       ResumeAtTeenageColor,
				hibernateState: enums.HibernationPending,
			},
			clientListenAt: StartAtElectricYouth,
			profile:        "-> universal: listen pending and filtered",
		}),

		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: listen active and filtered",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			active: activeTE{
				resumeAt:       ResumeAtTeenageColor,
				hibernateState: enums.HibernationActive,
			},
			clientListenAt: StartAtElectricYouth,
			profile:        "-> universal: filtered",
		}),

		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "folders: listen pending and filtered",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			active: activeTE{
				resumeAt:       ResumeAtTeenageColor,
				hibernateState: enums.HibernationPending,
			},
			clientListenAt: StartAtElectricYouth,
			profile:        "-> folders: listen pending and filtered",
		}),

		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "folders: listen active and filtered",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			active: activeTE{
				resumeAt:       ResumeAtTeenageColor,
				hibernateState: enums.HibernationActive,
			},
			clientListenAt: StartAtElectricYouth,
			profile:        "-> folders: filtered",
		}),

		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "files: listen pending and filtered",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:       ResumeAtCanYouKissMeFirst,
				hibernateState: enums.HibernationPending,
			},
			clientListenAt: StartAtBeforeLife,
			profile:        "-> files: listen pending and filtered",
		}),

		Entry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "files: listen active and filtered",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:       ResumeAtCanYouKissMeFirst,
				hibernateState: enums.HibernationActive,
			},
			clientListenAt: StartAtBeforeLife,
			profile:        "-> files: filtered",
		}),
	)
})
