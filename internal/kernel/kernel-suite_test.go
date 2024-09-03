package kernel_test

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/filtering"
	"github.com/snivilised/traverse/internal/helpers"
	"github.com/snivilised/traverse/internal/third/lo"
)

func TestKernel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kernel Suite")
}

const (
	RootPath    = "traversal-root-path"
	RestorePath = "/from-restore-path"
)

func UniversalCallback(name string) core.Client {
	return func(node *core.Node) error {
		depth := node.Extension.Depth
		GinkgoWriter.Printf(
			"---> 🌊 UNIVERSAL//%v-CALLBACK: (depth:%v) '%v'\n", name, depth, node.Path,
		)
		Expect(node.Extension).NotTo(BeNil(), helpers.Reason(node.Path))

		return nil
	}
}

func FoldersCallback(name string) core.Client {
	return func(node *core.Node) error {
		depth := node.Extension.Depth
		actualNoChildren := len(node.Children)
		GinkgoWriter.Printf(
			"---> 🔆 FOLDERS//CALLBACK%v: (depth:%v, children:%v) '%v'\n",
			name, depth, actualNoChildren, node.Path,
		)
		Expect(node.IsFolder()).To(BeTrue(),
			helpers.Because(node.Path, "node expected to be folder"),
		)
		Expect(node.Extension).NotTo(BeNil(), helpers.Reason(node.Path))

		return nil
	}
}

func FilesCallback(name string) core.Client {
	return func(node *core.Node) error {
		GinkgoWriter.Printf("---> 🌙 FILES//%v-CALLBACK: '%v'\n", name, node.Path)
		Expect(node.IsFolder()).To(BeFalse(),
			helpers.Because(node.Path, "node expected to be file"),
		)
		Expect(node.Extension).NotTo(BeNil(), helpers.Reason(node.Path))

		return nil
	}
}

func FoldersCaseSensitiveCallback(first, second string) core.Client {
	recording := make(helpers.RecordingMap)

	return func(node *core.Node) error {
		recording[node.Path] = len(node.Children)

		GinkgoWriter.Printf("---> 🔆 CASE-SENSITIVE-CALLBACK: '%v'\n", node.Path)
		Expect(node.IsFolder()).To(BeTrue())

		if strings.HasSuffix(node.Path, second) {
			GinkgoWriter.Printf("---> 💧 FIRST: '%v', 💧 SECOND: '%v'\n", first, second)

			paths := lo.Keys(recording)
			_, found := lo.Find(paths, func(s string) bool {
				return strings.HasSuffix(s, first)
			})

			Expect(found).To(BeTrue(), fmt.Sprintf("for node: '%v'", node.Extension.Name))
		}

		return nil
	}
}

func subscribes(subscription enums.Subscription, de fs.DirEntry) bool {
	isAnySubscription := (subscription == enums.SubscribeUniversal)
	files := de != nil && (subscription == enums.SubscribeFiles) && (!de.IsDir())
	folders := de != nil && ((subscription == enums.SubscribeFolders) ||
		subscription == enums.SubscribeFoldersWithFiles) && (de.IsDir())

	return isAnySubscription || files || folders
}

type TestOptions struct {
	FS          fstest.MapFS
	Recording   helpers.RecordingMap
	Path        string
	Result      core.TraverseResult
	Err         error
	ExpectedErr error
	Every       func(p string) bool
}

func AssertNavigation(entry *helpers.NaviTE, to *TestOptions) {
	if to.ExpectedErr != nil {
		Expect(to.Err).To(MatchError(to.ExpectedErr))
		return
	}

	Expect(to.Err).To(Succeed())

	visited := []string{}
	_ = to.Result.Session().StartedAt()
	_ = to.Result.Session().Elapsed()

	if entry.Visit {
		_ = fs.WalkDir(to.FS, to.Path, func(path string, de fs.DirEntry, _ error) error {
			if strings.HasSuffix(path, ".DS_Store") {
				return nil
			}

			if subscribes(entry.Subscription, de) {
				visited = append(visited, path)
			}

			return nil
		})

		every := lo.EveryBy(visited,
			lo.Ternary(to.Every != nil, to.Every, func(p string) bool {
				segments := strings.Split(p, string(filepath.Separator))
				name, err := lo.Last(segments)

				if err == nil {
					_, found := to.Recording[name]
					return found
				}

				return false
			}),
		)

		Expect(every).To(BeTrue())
	}

	for n, actualNoChildren := range entry.ExpectedNoOf.Children {
		expected := to.Recording[n]
		Expect(to.Recording[n]).To(Equal(actualNoChildren),
			helpers.BecauseQuantity(fmt.Sprintf("folder: '%v'", n),
				expected,
				actualNoChildren,
			),
		)
	}

	if entry.Mandatory != nil {
		for _, name := range entry.Mandatory {
			_, found := to.Recording[name]
			Expect(found).To(BeTrue(), helpers.Reason(name))
		}
	}

	if entry.Prohibited != nil {
		for _, name := range entry.Prohibited {
			_, found := to.Recording[name]
			Expect(found).To(BeFalse(), helpers.Reason(name))
		}
	}

	assertMetrics(entry, to)
}

func assertMetrics(entry *helpers.NaviTE, to *TestOptions) {
	Expect(to.Result.Metrics().Count(enums.MetricNoFilesInvoked)).To(
		Equal(entry.ExpectedNoOf.Files),
		helpers.BecauseQuantity("Incorrect no of files",
			int(entry.ExpectedNoOf.Files),                              //nolint:gosec // ok
			int(to.Result.Metrics().Count(enums.MetricNoFilesInvoked)), //nolint:gosec // ok
		),
	)

	Expect(to.Result.Metrics().Count(enums.MetricNoFoldersInvoked)).To(
		Equal(entry.ExpectedNoOf.Folders),
		helpers.BecauseQuantity("Incorrect no of folders",
			int(entry.ExpectedNoOf.Folders),                              //nolint:gosec // ok
			int(to.Result.Metrics().Count(enums.MetricNoFoldersInvoked)), //nolint:gosec // ok
		),
	)

	sum := lo.Sum(lo.Values(entry.ExpectedNoOf.Children))

	Expect(to.Result.Metrics().Count(enums.MetricNoChildFilesFound)).To(
		Equal(uint(sum)),
		helpers.BecauseQuantity("Incorrect total no of child files",
			sum,
			int(to.Result.Metrics().Count(enums.MetricNoChildFilesFound)), //nolint:gosec // ok
		),
	)
}

// customSamplingFilter is a custom sampling filter that just happens
// to use a glob as part of its implementation. The client can of course
// define their own custom implementation using filter.SampleFilter.
type customSamplingFilter struct {
	filtering.Sample
	description string
	pattern     string
}

func (f *customSamplingFilter) Description() string {
	return f.description
}

func (f *customSamplingFilter) Scope() enums.FilterScope {
	return f.Sample.Scope()
}

func (f *customSamplingFilter) Matching(children []fs.DirEntry) []fs.DirEntry {
	return f.Sample.Matching(children,
		func(entry fs.DirEntry, _ int) bool {
			matched, _ := filepath.Match(f.pattern, entry.Name())
			return matched
		},
	)
}
