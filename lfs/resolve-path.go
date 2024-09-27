package lfs

import (
	"os"
	"path/filepath"

	"github.com/snivilised/traverse/internal/third/lo"
)

// ResolvePath performs 2 forms of path resolution. The first is resolving a
// home path reference, via the ~ character; ~ is replaced by the user's
// home path. The second resolves ./ or ../ relative path. The overrides
// do not need to be provided.
func ResolvePath(path string, mocks ...ResolveMocks) string {
	result := path

	if len(mocks) > 0 {
		m := mocks[0]
		result = lo.TernaryF(result[0] == '~',
			func() string {
				if h, err := m.HomeFunc(); err == nil {
					return filepath.Join(h, result[1:])
				}

				return path
			},
			func() string {
				if a, err := m.AbsFunc(result); err == nil {
					return a
				}

				return path
			},
		)
	} else {
		result = lo.TernaryF(result[0] == '~',
			func() string {
				if h, err := os.UserHomeDir(); err == nil {
					return filepath.Join(h, result[1:])
				}

				return path
			},
			func() string {
				if a, err := filepath.Abs(result); err == nil {
					return a
				}

				return path
			},
		)
	}

	return result
}
