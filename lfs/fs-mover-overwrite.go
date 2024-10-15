package lfs

type overwriteMover struct {
	baseMover
}

func (m *overwriteMover) create() mover {
	m.actions = movers{
		{true, false, false, false}: m.moveItemWithName,         // from exists as file, to does not exist
		{true, false, true, false}:  m.moveItemWithName,         // from exists as dir, to does not exist
		{true, true, false, true}:   m.moveItemWithoutName,      // from exists as file,to exists as dir
		{true, true, true, true}:    m.moveItemWithoutNameClash, // from exists as dir, to exists as dir
		{true, true, false, false}:  m.noOp,                     // from and to refer to the same existing file
	}

	return m
}
