package lab

import (
	"path/filepath"
	"strings"
	"testing/fstest"

	. "github.com/onsi/gomega" //nolint:staticcheck // ok
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/internal/third/lo"
	"github.com/snivilised/nefilim/test/luna"
)

// TestOptions defines the options that can be usee for a unit test
type TestOptions struct {
	// FS the memory based file system
	FS *luna.MemFS

	// Record captures navigated nodes
	Recording Recall

	// Path where to navigate from
	Path string

	// Result is the navigation result
	Result core.TraverseResult

	// Err is the actual error returned from the navigation
	Err error

	// ExpectedErr is the error expected by the unit test
	ExpectedErr error

	// Every is an optional predicate function invoked for every
	// node visited.
	Every func(p string) bool

	// ByPassMetrics determines whether we by pass the assertions
	// that check metrics
	ByPassMetrics bool
}

// AssertNavigation is a compound asserter for navigation tests
func AssertNavigation(entry *NaviTE, to *TestOptions) {
	if to.ExpectedErr != nil {
		Expect(to.Err).To(MatchError(to.ExpectedErr))
		return
	}

	Expect(to.Err).To(Succeed())

	visited := []string{}
	_ = to.Result.Session().StartedAt()
	_ = to.Result.Session().Elapsed()

	if entry.Visit && to.FS != nil {
		for path, file := range to.FS.MapFS {
			if strings.HasPrefix(path, to.Path) {
				if subscribes(entry.Subscription, file) {
					visited = append(visited, path)
				}
			}
		}

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

		Expect(every).To(BeTrue(), "Not all expected items were invoked")
	}

	for name, expected := range entry.ExpectedNoOf.Children {
		Expect(to.Recording).To(HaveChildCountOf(ExpectedCount{
			Name:  name,
			Count: expected,
		}))
	}

	if entry.Mandatory != nil {
		for _, name := range entry.Mandatory {
			Expect(to.Recording).To(HaveInvokedNode(name))
		}
	}

	if entry.Prohibited != nil {
		for _, name := range entry.Prohibited {
			Expect(to.Recording).To(HaveNotInvokedNode(name))
		}
	}

	if !to.ByPassMetrics {
		assertMetrics(entry, to)
	}
}

func assertMetrics(entry *NaviTE, to *TestOptions) {
	Expect(to.Result).To(
		And(
			HaveMetricCountOf(ExpectedMetric{
				Type:  enums.MetricNoFilesInvoked,
				Count: entry.ExpectedNoOf.Files,
			}),
			HaveMetricCountOf(ExpectedMetric{
				Type:  enums.MetricNoDirectoriesInvoked,
				Count: entry.ExpectedNoOf.Directories,
			}),
			HaveMetricCountOf(ExpectedMetric{
				Type:  enums.MetricNoChildFilesFound,
				Count: uint(lo.Sum(lo.Values(entry.ExpectedNoOf.Children))), //nolint:gosec // ok
			}),
		),
	)
}

func subscribes(subscription enums.Subscription, mapFile *fstest.MapFile) bool {
	isUniversalSubscription := (subscription == enums.SubscribeUniversal)
	files := mapFile != nil && (subscription == enums.SubscribeFiles) && !mapFile.Mode.IsDir()
	directories := mapFile != nil && ((subscription == enums.SubscribeDirectories) ||
		subscription == enums.SubscribeDirectoriesWithFiles) && mapFile.Mode.IsDir()

	return isUniversalSubscription || files || directories
}
