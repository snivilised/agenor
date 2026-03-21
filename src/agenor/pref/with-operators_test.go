package pref_test

import (
	"io/fs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/li18ngo"
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
			Expect(agenor.WithCPU()).NotTo(BeNil())
		})
	})

	Context("WithDepth", func() {
		It("🧪 should: create option", func() {
			Expect(agenor.WithDepth(1)).NotTo(BeNil())
		})
	})

	Context("WithFaultHandler", func() {
		It("🧪 should: create option", func() {
			Expect(agenor.WithFaultHandler(&testFaultHandler{})).NotTo(BeNil())
		})
	})

	Context("WithFilter", func() {
		It("🧪 should: create option", func() {
			Expect(agenor.WithFilter(filterOptions)).NotTo(BeNil())
		})
	})

	Context("WithHibernationBehaviourExclusiveWake", func() {
		It("🧪 should: create option", func() {
			Expect(agenor.WithHibernationBehaviourExclusiveWake()).NotTo(BeNil())
		})
	})

	Context("WithHibernationBehaviourInclusiveSleep", func() {
		It("🧪 should: create option", func() {
			Expect(agenor.WithHibernationBehaviourInclusiveSleep()).NotTo(BeNil())
		})
	})

	Context("WithHibernationFilterSleep", func() {
		It("🧪 should: create option", func() {
			Expect(agenor.WithHibernationFilterSleep(nodeDef)).NotTo(BeNil())
		})
	})

	Context("WithHibernationFilterWake", func() {
		It("🧪 should: create option", func() {
			Expect(agenor.WithHibernationFilterWake(nodeDef)).NotTo(BeNil())
		})
	})

	Context("WithHibernationOptions", func() {
		It("🧪 should: create option", func() {
			Expect(agenor.WithHibernationOptions(&core.HibernateOptions{
				WakeAt: nodeDef,
			})).NotTo(BeNil())
		})
	})

	Context("WithHookSort", func() {
		It("🧪 should: create option", func() {
			option := agenor.WithHookSort(
				func([]fs.DirEntry, ...any) {},
			)
			Expect(option).NotTo(BeNil())
			_ = option(pref.DefaultOptions())
		})
	})

	Context("WithHookFileSubPath", func() {
		It("🧪 should: create option", func() {
			option := agenor.WithHookFileSubPath(func(*core.SubPathInfo) string {
				return ""
			})
			Expect(option).NotTo(BeNil())
			_ = option(pref.DefaultOptions())
		})
	})

	Context("WithHookDirectorySubPath", func() {
		It("🧪 should: create option", func() {
			option := agenor.WithHookDirectorySubPath(func(*core.SubPathInfo) string {
				return ""
			})
			Expect(option).NotTo(BeNil())
			_ = option(pref.DefaultOptions())
		})
	})

	Context("WithNavigationBehaviours", func() {
		It("🧪 should: create option", func() {
			option := agenor.WithNavigationBehaviours(
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
			)
			Expect(option).NotTo(BeNil())
			_ = option(pref.DefaultOptions())
		})
	})

	Context("WithPanicHandler", func() {
		It("🧪 should: create option", func() {
			option := agenor.WithPanicHandler(&testPanicHandler{})
			Expect(option).NotTo(BeNil())
			_ = option(pref.DefaultOptions())
		})
	})

	Context("WithNoRecurse", func() {
		It("🧪 should: create option", func() {
			Expect(agenor.WithNoRecurse()).NotTo(BeNil())
		})
	})

	Context("WithNoW", func() {
		It("🧪 should: create option", func() {
			Expect(agenor.WithNoW(3)).NotTo(BeNil())
		})
	})

	Context("WithSamplingOptions", func() {
		It("🧪 should: create option", func() {
			Expect(agenor.WithSamplingOptions(
				&pref.SamplingOptions{
					Type:      enums.SampleTypeFilter,
					InReverse: true,
					NoOf: pref.EntryQuantities{
						Files:       2,
						Directories: 3,
					},
				},
			)).NotTo(BeNil())
		})
	})

	Context("WithSkipHandler", func() {
		It("🧪 should: create option", func() {
			option := agenor.WithSkipHandler(&testSkipHandler{})
			Expect(option).NotTo(BeNil())
			_ = option(pref.DefaultOptions())
		})
	})

	Context("WithSortBehaviour", func() {
		It("🧪 should: create option", func() {
			option := agenor.WithSortBehaviour(&pref.SortBehaviour{
				IsCaseSensitive: true,
				SortFilesFirst:  true,
			})
			Expect(option).NotTo(BeNil())
			_ = option(pref.DefaultOptions())
		})
	})

	Context("WithSubPathBehaviour", func() {
		It("🧪 should: create option", func() {
			option := agenor.WithSubPathBehaviour(
				&pref.SubPathBehaviour{
					KeepTrailingSep: true,
				},
			)
			Expect(option).NotTo(BeNil())
			_ = option(pref.DefaultOptions())
		})
	})
})
