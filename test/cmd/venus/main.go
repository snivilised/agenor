package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/snivilised/li18ngo"
	nef "github.com/snivilised/nefilim"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/test/hydra"
)

const (
	usage   = `Usage: go run venus \n\t-path <relative-path> \n\t[-sub <file|dir>] \n\t[-filter]`
	verbose = false
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
	subscription enums.Subscription
	filters      string
	path         string
	pattern      string
}

func main() {
	var filters multiFlag

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

	flag.Parse()

	if *path == "" {
		fmt.Println(usage)

		flag.PrintDefaults()
		os.Exit(1)
	}

	n := &navigation{
		subscription: subscribe(*sub),
		filters:      filters.String(),
		path:         *path,
		pattern:      *pattern,
	}

	if n.filters != "" {
		fmt.Printf("🥝 filters: %v\n", n.filters)
	}

	if n.pattern != "" {
		fmt.Printf("🍒 pattern: %v\n", n.pattern)
	}

	navigate(n)
}

func navigate(n *navigation) {
	if err := li18ngo.Use(); err != nil {
		fmt.Printf("%v \n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	fS := hydra.Nuxx(verbose, strings.Split(n.filters, ",")...)

	result, err := tv.Walk().Configure().Extent(tv.Prime(
		&pref.Using{
			Tree:         n.path,
			Subscription: tv.SubscribeFiles,
			Head: pref.Head{
				Handler: func(servant tv.Servant) error {
					node := servant.Node()
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
									"\t\t%v  %v\n",
									indicator, node.Extension.Name,
								)
							},
						),
					)

					return nil
				},
				GetForest: func(_ string) *core.Forest {
					return &core.Forest{
						T: fS,
						R: nef.NewTraverseABS(),
					}
				},
			},
		},
		tv.WithOnBegin(lab.Begin("🔊")),
		tv.WithOnEnd(lab.End("🏁")),
		tv.IfOptionF(n.pattern != "", func() pref.Option {
			return tv.WithFilter(&pref.FilterOptions{
				Node: &core.FilterDef{
					Type:        enums.FilterTypeGlobEx,
					Description: "as selected by user",
					Pattern:     n.pattern,
					Scope:       enums.ScopeAll,
				},
			})
		}),
		tv.WithHookReadDirectory(readEntriesHook),
	)).Navigate(ctx)

	if err != nil {
		fmt.Printf("%v \n", err)
		os.Exit(1)
	}

	fmt.Printf("===> 🍭 invoked '%v' directories, '%v' files.\n",
		result.Metrics().Count(enums.MetricNoDirectoriesInvoked),
		result.Metrics().Count(enums.MetricNoFilesInvoked),
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
		return name != ".DS_Store" && name != "."
	})

	return filtered, nil
}
