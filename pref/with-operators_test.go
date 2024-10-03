package pref_test

import (
	"io/fs"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/pref"
)

var (
	nodeDef = &core.FilterDef{
		Type:        enums.FilterTypeGlob,
		Description: "items with '.flac' suffix",
		Pattern:     "*.flac",
		Scope:       enums.ScopeFile,
	}

	filterOptions = &pref.FilterOptions{
		Node: nodeDef,
	}
)

var _ = Describe("With Operators", Ordered, func() {
	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())
	})

	Context("WithCPU", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithCPU()).NotTo(BeNil())
		})
	})

	Context("WithDepth", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithDepth(1)).NotTo(BeNil())
		})
	})

	Context("WithFaultHandler", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithFaultHandler(&testFaultHandler{})).NotTo(BeNil())
		})
	})

	Context("WithFilter", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithFilter(filterOptions)).NotTo(BeNil())
		})
	})

	Context("WithHibernationBehaviourExclusiveWake", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithHibernationBehaviourExclusiveWake()).NotTo(BeNil())
		})
	})

	Context("WithHibernationBehaviourInclusiveSleep", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithHibernationBehaviourInclusiveSleep()).NotTo(BeNil())
		})
	})

	Context("WithHibernationFilterSleep", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithHibernationFilterSleep(nodeDef)).NotTo(BeNil())
		})
	})

	Context("WithHibernationFilterWake", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithHibernationFilterWake(nodeDef)).NotTo(BeNil())
		})
	})

	Context("WithHibernationOptions", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithHibernationOptions(&core.HibernateOptions{
				WakeAt: nodeDef,
			})).NotTo(BeNil())
		})
	})

	Context("WithHookSort", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithHookSort(
				func([]fs.DirEntry, ...any) {},
			)).NotTo(BeNil()) // !!no-test
		})
	})

	Context("WithHookFileSubPath", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithHookFileSubPath(func(*core.SubPathInfo) string {
				return ""
			})).NotTo(BeNil()) // !!no-test
		})
	})

	Context("WithHookFolderSubPath", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithHookFolderSubPath(func(*core.SubPathInfo) string {
				return ""
			})).NotTo(BeNil()) // !!no-test
		})
	})

	Context("WithNavigationBehaviours", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithNavigationBehaviours(
				&pref.NavigationBehaviours{
					SubPath: pref.SubPathBehaviour{
						KeepTrailingSep: true,
					},
					Sort: pref.SortBehaviour{
						IsCaseSensitive: true,
					},
					Cascade: pref.CascadeBehaviour{
						Depth: 2,
					},
				},
			)).NotTo(BeNil()) // !!no-test
		})
	})

	Context("WithPanicHandler", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithPanicHandler(&testPanicHandler{})).NotTo(BeNil()) // !!no-test
		})
	})

	Context("WithNoRecurse", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithNoRecurse()).NotTo(BeNil())
		})
	})

	Context("WithNoW", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithNoW(3)).NotTo(BeNil())
		})
	})

	Context("WithSamplingOptions", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithSamplingOptions(
				&pref.SamplingOptions{
					Type:      enums.SampleTypeFilter,
					InReverse: true,
					NoOf: pref.EntryQuantities{
						Files:   2,
						Folders: 3,
					},
				},
			)).NotTo(BeNil())
		})
	})

	Context("WithSkipHandler", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithSkipHandler(&testSkipHandler{})).NotTo(BeNil()) // !!no-test
		})
	})

	Context("WithSortBehaviour", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithSortBehaviour(&pref.SortBehaviour{
				IsCaseSensitive: true,
				SortFilesFirst:  true,
			})).NotTo(BeNil()) // !!no-test
		})
	})

	Context("WithSubPathBehaviour", func() {
		It("🧪 should: create option", func() {
			Expect(tv.WithSubPathBehaviour(
				&pref.SubPathBehaviour{
					KeepTrailingSep: true,
				},
			)).NotTo(BeNil()) // !!no-test
		})
	})
})
