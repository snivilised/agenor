package resume_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/nefilim/luna"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/life"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/test/hydra"
)

const (
	verbose = false
)

var noOp = func(string) {}

var _ = Describe("Resume", Ordered, func() { // formerly resume-strategy_test
	var (
		jsonPath string
		tree     string
		fS       *luna.MemFS
	)

	BeforeAll(func() {
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "traverse"},
				}
			},
		)).To(Succeed())

		fS = hydra.Nuxx(verbose, lab.Static.RetroWave)
		jsonPath = lab.GetJSONPath()
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("walk",
		func(ctx SpecContext, entry *resumeTE) {
			invocations := strategyInvocations{}

			for _, strategy := range strategies {
				recording := make(lab.RecordingMap)
				profile, ok := profiles[entry.profile]

				if !ok {
					Fail(fmt.Sprintf("bad test, missing profile for '%v'", entry.profile))
				}

				restorer := func(o *pref.Options, ts *pref.TraversalState) error {
					// synthetic assignments: The client should not perform these
					// types of assignments. Only being done here for testing purposes
					// to avoid the need to create many restore files
					// (eg resume-state.json) for different test cases.
					//

					// this is akin to tampering for testing purpose; needs to be re-thought
					//
					ts.Tree = tree
					ts.CurrentPath = entry.Relative
					ts.Hibernation = entry.active.listenState

					if profile.filtered {
						// the json resume state contains a filter definition, so the
						// filtered flag determines if this filter should be applied
						// to the test case.
						//
						// However, test-restore.DEFAULT, reflects the default options
						// which does not have a filter defined. So we should reverse
						// this and define a filter if profile.filtered.
						// TODO: define a filter.
						noOp("waiting for a filter to be defined")
					}
					//
					// end of synthetic assignments

					if strategy == enums.ResumeStrategyFastward {
						o.Events.Begin.On(func(_ *life.BeginState) {
							Fail("begin handler should not be invoked because begin notification muted")
						})
					}
					GinkgoWriter.Printf("===> 🐚 restoring ...\n")

					return nil
				}

				once := func(node *tv.Node) error { //nolint:unparam // return nil error ok
					_, found := recording[node.Extension.Name]
					Expect(found).To(BeFalse())
					recording[node.Extension.Name] = len(node.Children)

					return nil
				}

				callback := func(servant tv.Servant) error {
					node := servant.Node()
					depth := node.Extension.Depth
					GinkgoWriter.Printf(
						"---> ⏩ %v: (depth:%v) '%v'\n", themes[strategy].label, depth, node.Path,
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

				// Do we have a WithRestore option, that also accepts
				// the active state?
				// Was contains the info that was in nav.Resumption
				//
				// MarshalRequest??
				//
				// nav.RunnerInfo => should be built into Run

				// the resume process starts off at the plugin
				//
				result, err := tv.Walk().Configure().Extent(tv.Resume(
					&tv.Was{
						Using: pref.Using{
							Tree:         entry.Relative,
							Subscription: entry.Subscription,
							Handler:      callback,
							GetForest: func(_ string) *core.Forest {
								return &core.Forest{
									T: fS,
									R: nef.NewTraverseABS(),
								}
							},
						},
						From:     jsonPath,
						Strategy: strategy,
						Restorer: restorer,
					},
				)).Navigate(ctx)

				if profile.mandatory != nil {
					for _, name := range profile.mandatory {
						_, found := recording[name]
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
					FS:          fS,
					Recording:   recording,
					Path:        entry.Relative,
					Result:      result,
					Err:         err,
					ExpectedErr: entry.ExpectedErr,
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
		XEntry(nil, &resumeTE{ // UNDER CONSTRUCTION !!!
			NaviTE: lab.NaviTE{
				Given:        "universal: listen pending",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			active: activeTE{
				resumeAt:    ResumeAtTeenageColor,
				listenState: enums.HibernationPending,
			},
			clientListenAt: StartAtElectricYouth,
			profile:        "-> universal(pending): unfiltered",
		}),

		XEntry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: listen active",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			active: activeTE{
				resumeAt:    ResumeAtTeenageColor,
				listenState: enums.HibernationActive,
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

		XEntry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "folders: listen pending",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			active: activeTE{
				resumeAt:    ResumeAtTeenageColor,
				listenState: enums.HibernationPending,
			},
			clientListenAt: StartAtElectricYouth,
			profile:        "-> folders(pending): unfiltered",
		}),

		XEntry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "folders: listen active",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			active: activeTE{
				resumeAt:    ResumeAtTeenageColor,
				listenState: enums.HibernationActive,
			},
			clientListenAt: StartAtClientAlreadyActive,
			profile:        "-> folders(active): unfiltered",
		}),

		XEntry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "files: listen pending",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:    ResumeAtCanYouKissMeFirst,
				listenState: enums.HibernationPending,
			},
			clientListenAt: StartAtBeforeLife,
			profile:        "-> files(pending): unfiltered",
		}),

		XEntry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "files: listen active",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:    ResumeAtCanYouKissMeFirst,
				listenState: enums.HibernationActive,
			},
			clientListenAt: StartAtClientAlreadyActive,
			profile:        "-> files(active): unfiltered",
		}),

		// === Filtering (uni/folder/file)

		XEntry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: listen not active/deaf",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			active: activeTE{
				resumeAt:    ResumeAtTeenageColor,
				listenState: enums.HibernationRetired, // TODO(listen not active):check Retired is correct enum!!!
			},
			profile: "-> universal: filtered",
		}),

		XEntry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "folders: listen not active/deaf",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			active: activeTE{
				resumeAt:    ResumeAtTeenageColor,
				listenState: enums.HibernationRetired, // TODO:check Retired
			},
			profile: "-> folders: filtered",
		}),

		XEntry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "files: listen not active/deaf",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:    ResumeAtCanYouKissMeFirst,
				listenState: enums.HibernationRetired, // TODO:check Retired
			},
			profile: "-> files: filtered",
		}),

		// === Listening and filtering (uni/folder/file)

		XEntry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: listen pending and filtered",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			active: activeTE{
				resumeAt:    ResumeAtTeenageColor,
				listenState: enums.HibernationPending,
			},
			clientListenAt: StartAtElectricYouth,
			profile:        "-> universal: listen pending and filtered",
		}),

		XEntry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: listen active and filtered",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
			},
			active: activeTE{
				resumeAt:    ResumeAtTeenageColor,
				listenState: enums.HibernationActive,
			},
			clientListenAt: StartAtElectricYouth,
			profile:        "-> universal: filtered",
		}),

		XEntry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "folders: listen pending and filtered",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			active: activeTE{
				resumeAt:    ResumeAtTeenageColor,
				listenState: enums.HibernationPending,
			},
			clientListenAt: StartAtElectricYouth,
			profile:        "-> folders: listen pending and filtered",
		}),

		XEntry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "folders: listen active and filtered",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
			},
			active: activeTE{
				resumeAt:    ResumeAtTeenageColor,
				listenState: enums.HibernationActive,
			},
			clientListenAt: StartAtElectricYouth,
			profile:        "-> folders: filtered",
		}),

		XEntry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "files: listen pending and filtered",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:    ResumeAtCanYouKissMeFirst,
				listenState: enums.HibernationPending,
			},
			clientListenAt: StartAtBeforeLife,
			profile:        "-> files: listen pending and filtered",
		}),

		XEntry(nil, &resumeTE{
			NaviTE: lab.NaviTE{
				Given:        "files: listen active and filtered",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:    ResumeAtCanYouKissMeFirst,
				listenState: enums.HibernationActive,
			},
			clientListenAt: StartAtBeforeLife,
			profile:        "-> files: filtered",
		}),
	)
})
