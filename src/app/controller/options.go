package controller

import (
	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/app/report"
	"github.com/snivilised/jaywalk/src/locale"
)

type FilterIntent struct {
	FilesExGlob      string
	FilesRegEx       string
	DirectoriesGlob  string
	DirectoriesRegEx string
}

type TraversalSettingsIntent struct {
	NoRecurse     bool
	Depth         uint
	IsSampling    bool
	NoFiles       uint
	NoDirectories uint
	Filter        FilterIntent
}

func BuildTraversalSettings(intent TraversalSettingsIntent, ui report.Presenter) []pref.Option {
	var opts []pref.Option

	if intent.NoRecurse {
		opts = append(opts, agenor.WithNoRecurse())
	}

	if intent.Depth > 0 {
		opts = append(opts, agenor.WithDepth(intent.Depth))
	}

	if intent.IsSampling {
		opts = append(opts, agenor.WithSamplingOptions(&pref.SamplingOptions{
			NoOf: pref.EntryQuantities{
				Files:       intent.NoFiles,
				Directories: intent.NoDirectories,
			},
		}))
	}

	if filterOption, ok := TranslateFilterIntent(intent.Filter); ok {
		opts = append(opts, filterOption)
	}

	opts = append(opts, pref.WithTraversalConfigurer(ui))

	return opts
}

func TranslateFilterIntent(intent FilterIntent) (pref.Option, bool) {
	hasFileGlob := intent.FilesExGlob != ""
	hasFileRegex := intent.FilesRegEx != ""
	hasFolderGlob := intent.DirectoriesGlob != ""
	hasFolderRegex := intent.DirectoriesRegEx != ""

	if !hasFileGlob && !hasFileRegex && !hasFolderGlob && !hasFolderRegex {
		return nil, false
	}

	fileDef := core.FilterDef{}
	if hasFileGlob {
		fileDef.Type = enums.FilterTypeGlob
		fileDef.Pattern = intent.FilesExGlob
	} else if hasFileRegex {
		fileDef.Type = enums.FilterTypeRegex
		fileDef.Pattern = intent.FilesRegEx
	}

	dirDef := core.FilterDef{}
	if hasFolderGlob {
		dirDef.Type = enums.FilterTypeGlob
		dirDef.Pattern = intent.DirectoriesGlob
	} else if hasFolderRegex {
		dirDef.Type = enums.FilterTypeRegex
		dirDef.Pattern = intent.DirectoriesRegEx
	}

	return pref.WithFilter(&pref.FilterOptions{
		Node: &core.FilterDef{
			Poly: &core.PolyFilterDef{
				File:      fileDef,
				Directory: dirDef,
			},
		},
	}), true
}

func ResolveSubscription(flag string) (enums.Subscription, error) {
	switch flag {
	case "files", "":
		return enums.SubscribeFiles, nil
	case "dirs":
		return enums.SubscribeDirectories, nil
	case "all":
		return enums.SubscribeUniversal, nil
	default:
		return 0, locale.ErrInvalidSubscription
	}
}
