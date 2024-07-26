package tapable

import (
	"io/fs"

	"github.com/snivilised/traverse/core"
)

type (
	SubPathBroadcaster func(def core.SubPathHook,
		provider listenerProvider[core.ChainSubPathHook],
	) core.SubPathHook
)

func GetSubPathBroadcaster(def core.SubPathHook,
	provider listenerProvider[core.ChainSubPathHook],
) core.SubPathHook {
	return func(info *core.SubPathInfo) string {
		result := def(info)

		for _, listener := range provider.get() {
			result = listener(result, info)
		}

		return result
	}
}

func SubPathAttacher(def core.SubPathHook,
	provider listenerProvider[core.ChainSubPathHook],
	broadcaster SubPathBroadcaster,
) core.SubPathHook {
	return func(info *core.SubPathInfo) string {
		return broadcaster(def, provider)(info)
	}
}

type (
	ReadDirectoryBroadcaster func(def core.ReadDirectoryHook,
		provider listenerProvider[core.ChainReadDirectoryHook],
	) core.ReadDirectoryHook
)

func GetReadDirectoryBroadcaster(def core.ReadDirectoryHook,
	provider listenerProvider[core.ChainReadDirectoryHook],
) core.ReadDirectoryHook {
	return func(rsys fs.ReadDirFS, dirname string) ([]fs.DirEntry, error) {
		result, err := def(rsys, dirname)

		for _, listener := range provider.get() {
			result, err = listener(result, err, rsys, dirname)
		}

		return result, err
	}
}

func ReadDirectoryAttacher(def core.ReadDirectoryHook,
	provider listenerProvider[core.ChainReadDirectoryHook],
	broadcaster ReadDirectoryBroadcaster,
) core.ReadDirectoryHook {
	return func(rsys fs.ReadDirFS, dirname string) ([]fs.DirEntry, error) {
		return broadcaster(def, provider)(rsys, dirname)
	}
}

type (
	QueryStatusBroadcaster func(def core.QueryStatusHook,
		provider listenerProvider[core.ChainQueryStatusHook],
	) core.QueryStatusHook
)

func GetQueryStatusBroadcaster(def core.QueryStatusHook,
	provider listenerProvider[core.ChainQueryStatusHook],
) core.QueryStatusHook {
	return func(qsys fs.StatFS, path string) (fs.FileInfo, error) {
		result, err := def(qsys, path)

		for _, listener := range provider.get() {
			result, err = listener(result, err, qsys, path)
		}

		return result, err
	}
}

func QueryStatusAttacher(def core.QueryStatusHook,
	provider listenerProvider[core.ChainQueryStatusHook],
	broadcaster QueryStatusBroadcaster,
) core.QueryStatusHook {
	return func(qsys fs.StatFS, path string) (fs.FileInfo, error) {
		return broadcaster(def, provider)(qsys, path)
	}
}

type (
	SortBroadcaster func(def core.SortHook,
		provider listenerProvider[core.ChainSortHook],
	) core.SortHook
)

func GetSortBroadcaster(def core.SortHook,
	provider listenerProvider[core.ChainSortHook],
) core.SortHook {
	return func(entries []fs.DirEntry, custom ...any) {
		def(entries, custom...)

		for _, listener := range provider.get() {
			listener(entries, custom...)
		}
	}
}

func SortAttacher(def core.SortHook,
	provider listenerProvider[core.ChainSortHook],
	broadcaster SortBroadcaster,
) core.SortHook {
	return func(entries []fs.DirEntry, custom ...any) {
		broadcaster(def, provider)(entries, custom...)
	}
}
