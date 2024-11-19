package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/test/hanno"
	"github.com/snivilised/agenor/tfs"
	"github.com/snivilised/li18ngo"
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
	fS := hanno.Nuxx(verbose, strings.Split(n.filters, ",")...)

	result, err := age.Walk().Configure().Extent(age.Prime(
		&pref.Using{
			Tree:         n.path,
			Subscription: n.subscription,
			Head: pref.Head{
				Handler: func(servant age.Servant) error {
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
						R: tfs.New(),
					}
				},
			},
		},
		age.WithOnBegin(lab.Begin("🔊")),
		age.WithOnEnd(lab.End("🏁")),
		age.IfOptionF(n.pattern != "", func() pref.Option {
			return age.WithFilter(&pref.FilterOptions{
				Node: &core.FilterDef{
					Type:        enums.FilterTypeGlobEx,
					Description: "as selected by user",
					Pattern:     n.pattern,
					Scope:       enums.ScopeAll,
				},
			})
		}),
		age.WithHookReadDirectory(readEntriesHook),
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
