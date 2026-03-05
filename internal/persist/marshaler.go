package persist

import (
	ejson "encoding/json"
	"io/fs"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/internal/enclave"
	json "github.com/snivilised/agenor/internal/opts/jason"
	"github.com/snivilised/agenor/pref"
	nef "github.com/snivilised/nefilim"
)

type (
	// TamperFunc provides a way for unit tests to modify the JSON before
	// it is un-marshaled. The unit tests marshal a default JSON object
	// instance, so a TamperFunc is used to allow modification of that
	// default. Typically a single test will focus on a single field,
	// so that the TamperFunc is expected to only update 1 of the members at a
	// time.
	TamperFunc func(result *MarshalResult)

	// MarshalRequest is the input for the Marshal function
	MarshalRequest struct {
		// Active is the active state to be marshaled
		Active *core.ActiveState

		// O is the options to be marshaled
		O *pref.Options

		// Path is the file path to write the marshaled data to
		Path string

		// Perm is the file permissions to use when writing the marshaled data
		Perm fs.FileMode

		// FS is the file system to use for writing the marshaled data
		FS nef.WriteFileFS
	}

	// MarshalResult is the output of the Marshal function
	MarshalResult struct {
		// Active is the active state that was marshaled
		Active *core.ActiveState

		//
		JO *json.Options
	}

	// UnmarshalRequest is the input for the Unmarshal function
	UnmarshalRequest struct {
		// Restore is the restore state containing the file path and file
		// system to read the marshaled data from
		Restore *enclave.RestoreState
	}

	// UnmarshalResult is the output of the Unmarshal function
	UnmarshalResult struct {
		// Active is the active state that was unmarshaled
		Active *core.ActiveState

		// JO is the JSON options that were unmarshaled
		JO *json.Options

		// O is the options that were unmarshaled
		O *pref.Options
	}

	// Comparison is a struct used to compare the JSON options and the pref options
	Comparison struct {
		// JO is the JSON options to compare
		JO *json.Options

		// O is the pref options to compare
		O *pref.Options
	}
)

// Marshal marshals the request
func Marshal(request *MarshalRequest) (*MarshalResult, error) {
	jo := ToJSON(request.O)
	result := &MarshalResult{
		JO:     jo,
		Active: request.Active.Clone(),
	}

	data, err := ejson.MarshalIndent(
		result,
		JSONMarshalNoPrefix, JSONMarshal2SpacesIndent,
	)
	if err != nil {
		return nil, err
	}

	if err := (&Comparison{
		O:  request.O,
		JO: jo,
	}).Equals(); err != nil {
		return result, err
	}

	return result, request.FS.WriteFile(request.Path, data, request.Perm)
}

// Unmarshal unpacks the incoming request
func Unmarshal(request *UnmarshalRequest,
	tampers ...TamperFunc,
) (*UnmarshalResult, error) {
	bytes, err := request.Restore.FS.ReadFile(request.Restore.Path)
	if err != nil {
		return nil, err
	}

	var (
		mr MarshalResult
	)

	if err := ejson.Unmarshal(bytes, &mr); err != nil {
		return nil, err
	}

	for _, fn := range tampers {
		fn(&mr)
	}

	result := UnmarshalResult{
		O:      FromJSON(mr.JO),
		Active: mr.Active,
		JO:     mr.JO,
	}

	return &result, (&Comparison{
		O:  result.O,
		JO: result.JO,
	}).Equals()
}
