// Package main implements a multi-threaded file system navigator for
// test/demonstration purposes.
package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io/fs"
	"math/big"

	"os"
	"strings"
	"sync"
	"time"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/internal/third/lo"
	"github.com/snivilised/jaywalk/src/agenor/life"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/agenor/test/hanno"
	"github.com/snivilised/jaywalk/src/agenor/tfs"
	"github.com/snivilised/li18ngo"
)

const (
	usage = `Usage: go run nexus \n\t-path <relative-path> \n\t[-sub <file|dir>] \n\t[-filter]`
)

type multiFlag []string

func (f *multiFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *multiFlag) Set(value string) error {
	*f = append(*f, strings.Split(value, ",")...)
	return nil
}

type navigation struct {
	subscription    enums.Subscription
	filters         string
	path            string
	pattern         string
	handler         core.Client
	isWalk, isPrime bool
	options         []pref.Option
}

func main() {
	if err := li18ngo.Use(); err != nil {
		fmt.Printf("%v \n", err)
		os.Exit(1)
	}

	var (
		filters multiFlag
	)

	path := flag.String("path", "",
		"path to navigate from",
	)

	sub := flag.String("sub", "universal",
		"subscription type [file|dir] (defaults to universal)",
	)

	pattern := flag.String("pattern", "",
		"glob-ex filter [parent pattern|file pattern], eg *|flac",
	)

	flag.Var(&filters,
		"filter",
		"Specify filter(s) (csv)",
	)

	now := flag.Int("now", 0, "no of workers")
	resume := flag.Bool("resume", false, "resume navigation (not supported yet)")
	_ = resume

	seconds := flag.Int("delay", 1, "no of seconds to represent randomised delay interval")

	flag.Parse()

	if *path == "" {
		fmt.Println(usage)

		flag.PrintDefaults()
		os.Exit(1)
	}

	size := uint(*now) //nolint:gosec // ok
	options := []pref.Option{
		agenor.WithOnBegin(func(state *life.BeginState) {
			fmt.Printf(
				"---> 🛡️ [nexus-traverse-navigator:BEGIN], tree: '%v'\n", state.Tree,
			)
		}),
		agenor.WithOnEnd(func(_ core.TraverseResult) {
			fmt.Println(
				"---> 🏁 [nexus-traverse-navigator:END]",
			)
		}),
		agenor.IfOptionF(*pattern != "", func() pref.Option {
			return agenor.WithFilter(&pref.FilterOptions{
				Node: &core.FilterDef{
					Type:        enums.FilterTypeGlobEx,
					Description: "as selected by user",
					Pattern:     *pattern,
					Scope:       enums.ScopeAll,
				},
			})
		}),
		agenor.IfOptionF(size > 0, func() pref.Option {
			return agenor.WithNoW(size)
		}),
		agenor.WithHookReadDirectory(readEntriesHook),
	}

	n := &navigation{
		subscription: subscribe(*sub),
		filters:      filters.String(),
		path:         *path,
		pattern:      *pattern,
		handler:      lo.Ternary(size == 0, sequentialHandler, fileWorker(*seconds)),
		isWalk:       size == 0,
		isPrime:      true, // !resume
		options:      options,
	}

	if n.filters != "" {
		fmt.Printf("🥝 filters: %v\n", n.filters)
	}

	if n.pattern != "" {
		fmt.Printf("🍒 pattern: %v\n", n.pattern)
	}

	if size == 0 {
		fmt.Println("... crawling like a tortoise 🐢")
	} else {
		fmt.Printf("!!! running like a hare 🐰, with %v workers\n", size)
	}

	navigate(n)
}

func navigate(n *navigation) {
	// custom forest only required because we're using a virtual in memory tree
	// instead of the local fs.
	//
	forest := func(_ string) *core.Forest {
		const quiet = false

		return &core.Forest{
			T: hanno.Nuxx(quiet),
			R: tfs.New(),
		}
	}

	facade := lo.TernaryF(n.isPrime,
		func() pref.Facade {
			return &pref.Using{
				Subscription: n.subscription,
				Head: pref.Head{
					Handler:   n.handler,
					GetForest: forest,
				},
				Tree: n.path,
			}
		},
		func() pref.Facade {
			return &pref.Relic{
				Head: pref.Head{
					Handler:   n.handler,
					GetForest: forest,
				},
				From:     "path-to-json-file",
				Strategy: agenor.ResumeStrategyFastward,
			}
		},
	)

	var (
		wg sync.WaitGroup
	)

	result, err := agenor.Hydra(
		n.isWalk,
		n.isPrime,
		&wg,
	)(facade, n.options...).Navigate(context.Background())

	wg.Wait()

	if err != nil {
		fmt.Printf("%v \n", err)
		os.Exit(1)
	}

	fmt.Printf("===> 🍭 invoked '%v' directories, '%v' files.\n",
		result.Metrics().Count(enums.MetricNoDirectoriesInvoked),
		result.Metrics().Count(enums.MetricNoFilesInvoked),
	)
}

func sequentialHandler(servant core.Servant) error {
	display(servant.Node(), "")

	return nil
}

func fileWorker(seconds int) core.Client {
	return func(servant core.Servant) error {
		node := servant.Node()
		actual := lo.Ternary(node.IsDirectory(), 0, random(int64(seconds)))
		interval := time.Second * time.Duration(actual)

		display(node, lo.Ternary(interval == 0,
			"",
			fmt.Sprintf(" 💤 (%v seconds)", actual)),
		)

		<-time.After(interval)

		return nil
	}
}

func display(node *core.Node, snooze string) {
	indicator := lo.Ternary(node.IsDirectory(), "📂", "🏷️")

	fmt.Print(
		lo.TernaryF(node.IsDirectory(),
			func() string {
				return fmt.Sprintf(
					"\t%v  %v\n",
					indicator, node.Path,
				)
			},
			func() string {
				return fmt.Sprintf(
					"\t\t%v  %v%v\n",
					indicator, node.Extension.Name, snooze,
				)
			},
		),
	)
}

func subscribe(sub string) enums.Subscription {
	subscription := enums.SubscribeUniversal

	switch sub {
	case "file", "f":
		subscription = enums.SubscribeFiles
	case "dir", "d":
		subscription = enums.SubscribeDirectories
	}

	return subscription
}

func readEntriesHook(sys fs.ReadDirFS,
	dirname string,
) ([]fs.DirEntry, error) {
	contents, err := fs.ReadDir(sys, dirname)
	if err != nil {
		return nil, err
	}

	filtered := lo.Filter(contents, func(item fs.DirEntry, _ int) bool {
		name := item.Name()
		return name != "."
	})

	return filtered, nil
}

func random(maximum int64) int {
	m := big.NewInt(maximum)
	result, _ := rand.Int(rand.Reader, m)

	return int(result.Int64())
}
