package persist

import (
	ejson "encoding/json"
	"io/fs"

	"github.com/snivilised/traverse/internal/opts/json"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/lfs"
	"github.com/snivilised/traverse/pref"
)

type (
	// TamperFunc provides a way for unit tests to modify the JSON before
	// it is un-marshaled. The unit tests marshal a default JSON object
	// instance, so a TamperFunc is used to allow modification of that
	// default. Typically a single test will focus on a single field,
	// so that the TamperFunc is expected to only update 1 of the members at a
	// time.
	TamperFunc func(jo *json.Options)

	MarshalState struct {
		O      *pref.Options
		Active *types.ActiveState
		JO     *json.Options
		Path   string
		Perm   fs.FileMode
		FS     lfs.WriteFileFS
	}

	JSONState struct {
		JO     *json.Options
		Active *types.ActiveState
	}
)

func Marshal(ms *MarshalState) (*json.Options, error) {
	jo := ToJSON(ms.O)
	state := &JSONState{
		JO:     jo,
		Active: ms.Active,
	}

	data, err := ejson.MarshalIndent(
		state,
		JSONMarshalNoPrefix, JSONMarshal2SpacesIndent,
	)

	if err != nil {
		return nil, err
	}

	if err := Equals(ms.O, jo); err != nil {
		return jo, err
	}

	return jo, ms.FS.WriteFile(ms.Path, data, ms.Perm)
}

func Unmarshal(rs *types.RestoreState, tampers ...TamperFunc) (*MarshalState, error) {
	bytes, err := rs.FS.ReadFile(rs.Path)

	if err != nil {
		return nil, err
	}

	var (
		js JSONState
	)

	if err := ejson.Unmarshal(bytes, &js); err != nil {
		return nil, err
	}

	for _, fn := range tampers {
		fn(js.JO)
	}

	ms := MarshalState{
		O:      FromJSON(js.JO),
		Active: js.Active,
		JO:     js.JO,
	}

	return &ms, Equals(ms.O, js.JO)
}
