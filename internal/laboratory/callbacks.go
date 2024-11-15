package lab

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2" //nolint:revive,stylecheck // ok
	. "github.com/onsi/gomega"    //nolint:revive,stylecheck // ok

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/life"
)

func Begin(em string) life.BeginHandler {
	return func(state *life.BeginState) {
		GinkgoWriter.Printf(
			"---> %v [traverse-navigator-test:BEGIN], tree: '%v'\n", em, state.Tree,
		)
	}
}

func End(em string) life.EndHandler {
	return func(result core.TraverseResult) {
		GinkgoWriter.Printf(
			"---> %v [traverse-navigator-test:END], err: '%v'\n", em, result.Error(),
		)
	}
}

func UniversalCallback(name string) core.Client {
	return func(servant age.Servant) error {
		node := servant.Node()
		depth := node.Extension.Depth
		GinkgoWriter.Printf(
			"---> 🌊 UNIVERSAL//%v-CALLBACK: (depth:%v) '%v'\n", name, depth, node.Path,
		)
		Expect(node.Extension).NotTo(BeNil(), Reason(node.Path))

		return nil
	}
}

func DirectoriesCallback(name string) core.Client {
	return func(servant age.Servant) error {
		node := servant.Node()
		depth := node.Extension.Depth
		actualNoChildren := len(node.Children)
		GinkgoWriter.Printf(
			"---> 🔆 DIRECTORIES//CALLBACK%v: (depth:%v, children:%v) '%v'\n",
			name, depth, actualNoChildren, node.Path,
		)
		Expect(node.IsDirectory()).To(BeTrue(),
			Because(node.Path, "node expected to be directory"),
		)
		Expect(node.Extension).NotTo(BeNil(), Reason(node.Path))

		return nil
	}
}

func FilesCallback(name string) core.Client {
	return func(servant age.Servant) error {
		node := servant.Node()
		GinkgoWriter.Printf("---> 🌙 FILES//%v-CALLBACK: '%v'\n", name, node.Path)
		Expect(node.IsDirectory()).To(BeFalse(),
			Because(node.Path, "node expected to be file"),
		)
		Expect(node.Extension).NotTo(BeNil(), Reason(node.Path))

		return nil
	}
}

func DirectoriesCaseSensitiveCallback(first, second string) core.Client {
	recall := make(Recall)

	return func(servant age.Servant) error {
		node := servant.Node()
		recall[node.Path] = len(node.Children)

		GinkgoWriter.Printf("---> 🔆 CASE-SENSITIVE-CALLBACK: '%v'\n", node.Path)
		Expect(node.IsDirectory()).To(BeTrue())

		if strings.HasSuffix(node.Path, second) {
			GinkgoWriter.Printf("---> 💧 FIRST: '%v', 💧 SECOND: '%v'\n", first, second)

			paths := lo.Keys(recall)
			_, found := lo.Find(paths, func(s string) bool {
				return strings.HasSuffix(s, first)
			})

			Expect(found).To(BeTrue(), fmt.Sprintf("for node: '%v'", node.Extension.Name))
		}

		return nil
	}
}

func PanicAt(at string) core.Client {
	return func(servant core.Servant) error {
		node := servant.Node()
		depth := node.Extension.Depth
		name := node.Extension.Name

		GinkgoWriter.Printf(
			"---> 👿 PANIC-AT//%v-CALLBACK: (depth:%v) '%v'\n",
			name, depth, node.Path,
		)

		if name == at {
			panic("foo bar")
		}

		return nil
	}
}
