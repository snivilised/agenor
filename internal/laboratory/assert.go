package lab

import (
	"path/filepath"
	"strings"
	"testing/fstest"

	. "github.com/onsi/gomega" //nolint:revive,stylecheck // ok
	"github.com/snivilised/nefilim/luna"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
)

type TestOptions struct {
	FS            *luna.MemFS
	Recording     RecordingMap
	Path          string
	Result        core.TraverseResult
	Err           error
	ExpectedErr   error
	Every         func(p string) bool
	ByPassMetrics bool
}

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
				Count: uint(lo.Sum(lo.Values(entry.ExpectedNoOf.Children))),
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
