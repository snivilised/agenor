package controller_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/app/controller"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// makeNode builds a minimal *core.Node with the given absolute path.
func makeNode(path string) *core.Node {
	return &core.Node{Path: path}
}

// ---------------------------------------------------------------------------
// Specs
// ---------------------------------------------------------------------------

var _ = Describe("expand", func() {

	// ------------------------------------------------------------------
	// Basic placeholder substitution - no breach possible
	// ------------------------------------------------------------------

	Context("given a node three levels below root", func() {
		var (
			root string
			node *core.Node
		)

		BeforeEach(func() {
			root = filepath.Join("/", "home", "user", "videos")
			node = makeNode(filepath.Join(root, "holiday", "day1", "clip.mp4"))
		})

		It("expands {{.path}} to the node's absolute path", func() {
			result := controller.Expand("echo {{.path}}", root, node)

			Expect(result.Skipped).To(BeFalse())
			Expect(result.Cmd).To(ContainSubstring(node.Path))
		})

		It("expands {{.name}} to the filename including extension", func() {
			result := controller.Expand("echo {{.name}}", root, node)

			Expect(result.Skipped).To(BeFalse())
			Expect(result.Cmd).To(ContainSubstring("clip.mp4"))
		})

		It("expands {{.stem}} to the filename without extension", func() {
			result := controller.Expand("echo {{.stem}}", root, node)

			Expect(result.Skipped).To(BeFalse())
			Expect(result.Cmd).To(ContainSubstring("clip"))
			Expect(result.Cmd).NotTo(ContainSubstring(".mp4"))
		})

		It("expands {{.ext}} to the extension including the dot", func() {
			result := controller.Expand("echo {{.ext}}", root, node)

			Expect(result.Skipped).To(BeFalse())
			Expect(result.Cmd).To(ContainSubstring(".mp4"))
		})

		It("expands {{.parent}} to the immediate parent directory", func() {
			result := controller.Expand("echo {{.parent}}", root, node)

			Expect(result.Skipped).To(BeFalse())
			Expect(result.Cmd).To(ContainSubstring(
				filepath.Join(root, "holiday", "day1"),
			))
		})

		It("expands {{.grand}} to the grandparent directory", func() {
			result := controller.Expand("echo {{.grand}}", root, node)

			Expect(result.Skipped).To(BeFalse())
			Expect(result.Cmd).To(ContainSubstring(
				filepath.Join(root, "holiday"),
			))
		})

		It("expands {{.great}} to the great-grandparent directory", func() {
			result := controller.Expand("echo {{.great}}", root, node)

			Expect(result.Skipped).To(BeFalse())
			Expect(result.Cmd).To(ContainSubstring(root))
		})

		It("expands {{.root}} to the traversal root", func() {
			result := controller.Expand("echo {{.root}}", root, node)

			Expect(result.Skipped).To(BeFalse())
			Expect(result.Cmd).To(ContainSubstring(root))
		})

		It("expands multiple placeholders in a single cmd string", func() {
			cmd := "ffmpeg -i {{.path}} -q:v 2 {{.parent}}/{{.stem}}.mp4"
			result := controller.Expand(cmd, root, node)

			Expect(result.Skipped).To(BeFalse())
			Expect(result.Cmd).To(ContainSubstring(node.Path))
			Expect(result.Cmd).To(ContainSubstring("clip.mp4"))
		})
	})

	// ------------------------------------------------------------------
	// Breach detection - ancestor climbs above root
	// ------------------------------------------------------------------

	Context("given a node one level below root", func() {
		var (
			root string
			node *core.Node
		)

		BeforeEach(func() {
			root = filepath.Join("/", "home", "user", "videos")
			node = makeNode(filepath.Join(root, "clip.mp4"))
		})

		It("skips when {{.grand}} would breach root", func() {
			result := controller.Expand("echo {{.grand}}", root, node)

			Expect(result.Skipped).To(BeTrue())
			Expect(result.Placeholder).To(Equal("{{.grand}}"))
		})

		It("skips when {{.great}} would breach root", func() {
			result := controller.Expand("echo {{.great}}", root, node)

			Expect(result.Skipped).To(BeTrue())
			Expect(result.Placeholder).To(Equal("{{.great}}"))
		})

		It("does not skip when only {{.parent}} is used (parent == root)", func() {
			// parent of a node directly under root is root itself - not a breach.
			result := controller.Expand("echo {{.parent}}", root, node)

			Expect(result.Skipped).To(BeFalse())
		})
	})

	Context("given a node two levels below root", func() {
		var (
			root string
			node *core.Node
		)

		BeforeEach(func() {
			root = filepath.Join("/", "home", "user", "videos")
			node = makeNode(filepath.Join(root, "holiday", "clip.mp4"))
		})

		It("does not skip when {{.grand}} resolves exactly to root", func() {
			result := controller.Expand("echo {{.grand}}", root, node)

			Expect(result.Skipped).To(BeFalse())
		})

		It("skips when {{.great}} would breach root", func() {
			result := controller.Expand("echo {{.great}}", root, node)

			Expect(result.Skipped).To(BeTrue())
			Expect(result.Placeholder).To(Equal("{{.great}}"))
		})
	})

	// ------------------------------------------------------------------
	// Breach reports the first offending placeholder
	// ------------------------------------------------------------------

	Context("when multiple placeholders would breach root", func() {
		It("reports the first offending placeholder in check order", func() {
			root := filepath.Join("/", "home", "user", "videos")
			// node is directly under root - both {{.grand}} and {{.great}} breach
			node := makeNode(filepath.Join(root, "clip.mp4"))

			// cmd references grand first, then great
			result := controller.Expand("echo {{.grand}} {{.great}}", root, node)

			Expect(result.Skipped).To(BeTrue())
			Expect(result.Placeholder).To(Equal("{{.grand}}"))
		})
	})
})
