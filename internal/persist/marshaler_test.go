package persist_test

import (
	"errors"
	"fmt"
	"os"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/opts/json"
	"github.com/snivilised/traverse/internal/persist"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/lfs"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("Marshaler", Ordered, func() {
	const (
		foo  = "foo"
		flac = "*.flac"
		bar  = "*.bar"
	)

	var (
		FS                                       lfs.TraverseFS
		nodeFilterDef, polyNodeFilterDef         *core.FilterDef
		childFilterDef                           *core.ChildFilterDef
		samplingOptions                          *pref.SamplingOptions
		sampleFilterDef                          *core.SampleFilterDef
		jsonNodeFilterDef, jsonPolyNodeFilterDef json.FilterDef
		polyFilterDef                            *core.PolyFilterDef
	)

	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())
		nodeFilterDef = &core.FilterDef{
			Type:            enums.FilterTypeGlob,
			Description:     "items without .flac suffix",
			Pattern:         flac,
			Scope:           enums.ScopeAll,
			Negate:          true,
			IfNotApplicable: enums.TriStateBoolTrue,
		}

		childFilterDef = &core.ChildFilterDef{
			Type:        enums.FilterTypeGlob,
			Description: "items without .flac suffix",
			Pattern:     flac,
			Negate:      true,
		}

		samplingOptions = &pref.SamplingOptions{
			Type:      enums.SampleTypeFilter,
			InReverse: true,
			NoOf: pref.EntryQuantities{
				Files:   2,
				Folders: 3,
			},
		}

		sampleFilterDef = &core.SampleFilterDef{
			Type:        enums.FilterTypeGlob,
			Description: "items without .flac suffix",
			Pattern:     flac,
			Scope:       enums.ScopeAll,
			Negate:      true,
			Poly: &core.PolyFilterDef{
				File:   *nodeFilterDef,
				Folder: *nodeFilterDef,
			},
		}

		jsonNodeFilterDef = json.FilterDef{
			Type:            enums.FilterTypeGlob,
			Description:     "items without .flac suffix",
			Pattern:         flac,
			Scope:           enums.ScopeAll,
			Negate:          true,
			IfNotApplicable: enums.TriStateBoolTrue,
		}

		jsonPolyNodeFilterDef = json.FilterDef{
			Type:            enums.FilterTypePoly,
			Description:     "items without .flac suffix",
			Pattern:         flac,
			Scope:           enums.ScopeAll,
			Negate:          true,
			IfNotApplicable: enums.TriStateBoolTrue,
			Poly: &json.PolyFilterDef{
				File:   jsonNodeFilterDef,
				Folder: jsonNodeFilterDef,
			},
		}

		polyFilterDef = &core.PolyFilterDef{
			File:   *nodeFilterDef,
			Folder: *nodeFilterDef,
		}

		polyNodeFilterDef = &core.FilterDef{
			Type:            enums.FilterTypePoly,
			Description:     "items without .flac suffix",
			Pattern:         flac,
			Scope:           enums.ScopeAll,
			Negate:          true,
			IfNotApplicable: enums.TriStateBoolTrue,
			Poly:            polyFilterDef,
		}
	})

	BeforeEach(func() {
		FS = &lab.TestTraverseFS{
			MapFS: fstest.MapFS{
				home: &fstest.MapFile{
					Mode: os.ModeDir,
				},
			},
		}

		_ = FS.MkDirAll(to, permDir|os.ModeDir)
	})

	Context("map-fs", func() {
		DescribeTable("marshal",
			func(entry *marshalTE) {
				// success:
				o, _, err := opts.Get(
					pref.IfOptionF(entry.option != nil, func() pref.Option {
						return entry.option()
					}),
				)
				Expect(err).To(Succeed())

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
					writePath, permFile, FS,
				)

				Expect(err).To(Succeed())
				Expect(jo).NotTo(BeNil())

				// unequal error:
				if entry.tweak != nil {
					entry.tweak(jo)
				}
				equals, err := persist.Equals(o, jo)
				Expect(equals).To(BeFalse(), "should not compare equal")
				Expect(err).NotTo(Succeed())
			},
			func(entry *marshalTE) string {
				return fmt.Sprintf("given: %v, 🧪 should: marshal successfully", entry.given)
			},

			// 🍉 NavigationBehaviours:
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.SubPathBehaviour",
				},
				option: func() pref.Option {
					return pref.WithSubPathBehaviour(&pref.SubPathBehaviour{
						KeepTrailingSep: false,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Behaviours.SubPath.KeepTrailingSep = true
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.WithSortBehaviour",
				},
				option: func() pref.Option {
					return pref.WithSortBehaviour(&pref.SortBehaviour{
						IsCaseSensitive: true,
						SortFilesFirst:  true,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Behaviours.Sort.IsCaseSensitive = false
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.CascadeBehaviour.WithDepth",
				},
				option: func() pref.Option {
					return pref.WithDepth(4)
				},
				tweak: func(jo *json.Options) {
					jo.Behaviours.Cascade.Depth = 99
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.CascadeBehaviour.NoRecurse",
				},
				option: pref.WithNoRecurse,
				tweak: func(jo *json.Options) {
					jo.Behaviours.Cascade.NoRecurse = false
				},
			}),

			// 🍉 SamplingOptions:
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.SamplingOptions.InReverse",
				},
				option: func() pref.Option {
					return pref.WithSamplingOptions(&pref.SamplingOptions{
						Type:      enums.SampleTypeFilter,
						InReverse: true,
						NoOf: pref.EntryQuantities{
							Files:   3,
							Folders: 4,
						},
					})
				},
				tweak: func(jo *json.Options) {
					jo.Sampling.InReverse = false
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.SamplingOptions.SampleType",
				},
				option: func() pref.Option {
					return pref.WithSamplingOptions(&pref.SamplingOptions{
						Type:      enums.SampleTypeFilter,
						InReverse: true,
						NoOf: pref.EntryQuantities{
							Files:   3,
							Folders: 4,
						},
					})
				},
				tweak: func(jo *json.Options) {
					jo.Sampling.Type = enums.SampleTypeSlice
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.SamplingOptions.NoOf.Files",
				},
				option: func() pref.Option {
					return pref.WithSamplingOptions(&pref.SamplingOptions{
						Type:      enums.SampleTypeFilter,
						InReverse: true,
						NoOf: pref.EntryQuantities{
							Files:   3,
							Folders: 4,
						},
					})
				},
				tweak: func(jo *json.Options) {
					jo.Sampling.NoOf.Files = 99
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.SamplingOptions.NoOf.Folders",
				},
				option: func() pref.Option {
					return pref.WithSamplingOptions(&pref.SamplingOptions{
						Type:      enums.SampleTypeFilter,
						InReverse: true,
						NoOf: pref.EntryQuantities{
							Files:   3,
							Folders: 4,
						},
					})
				},
				tweak: func(jo *json.Options) {
					jo.Sampling.NoOf.Folders = 99
				},
			}),

			// 🍉 FilterOptions.Node
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - nil:pref.Options",
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Node = &json.FilterDef{
						Type:        enums.FilterTypeRegex,
						Description: foo,
						Pattern:     flac,
						Scope:       enums.ScopeFile,
					}
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - nil:json.Options",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: nodeFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Node = nil
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.Type",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: nodeFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Node.Type = enums.FilterTypeExtendedGlob
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.Description",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: nodeFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Node.Description = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.Pattern",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: nodeFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Node.Pattern = bar
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.Scope",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: nodeFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Node.Scope = enums.ScopeFile
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.Negate",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: nodeFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Node.Negate = false
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.IfNotApplicable",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: nodeFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Node.IfNotApplicable = enums.TriStateBoolFalse
				},
			}),

			// 🍉 FilterOptions.Node.Poly
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions.Node.Poly - nil:pref.Options",
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Node = &jsonPolyNodeFilterDef
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions.Node.Poly - nil:json.Options",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: polyNodeFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Node = nil
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.Poly.File",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: &core.FilterDef{
							Type: enums.FilterTypePoly,
							Poly: &core.PolyFilterDef{
								File:   *nodeFilterDef,
								Folder: *nodeFilterDef,
							},
						},
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Node.Poly.File.Description = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.Poly.Folder",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: &core.FilterDef{
							Type: enums.FilterTypePoly,
							Poly: &core.PolyFilterDef{
								File:   *nodeFilterDef,
								Folder: *nodeFilterDef,
							},
						},
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Node.Poly.Folder.Description = foo
				},
			}),

			// 🍉 FilterOptions.Child
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - nil:pref.Options",
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Child = &json.ChildFilterDef{
						Type:        enums.FilterTypeGlob,
						Description: "items without .flac suffix",
						Pattern:     flac,
						Negate:      true,
					}
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - nil:json.Options",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Child: childFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Child = nil
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Child.Type",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Child: childFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Child.Type = enums.FilterTypeExtendedGlob
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Child.Description",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Child: childFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Child.Description = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Child.Pattern",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Child: childFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Child.Pattern = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Child.Negate",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Child: childFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Child.Negate = false
				},
			}),

			// 🍉 FilterOptions.Sample
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - nil:pref.Options",
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Sample = &json.SampleFilterDef{
						Type:        enums.FilterTypeRegex,
						Description: foo,
						Pattern:     flac,
						Scope:       enums.ScopeFile,
					}
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - nil:json.Options",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sampleFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Sample = nil
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Sample.Type",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sampleFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Sample.Type = enums.FilterTypeExtendedGlob
				},
			}),
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Sample.Description",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sampleFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Sample.Description = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Sample.Pattern",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sampleFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Sample.Pattern = bar
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Sample.Scope",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sampleFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Sample.Scope = enums.ScopeFile
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Sample.Negate",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sampleFilterDef,
					})
				},
				tweak: func(jo *json.Options) {
					jo.Filter.Sample.Negate = false
				},
			}),

			// SamplingOptions:
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "Sampling - SampleOptions.Type",
				},
				option: func() pref.Option {
					return pref.WithSamplingOptions(samplingOptions)
				},
				tweak: func(jo *json.Options) {
					jo.Sampling.Type = enums.SampleTypeSlice
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "Sampling - SampleOptions.InReverse",
				},
				option: func() pref.Option {
					return pref.WithSamplingOptions(samplingOptions)
				},
				tweak: func(jo *json.Options) {
					jo.Sampling.InReverse = false
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "Sampling - SampleOptions.NoOf.Files",
				},
				option: func() pref.Option {
					return pref.WithSamplingOptions(samplingOptions)
				},
				tweak: func(jo *json.Options) {
					jo.Sampling.NoOf.Files = 99
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "Sampling - SampleOptions.NoOf.Folders",
				},
				option: func() pref.Option {
					return pref.WithSamplingOptions(samplingOptions)
				},
				tweak: func(jo *json.Options) {
					jo.Sampling.NoOf.Folders = 99
				},
			}),

			// 🍉 HibernateOptions:
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "HibernateOptions.Behaviour.InclusiveWake",
				},
				option: func() pref.Option {
					return pref.WithHibernationFilterWake(nodeFilterDef)
				},
				tweak: func(jo *json.Options) {
					jo.Hibernate.WakeAt.Description = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "HibernateOptions.Behaviour.InclusiveSleep",
				},
				option: func() pref.Option {
					return pref.WithHibernationFilterSleep(nodeFilterDef)
				},
				tweak: func(jo *json.Options) {
					jo.Hibernate.SleepAt.Description = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "HibernateOptions.Behaviour.InclusiveWake",
				},
				option: pref.WithHibernationBehaviourExclusiveWake,
				tweak: func(jo *json.Options) {
					jo.Hibernate.Behaviour.InclusiveWake = true
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "HibernateOptions.Behaviour.InclusiveSleep",
				},
				option: pref.WithHibernationBehaviourInclusiveSleep,
				tweak: func(jo *json.Options) {
					jo.Hibernate.Behaviour.InclusiveSleep = false
				},
			}),

			// 🍉 ConcurrencyOptions:
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "ConcurrencyOptions.NoW",
				},
				option: func() pref.Option {
					return pref.WithNoW(5)
				},
				tweak: func(jo *json.Options) {
					jo.Concurrency.NoW = 99
				},
			}),
		)

		Context("UnequalPtrError", func() {
			When("pref.Options is nil", func() {
				It("🧪 should: return UnequalPtrError", func() {
					equals, err := persist.Equals(nil, &json.Options{})
					Expect(equals).To(BeFalse(), "should not compare equal")
					Expect(err).NotTo(Succeed())
					Expect(errors.Is(err, locale.ErrUnEqualConversion)).To(BeTrue(),
						"error should be a locale.ErrUnEqualConversion",
					)
				})
			})

			When("json FilterDef is nil", func() {
				It("🧪 should: return UnequalPtrError", func() {
					o, _, _ := opts.Get()
					equals, err := persist.Equals(o, nil)
					Expect(equals).To(BeFalse(), "should not compare equal")
					Expect(err).NotTo(Succeed())
					Expect(errors.Is(err, locale.ErrUnEqualConversion)).To(BeTrue(),
						"error should be a locale.ErrUnEqualConversion",
					)
				})
			})
		})
	})
})
