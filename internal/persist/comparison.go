package persist

import (
	"fmt"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/opts/json"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

const (
	Anon = "ANON"
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

func (c *Comparison) Equals() error {
	return equalOptions(c.O, c.JO)
}

// equalOptions compare the pref.Options instance to the derived json instance json.Options.
// We can't use DeepEquals on the structs, because even though the structs
// may have the same members, DeepEqual will still fail because the host struct is
// different; eg: pref.NavigationBehaviours and json.NavigationBehaviours contain
// the same members, but they are different structs; which means comparison has to be
// done manually. equalOptions is only required because we have a custom mapping between
// pref.Options and json.Options, in the form of ToJSON/FromJSON
//
// We do need to apply the same technique to Active state, because there is no json
// version of ActiveState, so there is no custom functionality we need to check.
func equalOptions(o *pref.Options, jo *json.Options) error {
	if o == nil {
		return fmt.Errorf("pref.Options %w",
			UnequalPtrError[pref.Options, json.Options]{
				Field: "[nil pref.Options]",
				Value: o,
				Other: jo,
			},
		)
	}

	if jo == nil {
		return fmt.Errorf("json.Options %w",
			UnequalPtrError[pref.Options, json.Options]{
				Field: "[nil json.Options]",
				Value: o,
				Other: jo,
			},
		)
	}

	if err := equalBehaviours(&o.Behaviours, &jo.Behaviours); err != nil {
		return err
	}

	if err := equalSamplingOptions(&o.Sampling, &jo.Sampling); err != nil {
		return err
	}

	if err := equalFilterOptions(&o.Filter, &jo.Filter); err != nil {
		return err
	}

	if err := equalFilterDef("wake-at",
		o.Hibernate.WakeAt, jo.Hibernate.WakeAt,
	); err != nil {
		return err
	}

	if err := equalFilterDef("sleep-at",
		o.Hibernate.SleepAt, jo.Hibernate.SleepAt,
	); err != nil {
		return err
	}

	if o.Hibernate.Behaviour.InclusiveWake != jo.Hibernate.Behaviour.InclusiveWake {
		return fmt.Errorf("hibernate-behaviour %w", UnequalValueError[bool]{
			Field: "InclusiveWake",
			Value: o.Hibernate.Behaviour.InclusiveWake,
			Other: jo.Hibernate.Behaviour.InclusiveWake,
		})
	}

	if o.Hibernate.Behaviour.InclusiveSleep != jo.Hibernate.Behaviour.InclusiveSleep {
		return fmt.Errorf("hibernate-behaviour %w", UnequalValueError[bool]{
			Field: "InclusiveSleep",
			Value: o.Hibernate.Behaviour.InclusiveSleep,
			Other: jo.Hibernate.Behaviour.InclusiveSleep,
		})
	}

	if o.Concurrency.NoW != jo.Concurrency.NoW {
		return fmt.Errorf("concurrency %w", UnequalValueError[uint]{
			Field: "NoW",
			Value: o.Concurrency.NoW,
			Other: jo.Concurrency.NoW,
		})
	}

	return nil
}

func equalBehaviours(o *pref.NavigationBehaviours, jo *json.NavigationBehaviours) error {
	if o.SubPath.KeepTrailingSep != jo.SubPath.KeepTrailingSep {
		return fmt.Errorf("subPath %w", UnequalValueError[bool]{
			Field: "KeepTrailingSep",
			Value: o.SubPath.KeepTrailingSep,
			Other: jo.SubPath.KeepTrailingSep,
		})
	}

	if o.Sort.IsCaseSensitive != jo.Sort.IsCaseSensitive {
		return fmt.Errorf("sort %w", UnequalValueError[bool]{
			Field: "IsCaseSensitive",
			Value: o.Sort.IsCaseSensitive,
			Other: jo.Sort.IsCaseSensitive,
		})
	}

	if o.Sort.SortFilesFirst != jo.Sort.SortFilesFirst {
		return fmt.Errorf("sort %w", UnequalValueError[bool]{
			Field: "SortFilesFirst",
			Value: o.Sort.SortFilesFirst,
			Other: jo.Sort.SortFilesFirst,
		})
	}

	if o.Cascade.Depth != jo.Cascade.Depth {
		return fmt.Errorf("cascade %w", UnequalValueError[uint]{
			Field: "Depth",
			Value: o.Cascade.Depth,
			Other: jo.Cascade.Depth,
		})
	}

	if o.Cascade.NoRecurse != jo.Cascade.NoRecurse {
		return fmt.Errorf("cascade %w", UnequalValueError[bool]{
			Field: "NoRecurse",
			Value: o.Cascade.NoRecurse,
			Other: jo.Cascade.NoRecurse,
		})
	}

	// sort behaviour??

	return nil
}

func equalSamplingOptions(o *pref.SamplingOptions, jo *json.SamplingOptions) error {
	if o.Type != jo.Type {
		return fmt.Errorf("sampling %w", UnequalValueError[enums.SampleType]{
			Field: "Type",
			Value: o.Type,
			Other: jo.Type,
		})
	}

	if o.InReverse != jo.InReverse {
		return fmt.Errorf("sampling %w", UnequalValueError[bool]{
			Field: "InReverse",
			Value: o.InReverse,
			Other: jo.InReverse,
		})
	}

	if o.NoOf.Files != jo.NoOf.Files {
		return fmt.Errorf("sampling.noOf %w", UnequalValueError[uint]{
			Field: "Files",
			Value: o.NoOf.Files,
			Other: jo.NoOf.Files,
		})
	}

	if o.NoOf.Directories != jo.NoOf.Directories {
		return fmt.Errorf("sampling.noOf %w", UnequalValueError[uint]{
			Field: "Directories",
			Value: o.NoOf.Directories,
			Other: jo.NoOf.Directories,
		})
	}

	return nil
}

func equalFilterOptions(o *pref.FilterOptions, jo *json.FilterOptions) error {
	if err := equalFilterDef("node", o.Node, jo.Node); err != nil {
		return err
	}

	if err := equalChildFilterDef("child", o.Child, jo.Child); err != nil {
		return err
	}

	if err := equalSampleFilterDef("sample-filter", o.Sample, jo.Sample); err != nil {
		return err
	}

	return nil
}

func equalFilterDef(filterName string,
	def *core.FilterDef, jdef *json.FilterDef,
) error {
	if def == nil && jdef == nil {
		return nil
	}

	if def == nil && jdef != nil {
		return fmt.Errorf("filter-def %w",
			UnequalPtrError[core.FilterDef, json.FilterDef]{
				Field: "[nil def]",
				Value: def,
				Other: jdef,
			},
		)
	}

	if def != nil && jdef == nil {
		return fmt.Errorf("json-filter-def %w",
			UnequalPtrError[core.FilterDef, json.FilterDef]{
				Field: "[nil jdef]",
				Value: def,
				Other: jdef,
			},
		)
	}

	if def.Type != jdef.Type {
		return fmt.Errorf("%q filter-def %w", filterName,
			UnequalValueError[enums.FilterType]{
				Field: "Type",
				Value: def.Type,
				Other: jdef.Type,
			},
		)
	}

	if def.Description != jdef.Description {
		return fmt.Errorf("%q filter-def %w", filterName,
			UnequalValueError[string]{
				Field: "Description",
				Value: def.Description,
				Other: jdef.Description,
			},
		)
	}

	if def.Pattern != jdef.Pattern {
		return fmt.Errorf("%q filter-def %w", filterName,
			UnequalValueError[string]{
				Field: "Pattern",
				Value: def.Pattern,
				Other: jdef.Pattern,
			},
		)
	}

	if def.Scope != jdef.Scope {
		return fmt.Errorf("%q filter-def %w", filterName,
			UnequalValueError[enums.FilterScope]{
				Field: "Scope",
				Value: def.Scope,
				Other: jdef.Scope,
			},
		)
	}

	if def.Negate != jdef.Negate {
		return fmt.Errorf("%q filter-def %w", filterName,
			UnequalValueError[bool]{
				Field: "Negate",
				Value: def.Negate,
				Other: jdef.Negate,
			},
		)
	}

	if def.IfNotApplicable != jdef.IfNotApplicable {
		return fmt.Errorf("%q filter-def %w", filterName,
			UnequalValueError[enums.TriStateBool]{
				Field: "IfNotApplicable",
				Value: def.IfNotApplicable,
				Other: jdef.IfNotApplicable,
			},
		)
	}

	if def.Poly != nil && jdef.Poly != nil {
		if err := equalFilterDef("poly", &def.Poly.File, &jdef.Poly.File); err != nil {
			return err
		}

		if err := equalFilterDef("poly", &def.Poly.Directory, &jdef.Poly.Directory); err != nil {
			return err
		}
	}

	return nil
}

func equalChildFilterDef(filterName string,
	def *core.ChildFilterDef, jdef *json.ChildFilterDef,
) error {
	if def == nil && jdef == nil {
		return nil
	}

	if def == nil && jdef != nil {
		return fmt.Errorf("filter-def %w",
			UnequalPtrError[core.ChildFilterDef, json.ChildFilterDef]{
				Field: "[nil def]",
				Value: def,
				Other: jdef,
			},
		)
	}

	if def != nil && jdef == nil {
		return fmt.Errorf("filter-def %w",
			UnequalPtrError[core.ChildFilterDef, json.ChildFilterDef]{
				Field: "[nil jdef]",
				Value: def,
				Other: jdef,
			},
		)
	}

	if def.Type != jdef.Type {
		return fmt.Errorf("%q child-filter-def %w", filterName,
			UnequalValueError[enums.FilterType]{
				Field: "Type",
				Value: def.Type,
				Other: jdef.Type,
			},
		)
	}

	if def.Description != jdef.Description {
		return fmt.Errorf("%q child-filter-def %w", filterName,
			UnequalValueError[string]{
				Field: "Description",
				Value: def.Description,
				Other: jdef.Description,
			},
		)
	}

	if def.Pattern != jdef.Pattern {
		return fmt.Errorf("%q child-filter-def %w", filterName,
			UnequalValueError[string]{
				Field: "Pattern",
				Value: def.Pattern,
				Other: jdef.Pattern,
			},
		)
	}

	if def.Negate != jdef.Negate {
		return fmt.Errorf("%q child-filter-def %w", filterName,
			UnequalValueError[bool]{
				Field: "Negate",
				Value: def.Negate,
				Other: jdef.Negate,
			},
		)
	}

	return nil
}

func equalSampleFilterDef(filterName string,
	def *core.SampleFilterDef, jdef *json.SampleFilterDef,
) error {
	if def == nil && jdef == nil {
		return nil
	}

	if def == nil && jdef != nil {
		return fmt.Errorf("filter-def %w",
			UnequalPtrError[core.SampleFilterDef, json.SampleFilterDef]{
				Field: "[nil def]",
				Value: def,
				Other: jdef,
			},
		)
	}

	if def != nil && jdef == nil {
		return fmt.Errorf("filter-def %w",
			UnequalPtrError[core.SampleFilterDef, json.SampleFilterDef]{
				Field: "[nil jdef]",
				Value: def,
				Other: jdef,
			},
		)
	}

	if def.Type != jdef.Type {
		return fmt.Errorf("%q sample-filter-def %w", filterName,
			UnequalValueError[enums.FilterType]{
				Field: "Type",
				Value: def.Type,
				Other: jdef.Type,
			},
		)
	}

	if def.Description != jdef.Description {
		return fmt.Errorf("%q sample-filter-def %w", filterName,
			UnequalValueError[string]{
				Field: "Description",
				Value: def.Description,
				Other: jdef.Description,
			},
		)
	}

	if def.Pattern != jdef.Pattern {
		return fmt.Errorf("%q sample-filter-def %w", filterName,
			UnequalValueError[string]{
				Field: "Pattern",
				Value: def.Pattern,
				Other: jdef.Pattern,
			},
		)
	}

	if def.Scope != jdef.Scope {
		return fmt.Errorf("%q sample-filter-def %w", filterName,
			UnequalValueError[enums.FilterScope]{
				Field: "Scope",
				Value: def.Scope,
				Other: jdef.Scope,
			},
		)
	}

	if def.Negate != jdef.Negate {
		return fmt.Errorf("%q sample-filter-def %w", filterName,
			UnequalValueError[bool]{
				Field: "Negate",
				Value: def.Negate,
				Other: jdef.Negate,
			},
		)
	}

	if def.Poly != nil && jdef.Poly != nil {
		if err := equalFilterDef("poly", &def.Poly.File, &jdef.Poly.File); err != nil {
			return err
		}

		if err := equalFilterDef("poly", &def.Poly.Directory, &jdef.Poly.Directory); err != nil {
			return err
		}
	}

	return nil
}
