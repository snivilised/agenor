// Package hanno provides utilities for building a virtual file system tree
// based on an XML index file. It includes functions for creating a virtual
// tree from a predefined XML structure (like Musico) or from a custom XML
// file. The virtual tree is built in memory using a luna.MemFS, and the code
// allows for filtering which paths to include in the tree based on specified
// portions of the path. The IOProvider struct is used to abstract the reading
// and writing of files and directories, allowing for flexibility in how the
// virtual tree is constructed and displayed.
package hanno

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing/fstest"

	"github.com/snivilised/agenor/collections"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/nefilim/test/luna"
)

// Nuxx is a luna.MemFS factory hardcoded to Musico
func Nuxx(verbose bool, portions ...string) (fS *luna.MemFS) {
	fS = luna.NewMemFS()

	musico(
		newMemWriteProvider(fS, os.ReadFile, verbose, portions...),
		verbose,
	)

	return fS
}

// CustomTree is a luna.MemFS factory, equivalent to Nuxx that is populated by an
// alternative xml file. index is the full path the xml index file and
// tree is the name of the root element in the file (for Nuxx, this would be
// "MUSICO"). The tree can be filtered by specifying 'portions'.
func CustomTree(index, element string, verbose bool,
	portions ...string,
) (fS *luna.MemFS, err error) {
	fS = luna.NewMemFS()

	err = custom(
		index,
		element,
		newMemWriteProvider(fS, os.ReadFile, verbose, portions...),
		verbose,
	)

	return fS, err
}

func custom(index, tree string, provider *IOProvider, verbose bool) error {
	if _, err := os.Stat(index); err != nil {
		return err
	}

	if err := ensure(index, tree, provider, verbose); err != nil {
		return err
	}

	return nil
}

const (
	offset  = 2
	tabSize = 2
	doWrite = true
)

func musico(provider *IOProvider, verbose bool) {
	repo := Repo("")
	index := Combine(repo, "test/data/musico-index.xml")

	if err := ensure(index, "MUSICO", provider, verbose); err != nil {
		fmt.Printf("provision failed %v\n", err.Error())
	}
}

func ensure(index, tree string, provider *IOProvider, verbose bool) error {
	builder := virtualTree{
		tree:     tree,
		stack:    collections.NewStack[string](),
		index:    index,
		doWrite:  doWrite,
		provider: provider,
		verbose:  verbose,
		show: func(path string, exists existsEntry) {
			if !verbose {
				return
			}

			status := lo.Ternary(exists(path), "✅", "❌")

			fmt.Printf("---> %v path: '%v'\n", status, path)
		},
	}

	return builder.walk()
}

func newMemWriteProvider(fS *luna.MemFS,
	indexReader readFile,
	verbose bool,
	portions ...string,
) *IOProvider {
	filter := lo.Ternary(len(portions) > 0,
		matcher(func(path string) bool {
			for _, portion := range portions {
				if strings.Contains(path, portion) {
					return true
				}
			}

			return false
		}),
		matcher(func(string) bool {
			return true
		}),
	)

	if verbose {
		fmt.Printf("\n🤖 re-generating tree (filters: '%v')\n\n",
			strings.Join(portions, ", "),
		)
	}

	// PS: to check the existence of a path in an fs in production
	// code, use fs.Stat(fsys, path) instead of os.Stat/os.Lstat

	return &IOProvider{
		filter: filter,
		file: fileHandler{
			in: indexReader,
			out: writeFile(func(name string, data []byte, mode os.FileMode, show display) error {
				if name == "" {
					return nil
				}

				if filter(name) {
					trimmed := name
					fS.MapFS[trimmed] = &fstest.MapFile{
						Data: data,
						Mode: mode,
					}
					show(trimmed, func(path string) bool {
						entry, ok := fS.MapFS[path]
						return ok && !entry.Mode.IsDir()
					})
				}

				return nil
			}),
		},
		directory: directoryHandler{
			out: writeDirectory(func(path string, mode os.FileMode, show display, isRoot bool) error {
				if path == "" {
					return nil
				}

				if isRoot || filter(path) {
					trimmed := path
					fS.MapFS[trimmed] = &fstest.MapFile{
						Mode: mode | os.ModeDir,
					}
					show(trimmed, func(path string) bool {
						entry, ok := fS.MapFS[path]
						return ok && entry.Mode.IsDir()
					})
				}

				return nil
			}),
		},
	}
}

type (
	existsEntry func(path string) bool

	display func(path string, exists existsEntry)

	fileReader interface {
		read(name string) ([]byte, error)
	}

	readFile func(name string) ([]byte, error)

	fileWriter interface {
		write(name string, data []byte, perm os.FileMode, show display) error
	}

	writeFile func(name string, data []byte, perm os.FileMode, show display) error

	directoryWriter interface {
		write(path string, perm os.FileMode, show display, isRoot bool) error
	}

	writeDirectory func(path string, perm os.FileMode, show display, isRoot bool) error

	filter interface {
		match(portion string) bool
	}

	matcher func(portion string) bool

	fileHandler struct {
		in  fileReader
		out fileWriter
	}

	directoryHandler struct {
		out directoryWriter
	}

	// IOProvider provides the necessary handlers for reading and writing files
	// and directories, as well as a filter for determining which paths to include
	// when building the virtual tree. It is used by the virtualTree to interact
	// with the file system during the walk process.
	IOProvider struct {
		filter    filter
		file      fileHandler
		directory directoryHandler
	}

	// Tree represents the structure of the XML file used to build the virtual tree.
	Tree struct {
		XMLName xml.Name  `xml:"tree"`
		Root    Directory `xml:"directory"`
	}

	// Directory represents a directory in the virtual tree, which can contain files
	// and subdirectories.
	Directory struct {
		XMLName     xml.Name    `xml:"directory"`
		Name        string      `xml:"name,attr"`
		Files       []File      `xml:"file"`
		Directories []Directory `xml:"directory"`
	}

	// File represents a file in the virtual tree, which has a name and text content.
	File struct {
		XMLName xml.Name `xml:"file"`
		Name    string   `xml:"name,attr"`
		Text    string   `xml:",chardata"`
	}
)

func (fn readFile) read(name string) ([]byte, error) {
	return fn(name)
}

func (fn writeFile) write(name string, data []byte, perm os.FileMode, show display) error {
	return fn(name, data, perm, show)
}

func (fn writeDirectory) write(path string, perm os.FileMode, show display, isRoot bool) error {
	return fn(path, perm, show, isRoot)
}

func (fn matcher) match(portion string) bool {
	return fn(portion)
}

// virtualTree
type virtualTree struct {
	tree     string
	full     string
	stack    *collections.Stack[string]
	index    string
	doWrite  bool
	depth    int
	padding  string
	provider *IOProvider
	verbose  bool
	show     display
}

func (r *virtualTree) read() (*Directory, error) {
	data, err := r.provider.file.in.read(r.index)
	if err != nil {
		return nil, err
	}

	var tree Tree

	if ue := xml.Unmarshal(data, &tree); ue != nil {
		return nil, ue
	}

	return &tree.Root, nil
}

func (r *virtualTree) pad() string {
	return string(bytes.Repeat([]byte{' '}, (r.depth+offset)*tabSize))
}

func (r *virtualTree) refill() string {
	segments := r.stack.Content()
	return filepath.Join(segments...)
}

func (r *virtualTree) inc(name string) {
	r.stack.Push(name)
	r.full = r.refill()

	r.depth++
	r.padding = r.pad()
}

func (r *virtualTree) dec() {
	_, _ = r.stack.Pop()
	r.full = r.refill()

	r.depth--
	r.padding = r.pad()
}

func (r *virtualTree) walk() error {
	top, err := r.read()
	if err != nil {
		return err
	}

	r.full = r.tree

	return r.dir(*top, true)
}

func (r *virtualTree) dir(dir Directory, isRoot bool) error {
	if !isRoot {
		// We dont to add the root because only the descendents of the root
		// should be added
		r.inc(dir.Name)
	}

	if r.doWrite {
		if err := r.provider.directory.out.write(
			r.full,
			os.ModePerm,
			r.show,
			isRoot,
		); err != nil {
			return err
		}
	}

	for _, directory := range dir.Directories {
		if err := r.dir(directory, false); err != nil {
			return err
		}
	}

	for _, file := range dir.Files {
		full := Combine(r.full, file.Name)

		if r.doWrite {
			if err := r.provider.file.out.write(
				full,
				[]byte(file.Text),
				os.ModePerm,
				r.show,
			); err != nil {
				return err
			}
		}
	}

	r.dec()

	return nil
}

// Combine creates a path from the parent combined with the relative path. The relative
// path is a file system path so should only contain forward slashes, not the standard
// file path separator as denoted by filepath.Separator, typically used when interacting
// with the local file system. Do not use trailing "/".
func Combine(parent, relative string) string {
	if relative == "" {
		return parent
	}

	return parent + "/" + relative
}

// Repo gets the path of the repo with relative joined on
func Repo(relative string) string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, _ := cmd.Output()

	repo := strings.TrimSpace(string(output))

	return Combine(repo, relative)
}
