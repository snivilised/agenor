package persist

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/opts/json"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
)

// 📦 pkg: persist - defines marshalling functionality. This package is
// required in order to avoid the circular dependency that would be created
// if these functions were defined in opts.json.

const (
	JSONMarshalNoPrefix      = ""
	JSONMarshal2SpacesIndent = "  "
)

func ToJSON(o *pref.Options) *json.Options {
	return &json.Options{
		Behaviours: json.NavigationBehaviours{
			SubPath: json.SubPathBehaviour{
				KeepTrailingSep: o.Behaviours.SubPath.KeepTrailingSep,
			},
			Sort: json.SortBehaviour{
				IsCaseSensitive: o.Behaviours.Sort.IsCaseSensitive,
				SortFilesFirst:  o.Behaviours.Sort.SortFilesFirst,
			},
			Cascade: json.CascadeBehaviour{
				Depth:     o.Behaviours.Cascade.Depth,
				NoRecurse: o.Behaviours.Cascade.NoRecurse,
			},
		},
		Sampling: json.SamplingOptions{
			Type:      o.Sampling.Type,
			InReverse: o.Sampling.InReverse,
			NoOf: json.EntryQuantities{
				Files:   o.Sampling.NoOf.Files,
				Folders: o.Sampling.NoOf.Folders,
			},
		},
		Filter: json.FilterOptions{
			Node: NodeFilterDefToJSON(o.Filter.Node),
			Child: lo.TernaryF(o.Filter.Child != nil,
				func() *json.ChildFilterDef {
					return &json.ChildFilterDef{
						Type:        o.Filter.Child.Type,
						Description: o.Filter.Child.Description,
						Pattern:     o.Filter.Child.Pattern,
						Negate:      o.Filter.Child.Negate,
					}
				},
				func() *json.ChildFilterDef {
					return nil
				},
			),
			Sample: lo.TernaryF(o.Filter.Sample != nil,
				func() *json.SampleFilterDef {
					return &json.SampleFilterDef{
						Type:        o.Filter.Sample.Type,
						Description: o.Filter.Sample.Description,
						Pattern:     o.Filter.Sample.Pattern,
						Scope:       o.Filter.Sample.Scope,
						Negate:      o.Filter.Sample.Negate,
					}
				},
				func() *json.SampleFilterDef {
					return nil
				},
			),
		},
		Hibernate: json.HibernateOptions{
			WakeAt:  NodeFilterDefToJSON(o.Hibernate.WakeAt),
			SleepAt: NodeFilterDefToJSON(o.Hibernate.SleepAt),
			Behaviour: json.HibernationBehaviour{
				InclusiveWake:  o.Hibernate.Behaviour.InclusiveWake,
				InclusiveSleep: o.Hibernate.Behaviour.InclusiveSleep,
			},
		},
		Concurrency: json.ConcurrencyOptions{
			NoW: o.Concurrency.NoW,
		},
	}
}

func NodeFilterDefToJSON(def *core.FilterDef) *json.FilterDef {
	return lo.TernaryF(def != nil,
		func() *json.FilterDef {
			return &json.FilterDef{
				Type:            def.Type,
				Description:     def.Description,
				Pattern:         def.Pattern,
				Negate:          def.Negate,
				Scope:           def.Scope,
				IfNotApplicable: def.IfNotApplicable,
				Poly:            NodePolyDefToJSON(def.Poly),
			}
		},
		func() *json.FilterDef { return nil },
	)
}

func NodePolyDefToJSON(poly *core.PolyFilterDef) *json.PolyFilterDef {
	if poly == nil {
		return nil
	}

	return &json.PolyFilterDef{
		File:   *NodeFilterDefToJSON(&poly.File),
		Folder: *NodeFilterDefToJSON(&poly.Folder),
	}
}

func FromJSON(o *json.Options) *pref.Options {
	return &pref.Options{
		Behaviours: pref.NavigationBehaviours{
			SubPath: pref.SubPathBehaviour{
				KeepTrailingSep: o.Behaviours.SubPath.KeepTrailingSep,
			},
			Sort: pref.SortBehaviour{
				IsCaseSensitive: o.Behaviours.Sort.IsCaseSensitive,
				SortFilesFirst:  o.Behaviours.Sort.SortFilesFirst,
			},
			Cascade: pref.CascadeBehaviour{
				Depth:     o.Behaviours.Cascade.Depth,
				NoRecurse: o.Behaviours.Cascade.NoRecurse,
			},
		},
		Sampling: pref.SamplingOptions{
			Type:      o.Sampling.Type,
			InReverse: o.Sampling.InReverse,
			NoOf: pref.EntryQuantities{
				Files:   o.Sampling.NoOf.Files,
				Folders: o.Sampling.NoOf.Folders,
			},
		},
		Filter: pref.FilterOptions{
			Node: NodeFilterDefFromJSON(o.Filter.Node),
			Child: lo.TernaryF(o.Filter.Child != nil,
				func() *core.ChildFilterDef {
					return &core.ChildFilterDef{
						Type:        o.Filter.Child.Type,
						Description: o.Filter.Child.Description,
						Pattern:     o.Filter.Child.Pattern,
						Negate:      o.Filter.Child.Negate,
					}
				},
				func() *core.ChildFilterDef {
					return nil
				},
			),
			Sample: lo.TernaryF(o.Filter.Sample != nil,
				func() *core.SampleFilterDef {
					return &core.SampleFilterDef{
						Type:        o.Filter.Sample.Type,
						Description: o.Filter.Sample.Description,
						Pattern:     o.Filter.Sample.Pattern,
						Scope:       o.Filter.Node.Scope,
					}
				},
				func() *core.SampleFilterDef {
					return nil
				},
			),
		},
		Hibernate: core.HibernateOptions{
			WakeAt:  NodeFilterDefFromJSON(o.Hibernate.WakeAt),
			SleepAt: NodeFilterDefFromJSON(o.Hibernate.SleepAt),
			Behaviour: core.HibernationBehaviour{
				InclusiveWake:  o.Hibernate.Behaviour.InclusiveWake,
				InclusiveSleep: o.Hibernate.Behaviour.InclusiveSleep,
			},
		},
		Concurrency: pref.ConcurrencyOptions{
			NoW: o.Concurrency.NoW,
		},
	}
}

func NodeFilterDefFromJSON(def *json.FilterDef) *core.FilterDef {
	return lo.TernaryF(def != nil,
		func() *core.FilterDef {
			return &core.FilterDef{
				Type:        def.Type,
				Description: def.Description,
				Pattern:     def.Pattern,
				Negate:      def.Negate,
			}
		},
		func() *core.FilterDef {
			return nil
		},
	)
}
