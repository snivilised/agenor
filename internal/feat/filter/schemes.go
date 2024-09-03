package filter

import (
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/filtering"
	"github.com/snivilised/traverse/internal/measure"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/nfs"
	"github.com/snivilised/traverse/pref"
)

type (
	scheme interface {
		create() error
		init(pi *types.PluginInit, crate *measure.Crate)
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
	crate *measure.Crate
}

func (f *common) init(_ *types.PluginInit, crate *measure.Crate) {
	f.crate = crate
}

type nativeScheme struct {
	common
	filter core.TraverseFilter
}

func (f *nativeScheme) create() error {
	filter, err := filtering.NewNodeFilter(f.o.Filter.Node, &f.o.Filter)
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
	return matchNext(f.filter, node, f.crate)
}

type childScheme struct {
	common
	filter core.ChildTraverseFilter
}

func (f *childScheme) create() error {
	filter, err := filtering.NewChildFilter(f.o.Filter.Child)

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

func (f *childScheme) init(pi *types.PluginInit, crate *measure.Crate) {
	f.common.init(pi, crate)

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
			f.crate.Mums[enums.MetricNoChildFilesFilteredOut].Times(uint(filteredOut))
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
	filter, err := filtering.NewSampleFilter(f.o.Filter.Sample, &f.o.Sampling)

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

func (f *samplerScheme) init(pi *types.PluginInit, crate *measure.Crate) {
	f.common.init(pi, crate)

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
			f.crate.Mums[enums.MetricNoChildFilesFound],
			f.crate.Mums[enums.MetricNoChildFilesFilteredOut],
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
	f.filter = f.o.Filter.Custom

	if f.o.Filter.Sink != nil {
		f.o.Filter.Sink(pref.FilterReply{
			Node: f.filter,
		})
	}

	return f.filter.Validate()
}

func (f *customScheme) next(node *core.Node, _ core.Inspection) (bool, error) {
	return matchNext(f.filter, node, f.crate)
}

func matchNext(filter core.TraverseFilter, node *core.Node, crate *measure.Crate) (bool, error) {
	matched := filter.IsMatch(node)

	if !matched {
		filteredOutMetric := lo.Ternary(node.IsFolder(),
			enums.MetricNoFoldersFilteredOut,
			enums.MetricNoFilesFilteredOut,
		)
		crate.Mums[filteredOutMetric].Tick()
	}

	return matched, nil
}
