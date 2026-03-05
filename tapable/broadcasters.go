package tapable

import (
	"io/fs"

	"github.com/snivilised/agenor/core"
)

type (
	// SubPathBroadcaster is a function type that defines the signature for broadcasting sub-path hooks.
	SubPathBroadcaster func(def core.SubPathHook,
		provider listenerProvider[core.ChainSubPathHook],
	) core.SubPathHook
)

// GetSubPathBroadcaster creates a new SubPathHook that broadcasts to all registered listeners.
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

// SubPathAttacher creates a new SubPathHook that attaches the broadcaster to
// the provided definition and provider.
func SubPathAttacher(def core.SubPathHook,
	provider listenerProvider[core.ChainSubPathHook],
	broadcaster SubPathBroadcaster,
) core.SubPathHook {
	return func(info *core.SubPathInfo) string {
		return broadcaster(def, provider)(info)
	}
}

type (
	// ReadDirectoryBroadcaster is a function type that defines the signature for
	// broadcasting read directory hooks.
	ReadDirectoryBroadcaster func(def core.ReadDirectoryHook,
		provider listenerProvider[core.ChainReadDirectoryHook],
	) core.ReadDirectoryHook
)

// GetReadDirectoryBroadcaster creates a new ReadDirectoryHook that
// broadcasts to all registered listeners.
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

// ReadDirectoryAttacher creates a new ReadDirectoryHook that attaches the broadcaster to
// the provided definition and provider.
func ReadDirectoryAttacher(def core.ReadDirectoryHook,
	provider listenerProvider[core.ChainReadDirectoryHook],
	broadcaster ReadDirectoryBroadcaster,
) core.ReadDirectoryHook {
	return func(rsys fs.ReadDirFS, dirname string) ([]fs.DirEntry, error) {
		return broadcaster(def, provider)(rsys, dirname)
	}
}

type (
	// QueryStatusBroadcaster is a function type that defines the signature for
	// broadcasting query status hooks.
	QueryStatusBroadcaster func(def core.QueryStatusHook,
		provider listenerProvider[core.ChainQueryStatusHook],
	) core.QueryStatusHook
)

// GetQueryStatusBroadcaster creates a new QueryStatusHook that broadcasts to
// all registered listeners.
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

// QueryStatusAttacher creates a new QueryStatusHook that attaches the broadcaster to
// the provided definition and provider.
func QueryStatusAttacher(def core.QueryStatusHook,
	provider listenerProvider[core.ChainQueryStatusHook],
	broadcaster QueryStatusBroadcaster,
) core.QueryStatusHook {
	return func(qsys fs.StatFS, path string) (fs.FileInfo, error) {
		return broadcaster(def, provider)(qsys, path)
	}
}

type (
	// SortBroadcaster is a function type that defines the signature for
	// broadcasting sort hooks.
	SortBroadcaster func(def core.SortHook,
		provider listenerProvider[core.ChainSortHook],
	) core.SortHook
)

// GetSortBroadcaster creates a new SortHook that broadcasts to all
// registered listeners.
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

// SortAttacher creates a new SortHook that attaches the broadcaster to
// the provided definition and provider.
func SortAttacher(def core.SortHook,
	provider listenerProvider[core.ChainSortHook],
	broadcaster SortBroadcaster,
) core.SortHook {
	return func(entries []fs.DirEntry, custom ...any) {
		broadcaster(def, provider)(entries, custom...)
	}
}
