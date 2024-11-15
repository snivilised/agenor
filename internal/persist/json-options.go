package persist

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/internal/opts/json"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/pref"
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
				Files:       o.Sampling.NoOf.Files,
				Directories: o.Sampling.NoOf.Directories,
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
		File:      *NodeFilterDefToJSON(&poly.File),
		Directory: *NodeFilterDefToJSON(&poly.Directory),
	}
}

func FromJSON(jo *json.Options) *pref.Options {
	o := pref.DefaultOptions()

	o.Behaviours = pref.NavigationBehaviours{
		SubPath: pref.SubPathBehaviour{
			KeepTrailingSep: jo.Behaviours.SubPath.KeepTrailingSep,
		},
		Sort: pref.SortBehaviour{
			IsCaseSensitive: jo.Behaviours.Sort.IsCaseSensitive,
			SortFilesFirst:  jo.Behaviours.Sort.SortFilesFirst,
		},
		Cascade: pref.CascadeBehaviour{
			Depth:     jo.Behaviours.Cascade.Depth,
			NoRecurse: jo.Behaviours.Cascade.NoRecurse,
		},
	}
	o.Sampling = pref.SamplingOptions{
		Type:      jo.Sampling.Type,
		InReverse: jo.Sampling.InReverse,
		NoOf: pref.EntryQuantities{
			Files:       jo.Sampling.NoOf.Files,
			Directories: jo.Sampling.NoOf.Directories,
		},
	}
	o.Filter = pref.FilterOptions{
		Node: NodeFilterDefFromJSON(jo.Filter.Node),
		Child: lo.TernaryF(jo.Filter.Child != nil,
			func() *core.ChildFilterDef {
				return &core.ChildFilterDef{
					Type:        jo.Filter.Child.Type,
					Description: jo.Filter.Child.Description,
					Pattern:     jo.Filter.Child.Pattern,
					Negate:      jo.Filter.Child.Negate,
				}
			},
			func() *core.ChildFilterDef {
				return nil
			},
		),
		Sample: lo.TernaryF(jo.Filter.Sample != nil,
			func() *core.SampleFilterDef {
				return &core.SampleFilterDef{
					Type:        jo.Filter.Sample.Type,
					Description: jo.Filter.Sample.Description,
					Pattern:     jo.Filter.Sample.Pattern,
					Scope:       jo.Filter.Sample.Scope,
					Negate:      jo.Filter.Sample.Negate,
					// TODO:Poly: tbd,
				}
			},
			func() *core.SampleFilterDef {
				return nil
			},
		),
	}
	o.Hibernate = core.HibernateOptions{
		WakeAt:  NodeFilterDefFromJSON(jo.Hibernate.WakeAt),
		SleepAt: NodeFilterDefFromJSON(jo.Hibernate.SleepAt),
		Behaviour: core.HibernationBehaviour{
			InclusiveWake:  jo.Hibernate.Behaviour.InclusiveWake,
			InclusiveSleep: jo.Hibernate.Behaviour.InclusiveSleep,
		},
	}
	o.Concurrency = pref.ConcurrencyOptions{
		NoW: jo.Concurrency.NoW,
	}

	return o
}

func NodeFilterDefFromJSON(def *json.FilterDef) *core.FilterDef {
	return lo.TernaryF(def != nil,
		func() *core.FilterDef {
			return &core.FilterDef{
				Type:            def.Type,
				Description:     def.Description,
				Pattern:         def.Pattern,
				Negate:          def.Negate,
				Scope:           def.Scope,
				IfNotApplicable: def.IfNotApplicable,
				Poly:            NodePolyDefFromJSON(def.Poly),
			}
		},
		func() *core.FilterDef {
			return nil
		},
	)
}

func NodePolyDefFromJSON(poly *json.PolyFilterDef) *core.PolyFilterDef {
	if poly == nil {
		return nil
	}

	return &core.PolyFilterDef{
		File:      *NodeFilterDefFromJSON(&poly.File),
		Directory: *NodeFilterDefFromJSON(&poly.Directory),
	}
}
