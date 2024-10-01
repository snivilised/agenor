package json

// TODO: can't have package name that is json as that clashes
// with the one in standard library at encoding/json, so need
// to rename; perhaps to js.

type (
	// ConcurrencyOptions specifies options used for current traversal sessions
	ConcurrencyOptions struct {
		// NoW specifies the number of go-routines to use in the worker
		// pool used for concurrent traversal sessions requested by using
		// the Run function.
		NoW uint `json:"no-of-workers"`
	}
)
