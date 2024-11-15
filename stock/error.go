package stock

import (
	"errors"
	"io/fs"
)

// 📦 pkg: stock, contain functions that can be considered addenda
// to the std golang library. Any functions/type created here should
// be self contained; no dependencies should be created towards any
// other package inside agenor.

// IsBenignError enables the distinction between a genuine err and
// a synthetic file system error, in this context described as being
// either fs.SkipDir or fs.SkipAll.
func IsBenignError(err error) bool {
	if err == nil {
		return true
	}

	return errors.Is(err, fs.SkipDir) || errors.Is(err, fs.SkipAll)
}
