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
	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/helpers"
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/pref"
)

func TestKernel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kernel Suite")
}

const (
	RootPath    = "traversal-root-path"
	RestorePath = "/from-restore-path"
)

type recordingMap map[string]int
type recordingScopeMap map[string]enums.FilterScope
type recordingOrderMap map[string]int

type quantities struct {
	files    uint
	folders  uint
	children map[string]int
}

type naviTE struct {
	given         string
	should        string
	relative      string
	once          bool
	visit         bool
	caseSensitive bool
	subscription  enums.Subscription
	callback      core.Client
	mandatory     []string
	prohibited    []string
	expectedNoOf  quantities
}

type filterTE struct {
	naviTE
	name            string
	pattern         string
	scope           enums.FilterScope
	negate          bool
	expectedErr     error
	errorContains   string
	ifNotApplicable enums.TriStateBool
	custom          core.TraverseFilter
}

type polyTE struct {
	naviTE
	file   core.FilterDef
	folder core.FilterDef
}

type sampleTE struct {
	naviTE
	sampleType enums.SampleType
	reverse    bool
	filter     *filterTE
	noOf       pref.EntryQuantities
	each       pref.EachDirectoryEntryPredicate
	while      pref.WhileDirectoryPredicate
}

type customFilter struct {
	name            string
	pattern         string
	scope           enums.FilterScope
	negate          bool
	ifNotApplicable bool
}

// Description describes filter
func (f *customFilter) Description() string {
	return f.name
}

func (f *customFilter) Validate() {
	if f.scope == enums.ScopeUndefined {
		f.scope = enums.ScopeAll
	}
}

func (f *customFilter) Source() string {
	return f.pattern
}

func (f *customFilter) invert(result bool) bool {
	return lo.Ternary(f.negate, !result, result)
}

func (f *customFilter) IsMatch(node *core.Node) bool {
	if f.IsApplicable(node) {
		matched, _ := filepath.Match(f.pattern, node.Extension.Name)
		return f.invert(matched)
	}

	return f.ifNotApplicable
}

func (f *customFilter) IsApplicable(node *core.Node) bool {
	return (f.scope & node.Extension.Scope) > 0
}

func (f *customFilter) Scope() enums.FilterScope {
	return f.scope
}

func begin(em string) cycle.BeginHandler {
	return func(state *cycle.BeginState) {
		GinkgoWriter.Printf(
			"---> %v [traverse-navigator-test:BEGIN], root: '%v'\n", em, state.Root,
		)
	}
}

func end(em string) cycle.EndHandler {
	return func(result core.TraverseResult) {
		GinkgoWriter.Printf(
			"---> %v [traverse-navigator-test:END], err: '%v'\n", em, result.Error(),
		)
	}
}

func universalCallback(name string) core.Client {
	return func(node *core.Node) error {
		depth := node.Extension.Depth
		GinkgoWriter.Printf(
			"---> 🌊 UNIVERSAL//%v-CALLBACK: (depth:%v) '%v'\n", name, depth, node.Path,
		)
		Expect(node.Extension).NotTo(BeNil(), helpers.Reason(node.Path))

		return nil
	}
}

func foldersCallback(name string) core.Client {
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

func filesCallback(name string) core.Client {
	return func(node *core.Node) error {
		GinkgoWriter.Printf("---> 🌙 FILES//%v-CALLBACK: '%v'\n", name, node.Path)
		Expect(node.IsFolder()).To(BeFalse(),
			helpers.Because(node.Path, "node expected to be file"),
		)
		Expect(node.Extension).NotTo(BeNil(), helpers.Reason(node.Path))

		return nil
	}
}

func foldersCaseSensitiveCallback(first, second string) core.Client {
	recording := make(recordingMap)

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

type testOptions struct {
	vfs       fstest.MapFS
	recording recordingMap
	path      string
	result    core.TraverseResult
	err       error
	every     func(p string) bool
}

func assertNavigation(entry *naviTE, to testOptions) {
	Expect(to.err).To(Succeed())

	visited := []string{}
	_ = to.result.Session().StartedAt()
	_ = to.result.Session().Elapsed()

	if entry.visit {
		_ = fs.WalkDir(to.vfs, to.path, func(path string, de fs.DirEntry, _ error) error {
			if strings.HasSuffix(path, ".DS_Store") {
				return nil
			}

			if subscribes(entry.subscription, de) {
				visited = append(visited, path)
			}

			return nil
		})

		every := lo.EveryBy(visited,
			lo.Ternary(to.every != nil, to.every, func(p string) bool {
				segments := strings.Split(p, string(filepath.Separator))
				name, err := lo.Last(segments)

				if err == nil {
					_, found := to.recording[name]
					return found
				}

				return false
			}),
		)

		Expect(every).To(BeTrue())
	}

	for n, actualNoChildren := range entry.expectedNoOf.children {
		expected := to.recording[n]
		Expect(to.recording[n]).To(Equal(actualNoChildren),
			helpers.BecauseQuantity(fmt.Sprintf("folder: '%v'", n),
				expected,
				actualNoChildren,
			),
		)
	}

	if entry.mandatory != nil {
		for _, name := range entry.mandatory {
			_, found := to.recording[name]
			Expect(found).To(BeTrue(), helpers.Reason(name))
		}
	}

	if entry.prohibited != nil {
		for _, name := range entry.prohibited {
			_, found := to.recording[name]
			Expect(found).To(BeFalse(), helpers.Reason(name))
		}
	}

	assertMetrics(entry, to)
}

func assertMetrics(entry *naviTE, to testOptions) {
	Expect(to.result.Metrics().Count(enums.MetricNoFilesInvoked)).To(
		Equal(entry.expectedNoOf.files),
		helpers.BecauseQuantity("Incorrect no of files",
			int(entry.expectedNoOf.files),
			int(to.result.Metrics().Count(enums.MetricNoFilesInvoked)),
		),
	)

	Expect(to.result.Metrics().Count(enums.MetricNoFoldersInvoked)).To(
		Equal(entry.expectedNoOf.folders),
		helpers.BecauseQuantity("Incorrect no of folders",
			int(entry.expectedNoOf.folders),
			int(to.result.Metrics().Count(enums.MetricNoFoldersInvoked)),
		),
	)

	sum := lo.Sum(lo.Values(entry.expectedNoOf.children))

	Expect(to.result.Metrics().Count(enums.MetricNoChildFilesFound)).To(
		Equal(uint(sum)),
		helpers.BecauseQuantity("Incorrect total no of child files",
			sum,
			int(to.result.Metrics().Count(enums.MetricNoChildFilesFound)),
		),
	)
}
