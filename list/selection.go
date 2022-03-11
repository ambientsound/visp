package list

// ManuallySelected returns true if the given song index is selected through manual selection.
func (s *Base) ManuallySelected(i int) bool {
	_, ok := s.selection[i]
	return ok
}

// VisuallySelected returns true if the given song index is selected through visual selection.
func (s *Base) VisuallySelected(i int) bool {
	return s.visualSelection[0] <= i && i <= s.visualSelection[1]
}

// Selected returns true if the given song index is selected, either through
// visual selection or manual selection. If the song is doubly selected, the
// selection is inversed.
func (s *Base) Selected(i int) bool {
	a := s.ManuallySelected(i)
	b := s.VisuallySelected(i)
	return (a || b) && a != b
}

// SelectionIndices returns a slice of ints holding the position of each
// element in the current selection. If no elements are selected, the cursor
// position is returned.
func (s *Base) SelectionIndices() []int {
	selection := make([]int, 0, s.Len())
	max := s.Len()
	for i := 0; i < max; i++ {
		if s.Selected(i) {
			selection = append(selection, i)
		}
	}
	if len(selection) == 0 && s.Len() > 0 {
		selection = append(selection, s.Cursor())
	}
	return selection
}

// Selection returns selected rows as a new list.
func (s *Base) Selection() List {
	indices := s.SelectionIndices()
	result := New()
	for _, i := range indices {
		// fixme: copy?
		result.Add(s.rows[i])
	}
	return result
}

// SetSelection sets the selected status of a single song.
func (s *Base) SetSelected(i int, selected bool) {
	var x struct{}
	_, ok := s.selection[i]
	if ok == selected {
		return
	}
	if selected {
		s.selection[i] = x
	} else {
		delete(s.selection, i)
	}
}

// CommitVisualSelection converts the visual selection to manual selection.
func (s *Base) CommitVisualSelection() {
	if !s.HasVisualSelection() {
		return
	}
	for key := s.visualSelection[0]; key <= s.visualSelection[1]; key++ {
		selected := s.Selected(key)
		s.SetSelected(key, selected)
	}
}

// ClearSelection removes all selection.
func (s *Base) ClearSelection() {
	s.selection = make(map[int]struct{}, 0)
	s.visualSelection = [3]int{-1, -1, -1}
}

// validateVisualSelection makes sure the visual selection stays in range of
// the songlist size.
func (s *Base) validateVisualSelection(ymin, ymax, ystart int) (int, int, int) {
	if s.Len() == 0 || ymin < 0 || ymax < 0 || !s.InRange(ystart) {
		return -1, -1, -1
	}
	if !s.InRange(ymin) {
		ymin = 0
	}
	if !s.InRange(ymax) {
		ymax = s.Len() - 1
	}
	return ymin, ymax, ystart
}

// VisualSelection returns the min, max, and start position of visual select.
func (s *Base) VisualSelection() (int, int, int) {
	return s.visualSelection[0], s.visualSelection[1], s.visualSelection[2]
}

// SetVisualSelection sets the range of the visual selection. Use negative
// integers to un-select all visually selected songs.
func (s *Base) SetVisualSelection(ymin, ymax, ystart int) {
	s.visualSelection[0], s.visualSelection[1], s.visualSelection[2] = s.validateVisualSelection(ymin, ymax, ystart)
}

// HasVisualSelection returns true if the songlist is in visual selection mode.
func (s *Base) HasVisualSelection() bool {
	return s.visualSelection[0] >= 0 && s.visualSelection[1] >= 0
}

// EnableVisualSelection sets start and stop of the visual selection to the
// cursor position.
func (s *Base) EnableVisualSelection() {
	cursor := s.Cursor()
	s.SetVisualSelection(cursor, cursor, cursor)
}

// DisableVisualSelection disables visual selection.
func (s *Base) DisableVisualSelection() {
	s.SetVisualSelection(-1, -1, -1)
}

// ToggleVisualSelection toggles visual selection on and off.
func (s *Base) ToggleVisualSelection() {
	if !s.HasVisualSelection() {
		s.EnableVisualSelection()
	} else {
		s.DisableVisualSelection()
	}
}

// expandVisualSelection sets the visual selection boundaries from where it
// started to the current cursor position.
func (s *Base) expandVisualSelection() {
	if !s.HasVisualSelection() {
		return
	}
	ymin, ymax, ystart := s.VisualSelection()
	switch {
	case s.Cursor() < ystart:
		ymin, ymax = s.Cursor(), ystart
	case s.Cursor() > ystart:
		ymin, ymax = ystart, s.Cursor()
	default:
		ymin, ymax = ystart, ystart
	}
	s.SetVisualSelection(ymin, ymax, ystart)
}
