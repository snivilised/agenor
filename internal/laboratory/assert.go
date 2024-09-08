package lab

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing/fstest"

	. "github.com/onsi/gomega" //nolint:revive,stylecheck // ok
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
)

type TestOptions struct {
	FS          fstest.MapFS
	Recording   RecordingMap
	Path        string
	Result      core.TraverseResult
	Err         error
	ExpectedErr error
	Every       func(p string) bool
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

	assertMetrics(entry, to)
}

func assertMetrics(entry *NaviTE, to *TestOptions) {
	Expect(to.Result).To(
		And(
			HaveMetricCountOf(ExpectedMetric{
				Type:  enums.MetricNoFilesInvoked,
				Count: entry.ExpectedNoOf.Files,
			}),
			HaveMetricCountOf(ExpectedMetric{
				Type:  enums.MetricNoFoldersInvoked,
				Count: entry.ExpectedNoOf.Folders,
			}),
			HaveMetricCountOf(ExpectedMetric{
				Type:  enums.MetricNoChildFilesFound,
				Count: uint(lo.Sum(lo.Values(entry.ExpectedNoOf.Children))),
			}),
		),
	)
}

func subscribes(subscription enums.Subscription, de fs.DirEntry) bool {
	isUniversalSubscription := (subscription == enums.SubscribeUniversal)
	files := de != nil && (subscription == enums.SubscribeFiles) && (!de.IsDir())
	folders := de != nil && ((subscription == enums.SubscribeFolders) ||
		subscription == enums.SubscribeFoldersWithFiles) && (de.IsDir())

	return isUniversalSubscription || files || folders
}
