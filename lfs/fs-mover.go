package lfs

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/locale"
)

const (
	moveOpName = "Move"
)

type (
	mover interface {
		create() mover
		move(from, to string) error
	}

	moveFunc func(from, to string) error

	moverBitmask struct {
		fromExists bool
		toExists   bool
		fromIsDir  bool
		toIsDir    bool
	}

	movers map[moverBitmask]moveFunc

	baseMover struct {
		root    string
		fS      UniversalFS
		actions movers
	}
)

func (m *baseMover) move(from, to string) error {
	mask := m.query(from, to)

	if action, exists := m.actions[mask]; exists {
		if err := action(from, to); err != nil {
			// we can't mask the original message
			//
			// errors.Wrap()
			// return locale.NewInvalidBinaryFsOpError(moveOpName, from, to)
			return err
		}
		return nil
	}

	return locale.NewInvalidBinaryFsOpError(moveOpName, from, to)
}

func (m *baseMover) query(from, to string) moverBitmask {
	fromExists, fromIsDir := m.peek(from)
	toExists, toIsDir := m.peek(to)

	return moverBitmask{
		fromExists: fromExists,
		toExists:   toExists,
		fromIsDir:  fromIsDir,
		toIsDir:    toIsDir,
	}
}

func (m *baseMover) peek(name string) (exists, isDir bool) {
	if m.fS.DirectoryExists(name) {
		return true, true
	}

	if m.fS.FileExists(name) {
		return true, false
	}

	return false, false
}

func (m *baseMover) moveItemWithName(from, to string) error {
	// 'to' includes the file name eg:
	// from/file.txt => to/file.txt
	//
	if filepath.Dir(from) == filepath.Dir(to) {
		return locale.NewRejectSameDirMoveError(moveOpName, from, to)
	}

	return os.Rename(
		filepath.Join(m.root, from),
		filepath.Join(m.root, to),
	)
}

func (m *baseMover) moveItemWithoutName(from, to string) error {
	// 'to' does not include the file name, so it has to be appended, eg:
	// from/file.txt => to/
	//
	return os.Rename(
		filepath.Join(m.root, from),
		filepath.Join(m.root, to, filepath.Base(from)),
	)
}

func (m *baseMover) moveItemWithoutNameClash(from, to string) error {
	fromBase := filepath.Base(from)
	toBase := filepath.Base(to)

	if fromBase == toBase {
		// If there were a merge facility, this is where we would implement this,
		// ie merge the from directory with to, instead of returning an error.
		//
		return locale.NewRejectSameDirMoveError(moveOpName, from, to)
	}

	return m.moveItemWithoutName(from, to)
}

func (m *baseMover) noOp(_, _ string) error {
	return nil
}

type lazyMover struct {
	once  sync.Once
	mover mover
}

func (l *lazyMover) instance(root string, overwrite bool, fS UniversalFS) mover {
	l.once.Do(func() {
		l.mover = l.create(root, overwrite, fS)
	})

	return l.mover
}

func (l *lazyMover) create(root string, overwrite bool, fS UniversalFS) mover {
	return lo.TernaryF(overwrite,
		func() mover {
			return &overwriteMover{
				baseMover: baseMover{
					root: root,
					fS:   fS,
				},
			}
		},
		func() mover {
			return &tentativeMover{
				baseMover: baseMover{
					root: root,
					fS:   fS,
				},
			}
		},
	).create()
}
