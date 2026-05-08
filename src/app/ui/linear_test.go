package ui_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/app/report"
	"github.com/snivilised/jaywalk/src/app/ui"
	"github.com/snivilised/jaywalk/src/prism"
)

// ---------------------------------------------------------------------------
// Spy renderer
// ---------------------------------------------------------------------------

// spyRenderer captures calls made to it by the linear presenter so
// that tests can assert on the translated prism types without depending
// on terminal output or lipgloss rendering.
type spyRenderer struct {
	overture prism.Overture
	motifs   []prism.Motif
	summary  prism.Summary

	beginCalled bool
	endCalled   bool
}

func (s *spyRenderer) Begin(o prism.Overture) {
	s.overture = o
	s.beginCalled = true
}

func (s *spyRenderer) Show(m prism.Motif) {
	s.motifs = append(s.motifs, m)
}

func (s *spyRenderer) End(su prism.Summary) {
	s.summary = su
	s.endCalled = true
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// newLinearWithSpy constructs a ui.linear backed by the given spy and
// returns it as a report.Presenter. This uses the exported
// ui.NewLinearWithRenderer constructor added specifically for testing
// so that spies can be injected without going through the registry.
func newLinearWithSpy(spy prism.Renderer) report.Presenter {
	return ui.NewLinearWithRenderer(spy)
}

// stubNode builds a minimal *core.Node for use in event structs.
func stubNode(path, name string, _ bool, depth uint) *core.Node {
	return &core.Node{
		Path: path,
		Extension: core.Extension{
			Name:  name,
			Depth: core.TraversalDepth(depth),
		},
		// IsDirectory is derived from the node type in real agenor - here
		// we set it directly via the test helper on the node struct.
	}
}

// ---------------------------------------------------------------------------
// Specs
// ---------------------------------------------------------------------------

var _ = Describe("linear", Ordered, func() {
	var (
		spy       *spyRenderer
		presenter report.Presenter
	)

	BeforeEach(func() {
		spy = &spyRenderer{}
		presenter = newLinearWithSpy(spy)
	})

	// ------------------------------------------------------------------
	// OnBegin
	// ------------------------------------------------------------------

	Describe("OnBegin", func() {
		Context("for a prime traversal", func() {
			It("passes PrimeNavigation kind and root to renderer.Begin", func() {
				now := core.Now()

				presenter.OnBegin(&report.BeginEvent{
					Root:      "/home/user/docs",
					Caption:   "files and folders",
					StartedAt: now,
					IsPrime:   true,
				})

				Expect(spy.beginCalled).To(BeTrue())
				Expect(spy.overture.Root).To(Equal("/home/user/docs"))
				Expect(spy.overture.Caption).To(Equal("files and folders"))
				Expect(spy.overture.StartedAt).To(Equal(now))
				Expect(spy.overture.Kind).To(Equal(prism.PrimeNavigation))
				Expect(spy.overture.ResumeFrom).To(BeEmpty())
			})
		})

		Context("for a resume traversal", func() {
			It("passes ResumeNavigation kind and resume path to renderer.Begin", func() {
				presenter.OnBegin(&report.BeginEvent{
					Root:       "/home/user/docs",
					Caption:    "files only",
					StartedAt:  core.Now(),
					IsPrime:    false,
					ResumeFrom: "/home/user/docs/subdir",
				})

				Expect(spy.overture.Kind).To(Equal(prism.ResumeNavigation))
				Expect(spy.overture.ResumeFrom).To(Equal("/home/user/docs/subdir"))
			})
		})
	})

	// ------------------------------------------------------------------
	// OnNodeEvent
	// ------------------------------------------------------------------

	Describe("OnNodeEvent", func() {
		DescribeTable("translates node fields into Motif correctly",
			func(path, name string, isDir bool, depth uint) {
				presenter.OnBegin(&report.BeginEvent{
					Root:    path,
					IsPrime: true,
				})

				node := stubNode(path, name, isDir, depth)
				presenter.OnNodeEvent(&report.NeutralEvent{
					DisplayEvent: report.DisplayEvent{Node: node},
				})

				Expect(spy.motifs).To(HaveLen(1))
				m := spy.motifs[0]
				Expect(m.Path).To(Equal(path))
				Expect(m.Name).To(Equal(name))
				Expect(m.Depth).To(Equal(depth))
				Expect(m.ActionName).To(BeEmpty())
				Expect(m.PipelineName).To(BeEmpty())
				Expect(m.Skipped).To(BeFalse())
				Expect(m.Err).To(BeNil())
			},
			Entry("root directory at depth 0",
				"/home/user/docs", "docs", true, uint(0),
			),
			Entry("file at depth 1",
				"/home/user/docs/report.pdf", "report.pdf", false, uint(1),
			),
			Entry("nested directory at depth 3",
				"/a/b/c/d", "d", true, uint(3),
			),
		)
	})

	// ------------------------------------------------------------------
	// OnActionEvent
	// ------------------------------------------------------------------

	Describe("OnActionEvent", func() {
		Context("when the action succeeds", func() {
			It("sets ActionName and leaves Err nil", func() {
				node := stubNode("/docs/file.mp4", "file.mp4", false, 1)
				presenter.OnActionEvent(&report.ActionEvent{
					DisplayEvent: report.DisplayEvent{
						Node: node,
						Name: "encode",
					},
					Err: nil,
				})

				Expect(spy.motifs).To(HaveLen(1))
				m := spy.motifs[0]
				Expect(m.ActionName).To(Equal("encode"))
				Expect(m.Err).To(BeNil())
			})
		})

		Context("when the action fails", func() {
			It("sets ActionName and Err on the Motif", func() {
				actionErr := errors.New("ffmpeg: codec not found")
				node := stubNode("/docs/file.mp4", "file.mp4", false, 1)

				presenter.OnActionEvent(&report.ActionEvent{
					DisplayEvent: report.DisplayEvent{
						Node: node,
						Name: "encode",
					},
					Err: actionErr,
				})

				Expect(spy.motifs).To(HaveLen(1))
				m := spy.motifs[0]
				Expect(m.ActionName).To(Equal("encode"))
				Expect(m.Err).To(MatchError(actionErr))
			})
		})
	})

	// ------------------------------------------------------------------
	// OnPipelineEvent
	// ------------------------------------------------------------------

	Describe("OnPipelineEvent", func() {
		Context("when the pipeline succeeds", func() {
			It("sets PipelineName and leaves Err nil", func() {
				node := stubNode("/docs/file.mp4", "file.mp4", false, 1)

				presenter.OnPipelineEvent(&report.PipelineEvent{
					DisplayEvent: report.DisplayEvent{
						Node: node,
						Name: "encode-and-upload",
					},
				})

				Expect(spy.motifs).To(HaveLen(1))
				m := spy.motifs[0]
				Expect(m.PipelineName).To(Equal("encode-and-upload"))
				Expect(m.Err).To(BeNil())
			})
		})
	})

	// ------------------------------------------------------------------
	// OnSkipEvent
	// ------------------------------------------------------------------

	Describe("OnSkipEvent", func() {
		It("sets Skipped, Placeholder and ResolvedPath on the Motif", func() {
			node := stubNode("/docs/file.mp4", "file.mp4", false, 1)

			presenter.OnSkipEvent(&report.SkipEvent{
				DisplayEvent: report.DisplayEvent{
					Node: node,
					Name: "encode",
				},
				Placeholder:  "{{.path}}",
				ResolvedPath: "/",
			})

			Expect(spy.motifs).To(HaveLen(1))
			m := spy.motifs[0]
			Expect(m.Skipped).To(BeTrue())
			Expect(m.Placeholder).To(Equal("{{.path}}"))
			Expect(m.ResolvedPath).To(Equal("/"))
			Expect(m.ActionName).To(Equal("encode"))
		})
	})

	// ------------------------------------------------------------------
	// OnComplete
	// ------------------------------------------------------------------

	Describe("OnComplete", func() {
		Context("for a successful prime traversal", func() {
			It("passes counts, elapsed, and PrimeNavigation kind to renderer.End", func() {
				presenter.OnBegin(&report.BeginEvent{
					Root:    "/docs",
					IsPrime: true,
				})

				presenter.OnComplete(&report.Traversal{
					FilesVisited: 42,
					DirsVisited:  7,
					Elapsed:      3 * time.Second,
				})

				Expect(spy.endCalled).To(BeTrue())
				s := spy.summary
				Expect(s.FilesVisited).To(Equal(uint(42)))
				Expect(s.DirsVisited).To(Equal(uint(7)))
				Expect(s.Elapsed).To(Equal(3 * time.Second))
				Expect(s.Errors).To(BeEmpty())
				Expect(s.Kind).To(Equal(prism.PrimeNavigation))
			})
		})

		Context("when the traversal ended with an error", func() {
			It("includes the error in Summary.Errors", func() {
				traversalErr := errors.New("permission denied")

				presenter.OnBegin(&report.BeginEvent{
					Root:    "/docs",
					IsPrime: true,
				})

				presenter.OnComplete(&report.Traversal{
					FilesVisited: 5,
					Err:          traversalErr,
				})

				s := spy.summary
				Expect(s.Errors).To(HaveLen(1))
				Expect(s.Errors[0]).To(MatchError(traversalErr))
			})
		})

		Context("for a resume traversal", func() {
			It("passes ResumeNavigation kind to renderer.End", func() {
				presenter.OnBegin(&report.BeginEvent{
					Root:    "/docs",
					IsPrime: false,
				})

				presenter.OnComplete(&report.Traversal{})

				Expect(spy.summary.Kind).To(Equal(prism.ResumeNavigation))
			})
		})
	})
})
