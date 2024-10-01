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
	stateMarshaler interface {
		Marshal(path string) error
		Unmarshal(path string) error
	}

	MarshalState struct {
		O      *pref.Options
		Active *types.ActiveState
	}

	jsonState struct {
		JO     *json.Options
		Active *types.ActiveState
	}
)

func Marshal(ms *MarshalState, path string, perm fs.FileMode,
	wfs lfs.WriteFileFS,
) (*json.Options, error) {
	jo := ToJSON(ms.O)
	state := &jsonState{
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

	if equal, err := Equals(ms.O, jo); !equal {
		return jo, err
	}

	return jo, wfs.WriteFile(path, data, perm)
}

func Unmarshal(_ *types.RestoreState, path string,
	reader lfs.ReadFileFS,
) (*MarshalState, error) {
	_ = path
	_ = reader

	return &MarshalState{}, nil
}
