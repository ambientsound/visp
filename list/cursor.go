package list

type Cursor interface {
	Cursor() int
	CursorRow() Row
	MoveCursor(int)
	SetCursor(int)
	SetCursorByID(string) error
	ValidateCursor(int, int)
}

// MoveCursor moves the cursor by the specified offset.
func (s *Base) MoveCursor(i int) {
	s.SetCursor(s.Cursor() + i)
}

// SetCursor sets the cursor to an absolute position.
func (s *Base) SetCursor(i int) {
	s.cursor = i
	s.ValidateCursor(0, s.Len()-1)
	s.expandVisualSelection()
}

// SetCursor moves the cursor to the row with the specified ID.
func (s *Base) SetCursorByID(id string) error {
	rown, err := s.RowNum(id)
	if err != nil {
		return err
	}
	s.SetCursor(rown)
	return nil
}

// Cursor returns the cursor position.
func (s *Base) Cursor() int {
	return s.cursor
}

func (s *Base) CursorRow() Row {
	return s.Row(s.cursor)
}

// ValidateCursor makes sure the cursor is within minimum and maximum boundaries.
func (s *Base) ValidateCursor(ymin, ymax int) {
	if s.Cursor() < ymin {
		s.cursor = ymin
	}
	if s.Cursor() > ymax {
		s.cursor = ymax
	}
}
