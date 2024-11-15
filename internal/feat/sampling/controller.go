package sampling

import (
	"io/fs"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/pref"
	nef "github.com/snivilised/nefilim"
)

type controller struct {
	o      *pref.SamplingOptions
	filter core.ChildTraverseFilter
}

func (p *controller) Role() enums.Role {
	return enums.RoleSampler
}

func (p *controller) Next(_ core.Servant, _ enclave.Inspection) (bool, error) {
	return true, nil
}

func (p *controller) sample(result []fs.DirEntry, _ error,
	_ fs.ReadDirFS, _ string,
) ([]fs.DirEntry, error) {
	files, directories := nef.Separate(result)

	return union(&readResult{
		files:       files,
		directories: directories,
		o:           p.o,
	}), nil
}

type (
	samplerFunc func(n uint, entries []fs.DirEntry) []fs.DirEntry
)

type readResult struct {
	files       []fs.DirEntry
	directories []fs.DirEntry
	o           *pref.SamplingOptions
}

func union(r *readResult) []fs.DirEntry {
	noOfFiles := lo.Ternary(r.o.NoOf.Files == 0,
		uint(len(r.files)), r.o.NoOf.Files,
	)

	both := lo.Ternary(
		r.o.InReverse, last, first,
	)(noOfFiles, r.files)

	noOfDirectories := lo.Ternary(r.o.NoOf.Directories == 0,
		uint(len(r.directories)), r.o.NoOf.Directories,
	)
	both = append(both, lo.Ternary(
		r.o.InReverse, last, first,
	)(noOfDirectories, r.directories)...)

	return both
}

func first(n uint, entries []fs.DirEntry) []fs.DirEntry {
	return entries[:(min(n, uint(len(entries))))]
}

func last(n uint, entries []fs.DirEntry) []fs.DirEntry {
	return entries[uint(len(entries))-(min(n, uint(len(entries)))):]
}
