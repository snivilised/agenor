package persist

import (
	"fmt"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/opts/json"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

type UnequalValueError[T any] struct {
	Field string
	Value T
	Other T
}

func (e UnequalValueError[T]) Error() string {
	s := fmt.Sprintf("unequal(%v) => value: %v, other: %v", e.Field, e.Value, e.Other)

	return s
}

func (UnequalValueError[T]) Unwrap() error {
	return locale.ErrUnEqualConversion
}

type UnequalPtrError[T any, O any] struct {
	Field string
	Value *T
	Other *O
}

func (e UnequalPtrError[T, O]) Error() string {
	s := fmt.Sprintf("unequal-ptr(%v) => value: %v, other: %v", e.Field, e.Value, e.Other)

	return s
}

func (UnequalPtrError[T, O]) Unwrap() error {
	return locale.ErrUnEqualConversion
}

// Equals compare the pref.Options instance to the derived json instance json.Options.
// We can't use DeepEquals because, on the structs, because even though the structs
// may have te same members, DeepEqual will still fail because the host struct is
// different; eg: pref.NavigationBehaviours and json.NavigationBehaviours contain
// the same members, but they are different structs; which means comparison has to be
// done manually.
func Equals(o *pref.Options, jo *json.Options) (bool, error) {
	if o == nil {
		return false, fmt.Errorf("pref.Options %w",
			UnequalPtrError[pref.Options, json.Options]{
				Field: "[nil pref.Options]",
				Value: o,
				Other: jo,
			},
		)
	}

	if jo == nil {
		return false, fmt.Errorf("json.Options %w",
			UnequalPtrError[pref.Options, json.Options]{
				Field: "[nil json.Options]",
				Value: o,
				Other: jo,
			},
		)
	}

	if equal, err := equalBehaviours(&o.Behaviours, &jo.Behaviours); !equal {
		return false, err
	}

	if equal, err := equalSampling(&o.Sampling, &jo.Sampling); !equal {
		return false, err
	}

	if equal, err := equalFilterOptions(&o.Filter, &jo.Filter); !equal {
		return false, err
	}

	if equal, err := equalFilterDef("wake-at", o.Hibernate.WakeAt, jo.Hibernate.WakeAt); !equal {
		return equal, err
	}

	if equal, err := equalFilterDef("sleep-at", o.Hibernate.SleepAt, jo.Hibernate.SleepAt); !equal {
		return equal, err
	}

	if o.Hibernate.Behaviour.InclusiveWake != jo.Hibernate.Behaviour.InclusiveWake {
		return false, fmt.Errorf("hibernate-behaviour %w", UnequalValueError[bool]{
			Field: "InclusiveWake",
			Value: o.Hibernate.Behaviour.InclusiveWake,
			Other: jo.Hibernate.Behaviour.InclusiveWake,
		})
	}

	if o.Hibernate.Behaviour.InclusiveSleep != jo.Hibernate.Behaviour.InclusiveSleep {
		return false, fmt.Errorf("hibernate-behaviour %w", UnequalValueError[bool]{
			Field: "InclusiveSleep",
			Value: o.Hibernate.Behaviour.InclusiveSleep,
			Other: jo.Hibernate.Behaviour.InclusiveSleep,
		})
	}

	if o.Concurrency.NoW != jo.Concurrency.NoW {
		return false, fmt.Errorf("concurrency %w", UnequalValueError[uint]{
			Field: "NoW",
			Value: o.Concurrency.NoW,
			Other: jo.Concurrency.NoW,
		})
	}

	return true, nil
}

func equalBehaviours(o *pref.NavigationBehaviours, jo *json.NavigationBehaviours) (bool, error) {
	if o.SubPath.KeepTrailingSep != jo.SubPath.KeepTrailingSep {
		return false, fmt.Errorf("subPath %w", UnequalValueError[bool]{
			Field: "SubPath",
			Value: o.SubPath.KeepTrailingSep,
			Other: jo.SubPath.KeepTrailingSep,
		})
	}

	if o.Sort.IsCaseSensitive != jo.Sort.IsCaseSensitive {
		return false, fmt.Errorf("sort %w", UnequalValueError[bool]{
			Field: "IsCaseSensitive",
			Value: o.Sort.IsCaseSensitive,
			Other: jo.Sort.IsCaseSensitive,
		})
	}

	if o.Cascade.Depth != jo.Cascade.Depth {
		return false, fmt.Errorf("cascade %w", UnequalValueError[uint]{
			Field: "Depth",
			Value: o.Cascade.Depth,
			Other: jo.Cascade.Depth,
		})
	}

	if o.Cascade.NoRecurse != jo.Cascade.NoRecurse {
		return false, fmt.Errorf("cascade %w", UnequalValueError[bool]{
			Field: "NoRecurse",
			Value: o.Cascade.NoRecurse,
			Other: jo.Cascade.NoRecurse,
		})
	}

	return true, nil
}

func equalSampling(o *pref.SamplingOptions, jo *json.SamplingOptions) (bool, error) {
	if o.Type != jo.Type {
		return false, fmt.Errorf("sampling %w", UnequalValueError[enums.SampleType]{
			Field: "SampleType",
			Value: o.Type,
			Other: jo.Type,
		})
	}

	if o.InReverse != jo.InReverse {
		return false, fmt.Errorf("sampling %w", UnequalValueError[bool]{
			Field: "SampleInReverse",
			Value: o.InReverse,
			Other: jo.InReverse,
		})
	}

	if o.NoOf.Files != jo.NoOf.Files {
		return false, fmt.Errorf("sampling.noOf %w", UnequalValueError[uint]{
			Field: "Files",
			Value: o.NoOf.Files,
			Other: jo.NoOf.Files,
		})
	}

	if o.NoOf.Folders != jo.NoOf.Folders {
		return false, fmt.Errorf("sampling.noOf %w", UnequalValueError[uint]{
			Field: "Folders",
			Value: o.NoOf.Folders,
			Other: jo.NoOf.Folders,
		})
	}

	return true, nil
}

func equalFilterOptions(o *pref.FilterOptions, jo *json.FilterOptions) (bool, error) {
	if equal, err := equalFilterDef("node", o.Node, jo.Node); !equal {
		return equal, err
	}

	if equal, err := equalChildFilterDef("child", o.Child, jo.Child); !equal {
		return equal, err
	}

	if equal, err := equalSampleFilterDef("sample-filter", o.Sample, jo.Sample); !equal {
		return equal, err
	}

	return true, nil
}

func equalFilterDef(filterName string,
	def *core.FilterDef, jdef *json.FilterDef,
) (bool, error) {
	if def == nil && jdef == nil {
		return true, nil
	}

	if def == nil && jdef != nil {
		return false, fmt.Errorf("filter-def %w",
			UnequalPtrError[core.FilterDef, json.FilterDef]{
				Field: "[nil def]",
				Value: def,
				Other: jdef,
			},
		)
	}

	if def != nil && jdef == nil {
		return false, fmt.Errorf("json-filter-def %w",
			UnequalPtrError[core.FilterDef, json.FilterDef]{
				Field: "[nil jdef]",
				Value: def,
				Other: jdef,
			},
		)
	}

	if def.Type != jdef.Type {
		return false, fmt.Errorf("%q filter-def %w", filterName,
			UnequalValueError[enums.FilterType]{
				Field: "Type",
				Value: def.Type,
				Other: jdef.Type,
			},
		)
	}

	if def.Description != jdef.Description {
		return false, fmt.Errorf("%q filter-def %w", filterName,
			UnequalValueError[string]{
				Field: "Description",
				Value: def.Description,
				Other: jdef.Description,
			},
		)
	}

	if def.Pattern != jdef.Pattern {
		return false, fmt.Errorf("%q filter-def %w", filterName,
			UnequalValueError[string]{
				Field: "Pattern",
				Value: def.Pattern,
				Other: jdef.Pattern,
			},
		)
	}

	if def.Scope != jdef.Scope {
		return false, fmt.Errorf("%q filter-def %w", filterName,
			UnequalValueError[enums.FilterScope]{
				Field: "Scope",
				Value: def.Scope,
				Other: jdef.Scope,
			},
		)
	}

	if def.Negate != jdef.Negate {
		return false, fmt.Errorf("%q filter-def %w", filterName,
			UnequalValueError[bool]{
				Field: "Negate",
				Value: def.Negate,
				Other: jdef.Negate,
			},
		)
	}

	if def.IfNotApplicable != jdef.IfNotApplicable {
		return false, fmt.Errorf("%q filter-def %w", filterName,
			UnequalValueError[enums.TriStateBool]{
				Field: "IfNotApplicable",
				Value: def.IfNotApplicable,
				Other: jdef.IfNotApplicable,
			},
		)
	}

	if def.Poly != nil && jdef.Poly != nil {
		if equal, err := equalFilterDef("poly", &def.Poly.File, &jdef.Poly.File); !equal {
			return equal, err
		}

		if equal, err := equalFilterDef("poly", &def.Poly.Folder, &jdef.Poly.Folder); !equal {
			return equal, err
		}
	}

	return true, nil
}

func equalChildFilterDef(filterName string,
	def *core.ChildFilterDef, jdef *json.ChildFilterDef,
) (bool, error) {
	if def == nil && jdef == nil {
		return true, nil
	}

	if def == nil && jdef != nil {
		return false, fmt.Errorf("filter-def %w",
			UnequalPtrError[core.ChildFilterDef, json.ChildFilterDef]{
				Field: "[nil def]",
				Value: def,
				Other: jdef,
			},
		)
	}

	if def != nil && jdef == nil {
		return false, fmt.Errorf("filter-def %w",
			UnequalPtrError[core.ChildFilterDef, json.ChildFilterDef]{
				Field: "[nil jdef]",
				Value: def,
				Other: jdef,
			},
		)
	}

	if def.Type != jdef.Type {
		return false, fmt.Errorf("%q child-filter-def %w", filterName,
			UnequalValueError[enums.FilterType]{
				Field: "Type",
				Value: def.Type,
				Other: jdef.Type,
			},
		)
	}

	if def.Description != jdef.Description {
		return false, fmt.Errorf("%q child-filter-def %w", filterName,
			UnequalValueError[string]{
				Field: "Description",
				Value: def.Description,
				Other: jdef.Description,
			},
		)
	}

	if def.Pattern != jdef.Pattern {
		return false, fmt.Errorf("%q child-filter-def %w", filterName,
			UnequalValueError[string]{
				Field: "Pattern",
				Value: def.Pattern,
				Other: jdef.Pattern,
			},
		)
	}

	if def.Negate != jdef.Negate {
		return false, fmt.Errorf("%q child-filter-def %w", filterName,
			UnequalValueError[bool]{
				Field: "Negate",
				Value: def.Negate,
				Other: jdef.Negate,
			},
		)
	}

	return true, nil
}

func equalSampleFilterDef(filterName string,
	def *core.SampleFilterDef, jdef *json.SampleFilterDef,
) (bool, error) {
	if def == nil && jdef == nil {
		return true, nil
	}

	if def == nil && jdef != nil {
		return false, fmt.Errorf("filter-def %w",
			UnequalPtrError[core.SampleFilterDef, json.SampleFilterDef]{
				Field: "[nil def]",
				Value: def,
				Other: jdef,
			},
		)
	}

	if def != nil && jdef == nil {
		return false, fmt.Errorf("filter-def %w",
			UnequalPtrError[core.SampleFilterDef, json.SampleFilterDef]{
				Field: "[nil jdef]",
				Value: def,
				Other: jdef,
			},
		)
	}

	if def.Type != jdef.Type {
		return false, fmt.Errorf("%q sample-filter-def %w", filterName,
			UnequalValueError[enums.FilterType]{
				Field: "Type",
				Value: def.Type,
				Other: jdef.Type,
			},
		)
	}

	if def.Description != jdef.Description {
		return false, fmt.Errorf("%q sample-filter-def %w", filterName,
			UnequalValueError[string]{
				Field: "Description",
				Value: def.Description,
				Other: jdef.Description,
			},
		)
	}

	if def.Pattern != jdef.Pattern {
		return false, fmt.Errorf("%q sample-filter-def %w", filterName,
			UnequalValueError[string]{
				Field: "Pattern",
				Value: def.Pattern,
				Other: jdef.Pattern,
			},
		)
	}

	if def.Scope != jdef.Scope {
		return false, fmt.Errorf("%q sample-filter-def %w", filterName,
			UnequalValueError[enums.FilterScope]{
				Field: "Scope",
				Value: def.Scope,
				Other: jdef.Scope,
			},
		)
	}

	if def.Negate != jdef.Negate {
		return false, fmt.Errorf("%q sample-filter-def %w", filterName,
			UnequalValueError[bool]{
				Field: "Negate",
				Value: def.Negate,
				Other: jdef.Negate,
			},
		)
	}

	return true, nil
}
