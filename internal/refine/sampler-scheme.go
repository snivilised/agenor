package refine

import (
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/measure"
	"github.com/snivilised/traverse/nfs"
	"github.com/snivilised/traverse/pref"
)

type (
	scheme interface {
		create() error
		init(pi *types.PluginInit, owner *measure.Owned)
		next(node *core.Node, inspection core.Inspection) (bool, error)
	}
)

// need to narrow down the options
func newScheme(o *pref.Options) scheme {
	c := common{o: o}

	if o.Filter.IsNodeFilteringActive() {
		return &nativeScheme{
			common: c,
		}
	}

	if o.Filter.IsChildFilteringActive() {
		return &childScheme{
			common: c,
		}
	}

	if o.Filter.IsSampleFilteringActive() {
		return &samplerScheme{
			common: c,
		}
	}

	if o.Filter.IsCustomFilteringActive() {
		return &customScheme{
			common: c,
		}
	}

	return nil
}

type common struct {
	o     *pref.Options
	owner *measure.Owned
}

func (f *common) init(_ *types.PluginInit, owned *measure.Owned) {
	f.owner = owned
}

type nativeScheme struct {
	common
	filter core.TraverseFilter
}

func (f *nativeScheme) create() error {
	filter, err := newNodeFilter(f.o.Filter.Node, &f.o.Filter)
	if err != nil {
		return err
	}

	f.filter = filter

	if f.o.Filter.Sink != nil {
		f.o.Filter.Sink(pref.FilterReply{
			Node: f.filter,
		})
	}

	return nil
}

func (f *nativeScheme) next(node *core.Node, _ core.Inspection) (bool, error) {
	matched := f.filter.IsMatch(node)

	if !matched {
		filteredOutMetric := lo.Ternary(node.IsFolder(),
			enums.MetricNoFoldersFilteredOut,
			enums.MetricNoFilesFilteredOut,
		)
		f.owner.Mums[filteredOutMetric].Tick()
	}

	return matched, nil
}

type childScheme struct {
	common
	filter core.ChildTraverseFilter
}

func (f *childScheme) create() error {
	filter, err := newChildFilter(f.o.Filter.Child)

	if err != nil {
		return err
	}
	f.filter = filter

	if f.o.Filter.Sink != nil {
		f.o.Filter.Sink(pref.FilterReply{
			Child: f.filter,
		})
	}

	return nil
}

func (f *childScheme) init(pi *types.PluginInit, owner *measure.Owned) {
	f.common.init(pi, owner)

	// [KEEP-FILTER-IN-SYNC] keep this in sync with the default
	// behaviour in builders.override.Actions
	//
	pi.Actions.HandleChildren.Intercept(
		func(inspection core.Inspection, mums measure.MutableMetrics) {
			files := inspection.Sort(enums.EntryTypeFile)
			matching := f.filter.Matching(files)

			inspection.AssignChildren(matching)
			mums[enums.MetricNoChildFilesFound].Times(uint(len(files)))

			filteredOut := len(files) - len(matching)
			f.owner.Mums[enums.MetricNoChildFilesFilteredOut].Times(uint(filteredOut))
		},
	)
}

func (f *childScheme) next(_ *core.Node, _ core.Inspection) (bool, error) {
	return false, nil
}

type samplerScheme struct {
	common
	filter core.SampleTraverseFilter
}

func (f *samplerScheme) create() error {
	filter, err := newSampleFilter(f.o.Filter.Sample, &f.o.Sampling)

	if err != nil {
		return err
	}

	f.filter = filter

	// the filter plugin performs premature filtering (with fs.DirEntry as opposed
	// to core.Node) on behalf of the sampler.
	f.o.Hooks.ReadDirectory.Chain(
		func(result []fs.DirEntry, err error,
			_ fs.ReadDirFS, _ string,
		) ([]fs.DirEntry, error) {
			return f.filter.Matching(result), err
		})

	if f.o.Filter.Sink != nil {
		f.o.Filter.Sink(pref.FilterReply{
			Sampler: f.filter,
		})
	}

	return nil
}

func (f *samplerScheme) init(pi *types.PluginInit, owner *measure.Owned) {
	f.common.init(pi, owner)

	// [KEEP-FILTER-IN-SYNC] keep this in sync with the default
	// behaviour in builders.override.Actions
	//
	pi.Actions.HandleChildren.Intercept(
		func(inspection core.Inspection, _ measure.MutableMetrics) {
			files := inspection.Sort(enums.EntryTypeFile)
			matching := f.filter.Matching(files)

			inspection.AssignChildren(matching)
		},
	)
}

func (f *samplerScheme) next(node *core.Node, inspection core.Inspection) (bool, error) {
	if node.Extension.Scope.IsRoot() {
		matching := f.filter.Matching(
			[]fs.DirEntry{nfs.FromFileInfo(node.Info)},
		)
		result := len(matching) > 0

		lo.Ternary(result,
			f.owner.Mums[enums.MetricNoChildFilesFound],
			f.owner.Mums[enums.MetricNoChildFilesFilteredOut],
		).Times(uint(len(inspection.Contents().Files())))

		return result, nil
	}

	return true, nil
}

type customScheme struct {
	common
	filter core.TraverseFilter
}

func (f *customScheme) create() error {
	return f.o.Filter.Custom.Validate()
}

func (f *customScheme) next(_ *core.Node, _ core.Inspection) (bool, error) {
	if f.o.Filter.Sink != nil {
		f.o.Filter.Sink(pref.FilterReply{
			Node: f.filter,
		})
	}

	return false, nil
}
