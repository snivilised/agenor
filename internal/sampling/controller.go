package sampling

import (
	"fmt"
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/nfs"
	"github.com/snivilised/traverse/pref"
)

type controller struct {
	o      *samplingOptions
	filter core.ChildTraverseFilter
}

func (p *controller) Role() enums.Role {
	return enums.RoleSampler
}

func (p *controller) Next(node *core.Node, inspection core.Inspection) (bool, error) {
	_ = inspection
	// ensure invoke count is correct
	//
	fmt.Printf("sampler🌀 ~ name(Next): '%v' true\n", node.Extension.Name)
	return true, nil
}

func (p *controller) sample(result []fs.DirEntry, _ error,
	_ fs.ReadDirFS, _ string,
) ([]fs.DirEntry, error) {
	files, folders := nfs.Separate(result)

	return union(&readResult{
		files:   files,
		folders: folders,
		o:       p.o.sampling,
	}), nil
}

type (
	samplerFunc func(n uint, entries []fs.DirEntry) []fs.DirEntry
)

type readResult struct {
	files   []fs.DirEntry
	folders []fs.DirEntry
	o       *pref.SamplingOptions
}

func union(r *readResult) []fs.DirEntry {
	noOfFiles := lo.Ternary(r.o.NoOf.Files == 0,
		uint(len(r.files)), r.o.NoOf.Files,
	)

	both := lo.Ternary(
		r.o.SampleInReverse, last, first,
	)(noOfFiles, r.files)

	noOfFolders := lo.Ternary(r.o.NoOf.Folders == 0,
		uint(len(r.folders)), r.o.NoOf.Folders,
	)
	both = append(both, lo.Ternary(
		r.o.SampleInReverse, last, first,
	)(noOfFolders, r.folders)...)

	return both
}

func first(n uint, entries []fs.DirEntry) []fs.DirEntry {
	return entries[:(min(n, uint(len(entries))))]
}

func last(n uint, entries []fs.DirEntry) []fs.DirEntry {
	return entries[uint(len(entries))-(min(n, uint(len(entries)))):]
}
