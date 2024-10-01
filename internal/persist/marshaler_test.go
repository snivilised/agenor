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
	var (
		FS             lfs.TraverseFS
		nodeFilterDef  *core.FilterDef
		childFilterDef *core.ChildFilterDef
	)

	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())
		nodeFilterDef = &core.FilterDef{
			Type:            enums.FilterTypeGlob,
			Description:     "items without .flac suffix",
			Pattern:         "*.flac",
			Scope:           enums.ScopeAll,
			Negate:          true,
			IfNotApplicable: enums.TriStateBoolTrue,
		}

		childFilterDef = &core.ChildFilterDef{
			Type:        enums.FilterTypeGlob,
			Description: "items without .flac suffix",
			Pattern:     "*.flac",
			Negate:      true,
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
		Context("marshal", func() {
			DescribeTable("success",
				func(entry *errorTE) {
					o, _, err := opts.Get(
						entry.option(),
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

					equals, err := persist.Equals(o, &json.Options{})
					Expect(equals).To(BeFalse(), "should not compare equal")
					Expect(err).NotTo(Succeed())
				},
				func(entry *errorTE) string {
					return fmt.Sprintf("given: %v, 🧪 should: marshal successfully", entry.given)
				},

				// NavigationBehaviours:
				Entry(nil, &errorTE{
					marshalTE: marshalTE{
						given: "NavigationBehaviours.SubPathBehaviour",
						option: func() pref.Option {
							return pref.WithSubPathBehaviour(&pref.SubPathBehaviour{
								KeepTrailingSep: false,
							})
						},
					},
				}),

				Entry(nil, &errorTE{
					marshalTE: marshalTE{
						given: "NavigationBehaviours.WithSortBehaviour",
						option: func() pref.Option {
							return pref.WithSortBehaviour(&pref.SortBehaviour{
								IsCaseSensitive: true,
								SortFilesFirst:  true,
							})
						},
					},
				}),

				Entry(nil, &errorTE{
					marshalTE: marshalTE{
						given: "NavigationBehaviours.CascadeBehaviour.WithDepth",
						option: func() pref.Option {
							return pref.WithDepth(4)
						},
					},
				}),

				Entry(nil, &errorTE{
					marshalTE: marshalTE{
						given:  "NavigationBehaviours.CascadeBehaviour.NoRecurse",
						option: pref.WithNoRecurse,
					},
				}),

				// SamplingOptions:
				Entry(nil, &errorTE{
					marshalTE: marshalTE{
						given: "NavigationBehaviours.SamplingOptions",
						option: func() pref.Option {
							return pref.WithSamplingOptions(&pref.SamplingOptions{
								SampleType:      enums.SampleTypeFilter,
								SampleInReverse: true,
								NoOf: pref.EntryQuantities{
									Files:   3,
									Folders: 4,
								},
							})
						},
					},
				}),

				// FilterOptions:
				Entry(nil, &errorTE{
					marshalTE: marshalTE{
						given: "FilterOptions - Node",
						option: func() pref.Option {
							return pref.WithFilter(&pref.FilterOptions{
								Node: nodeFilterDef,
							})
						},
					},
				}),

				Entry(nil, &errorTE{
					marshalTE: marshalTE{
						given: "FilterOptions - Child",
						option: func() pref.Option {
							return pref.WithFilter(&pref.FilterOptions{
								Child: childFilterDef,
							})
						},
					},
				}),

				Entry(nil, &errorTE{
					marshalTE: marshalTE{
						given: "FilterOptions - Sample",
						option: func() pref.Option {
							return pref.WithFilter(&pref.FilterOptions{
								Sample: &core.SampleFilterDef{
									Type:        enums.FilterTypeGlob,
									Description: "items without .flac suffix",
									Pattern:     "*.flac",
									Scope:       enums.ScopeAll,
									Negate:      true,
									Poly: &core.PolyFilterDef{
										File:   *nodeFilterDef,
										Folder: *nodeFilterDef,
									},
								},
							})
						},
					},
				}),

				// HibernateOptions:
				Entry(nil, &errorTE{
					marshalTE: marshalTE{
						given: "HibernateOptions.Behaviour.InclusiveWake",
						option: func() pref.Option {
							return pref.WithHibernationFilterWake(nodeFilterDef)
						},
					},
				}),

				Entry(nil, &errorTE{
					marshalTE: marshalTE{
						given: "HibernateOptions.Behaviour.InclusiveSleep",
						option: func() pref.Option {
							return pref.WithHibernationFilterSleep(nodeFilterDef)
						},
					},
				}),

				Entry(nil, &errorTE{
					marshalTE: marshalTE{
						given:  "HibernateOptions.Behaviour.InclusiveWake",
						option: pref.WithHibernationBehaviourExclusiveWake,
					},
				}),

				Entry(nil, &errorTE{
					marshalTE: marshalTE{
						given:  "HibernateOptions.Behaviour.InclusiveSleep",
						option: pref.WithHibernationBehaviourInclusiveSleep,
					},
				}),

				// ConcurrencyOptions:
				Entry(nil, &errorTE{
					marshalTE: marshalTE{
						given: "ConcurrencyOptions.NoW",
						option: func() pref.Option {
							return pref.WithNoW(5)
						},
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
})
