package prism_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/prism"
)

var _ = Describe("StreamRenderer", func() {
	It("renders tree branches with default icons", func() {
		w := &bytes.Buffer{}
		palette := prism.Palette{}

		renderer, err := prism.New(prism.StreamView, palette, w)
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

		renderer, err := prism.New(prism.StreamView, palette, w, prism.WithIcons(map[string]string{
			prism.TreeIconRoot:           "R",
			prism.TreeIconDirectory:      "D",
			prism.TreeIconFile:           "F",
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

	It("returns a renderer when options are provided", func() {
		w := &bytes.Buffer{}
		palette := prism.Palette{}

		renderer, err := prism.New(prism.StreamView, palette, w, prism.WithIcons(nil))
		Expect(err).To(BeNil())
		Expect(renderer).NotTo(BeNil())
		Expect(renderer).To(Not(BeNil()))
		// exercise a simple event to ensure the renderer is operational.
		renderer.Show(prism.Motif{Name: "test", Depth: 0, VisualDepth: 0, IsDir: true, IsLast: true})
		Expect(w.String()).To(ContainSubstring("✻ test/"))
	})

	It("renders final directory children without vertical continuation", func() {
		w := &bytes.Buffer{}
		palette := prism.Palette{}

		renderer, err := prism.New(prism.StreamView, palette, w)
		Expect(err).To(BeNil())

		// Root directory
		renderer.Show(prism.Motif{
			Name:        "src",
			IsDir:       true,
			Depth:       0,
			VisualDepth: 0,
			IsLast:      true,
		})

		// First child directory (not last)
		renderer.Show(prism.Motif{
			Name:        "app",
			IsDir:       true,
			Depth:       1,
			VisualDepth: 1,
			IsLast:      false,
		})

		// File under first child
		renderer.Show(prism.Motif{
			Name:        "main.go",
			IsDir:       false,
			Depth:       2,
			VisualDepth: 2,
			IsLast:      true,
		})

		// Final child directory (last)
		renderer.Show(prism.Motif{
			Name:        "ui",
			IsDir:       true,
			Depth:       1,
			VisualDepth: 1,
			IsLast:      true,
		})

		// File under final child
		renderer.Show(prism.Motif{
			Name:        "doc.go",
			IsDir:       false,
			Depth:       2,
			VisualDepth: 2,
			IsLast:      true,
		})

		output := w.String()
		// The final directory should show with branch-last
		Expect(output).To(ContainSubstring("└── 📁 ui/"))
		// Its children should not have vertical continuation (no │ prefix)
		Expect(output).To(ContainSubstring("   └── 🔖 doc.go"))
	})

	It("applies BranchStyle from theme to branch characters", func() {
		w := &bytes.Buffer{}
		// Create a palette with explicit branch color
		palette := prism.SystemPalette()
		palette.Branch = prism.SemanticColour{ANSI16: "green"}

		renderer, err := prism.New(prism.StreamView, palette, w)
		Expect(err).To(BeNil())
		Expect(renderer).NotTo(BeNil())

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

		output := w.String()
		Expect(output).To(ContainSubstring("✻ root/\n"))
		Expect(output).To(ContainSubstring("└── 🔖 child\n"))
	})
})
