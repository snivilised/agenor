package kernel_test

import (
	"fmt"
	"io/fs"
	"strings"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/samber/lo"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/helpers"
)

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
	message       string
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

func begin(em string) cycle.BeginHandler {
	return func(root string) {
		GinkgoWriter.Printf(
			"---> %v [traverse-navigator-test:BEGIN], root: '%v'\n", em, root,
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
		Expect(node.IsFolder()).To(BeTrue())
		Expect(node.Extension).NotTo(BeNil(), helpers.Reason(node.Path))

		return nil
	}
}

func filesCallback(name string) core.Client {
	return func(node *core.Node) error {
		GinkgoWriter.Printf("---> 🌙 FILES//%v-CALLBACK: '%v'\n", name, node.Path)
		Expect(node.IsFolder()).To(BeFalse())
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
