package flow_test

import (
	"bytes"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/prism"
	"github.com/snivilised/jaywalk/src/prism/flow"
)

var _ = Describe("LinearRenderer", func() {
	It("renders tree branches with default icons", func() {
		w := &bytes.Buffer{}
		palette := prism.Palette{}

		renderer, err := flow.New(palette, w)
		Expect(err).To(BeNil())
		Expect(renderer).NotTo(BeNil())

		renderer.Show(prism.Motif{
			Name:        "app",
			IsDir:       true,
			Depth:       0,
			VisualDepth: 0,
			IsLast:      true,
		})

		renderer.Show(prism.Motif{
			Name:        "bedrock",
			IsDir:       true,
			Depth:       1,
			VisualDepth: 1,
			IsLast:      false,
		})

		renderer.Show(prism.Motif{
			Name:        "bedrock_suite_test.go",
			IsDir:       false,
			Depth:       1,
			VisualDepth: 2,
			IsLast:      true,
		})

		output := w.String()
		Expect(output).To(ContainSubstring("✻ app/\n"))
		Expect(output).To(ContainSubstring("├── 📁 bedrock/\n"))
		Expect(output).To(ContainSubstring("│  └── 🔖 bedrock_suite_test.go\n"))
	})

	It("applies WithIcons overrides", func() {
		w := &bytes.Buffer{}
		palette := prism.Palette{}

		renderer, err := flow.New(palette, w, flow.WithIcons(map[string]string{
			prism.TreeIconRoot:           "R",
			prism.TreeIconDirectory:      "D",
			prism.TreeIconFile:           "F",
			prism.TreeIconElapsed:        "E",
			prism.TreeIconBranchJoint:    "+-- ",
			prism.TreeIconBranchLast:     "L-- ",
			prism.TreeIconBranchVertical: "|",
			prism.TreeIconBranchIndent:   "  ",
		}))
		Expect(err).To(BeNil())

		renderer.Show(prism.Motif{
			Name:        "root",
			IsDir:       true,
			Depth:       0,
			VisualDepth: 0,
			IsLast:      true,
		})

		renderer.Show(prism.Motif{
			Name:        "child",
			IsDir:       false,
			Depth:       1,
			VisualDepth: 1,
			IsLast:      true,
		})

		Expect(w.String()).To(ContainSubstring("R root/\n"))
		Expect(w.String()).To(ContainSubstring("L-- F child\n"))
	})

	It("renders summary entries with tree icon prefixes", func() {
		w := &bytes.Buffer{}
		palette := prism.Palette{}

		renderer, err := flow.New(palette, w)
		Expect(err).To(BeNil())

		renderer.End(prism.Summary{
			Kind:         prism.PrimeNavigation,
			FilesVisited: 12,
			DirsVisited:  3,
			Elapsed:      2 * time.Second,
		})

		output := w.String()
		Expect(output).To(ContainSubstring("🔖 Files"))
		Expect(output).To(ContainSubstring("📁 Directories"))
		Expect(output).To(ContainSubstring("⏰ Elapsed"))
	})

	It("returns a renderer when options are provided", func() {
		w := &bytes.Buffer{}
		palette := prism.Palette{}

		renderer, err := flow.New(palette, w, flow.WithIcons(nil))
		Expect(err).To(BeNil())
		Expect(renderer).NotTo(BeNil())
		Expect(renderer).To(Not(BeNil()))
		renderer.Show(prism.Motif{Name: "test", Depth: 0, VisualDepth: 0, IsDir: true, IsLast: true})
		Expect(w.String()).To(ContainSubstring("✻ test/"))
	})

	It("renders the banner inside the summary border style", func() {
		w := &bytes.Buffer{}
		palette := prism.Palette{}

		renderer, err := flow.New(palette, w)
		Expect(err).To(BeNil())

		renderer.Begin(prism.Overture{
			Kind:      prism.PrimeNavigation,
			Root:      "./src/app",
			Caption:   "files and folders",
			StartedAt: time.Date(2026, time.May, 10, 11, 31, 7, 0, time.UTC),
		})

		output := w.String()
		Expect(output).To(ContainSubstring("╭"))
		Expect(output).To(ContainSubstring("jay  ./src/app"))
		Expect(output).To(ContainSubstring("files and folders  -"))
		Expect(output).To(ContainSubstring("╰"))
	})

	It("renders final directory children without vertical continuation", func() {
		w := &bytes.Buffer{}
		palette := prism.Palette{}

		renderer, err := flow.New(palette, w)
		Expect(err).To(BeNil())

		renderer.Show(prism.Motif{Name: "src", IsDir: true, Depth: 0, VisualDepth: 0, IsLast: true})
		renderer.Show(prism.Motif{Name: "app", IsDir: true, Depth: 1, VisualDepth: 1, IsLast: false})
		renderer.Show(prism.Motif{Name: "main.go", IsDir: false, Depth: 2, VisualDepth: 2, IsLast: true})
		renderer.Show(prism.Motif{Name: "ui", IsDir: true, Depth: 1, VisualDepth: 1, IsLast: true})
		renderer.Show(prism.Motif{Name: "doc.go", IsDir: false, Depth: 2, VisualDepth: 2, IsLast: true})

		output := w.String()
		Expect(output).To(ContainSubstring("└── 📁 ui/"))
		Expect(output).To(ContainSubstring("   └── 🔖 doc.go"))
	})

	It("applies BranchStyle from theme to branch characters", func() {
		w := &bytes.Buffer{}
		palette := prism.SystemPalette()
		palette.Branch = prism.SemanticColour{ANSI16: "green"}

		renderer, err := flow.New(palette, w)
		Expect(err).To(BeNil())
		Expect(renderer).NotTo(BeNil())

		renderer.Show(prism.Motif{Name: "root", IsDir: true, Depth: 0, VisualDepth: 0, IsLast: true})
		renderer.Show(prism.Motif{Name: "child", IsDir: false, Depth: 1, VisualDepth: 1, IsLast: true})

		output := w.String()
		Expect(output).To(ContainSubstring("✻ root/\n"))
		Expect(output).To(ContainSubstring("└── 🔖 child\n"))
	})
})
