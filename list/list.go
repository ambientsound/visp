package list

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type Selectable interface {
	ClearSelection()
	CommitVisualSelection()
	DisableVisualSelection()
	EnableVisualSelection()
	HasVisualSelection() bool
	Selected(int) bool
	SelectionIndices() []int
	SetSelected(int, bool)
	SetVisualSelection(int, int, int)
	ToggleVisualSelection()
}

type Metadata interface {
	ColumnNames() []string
	Columns([]string) []Column
	ID() string
	Name() string
	SetColumnNames([]string)
	SetID(string)
	SetName(string)
	SetVisibleColumns([]string)
	VisibleColumns() []string
}

type List interface {
	Cursor
	Metadata
	Selectable
	Add(Row)
	All() []Row
	Clear()
	InRange(int) bool
	InsertList(source List, position int) error
	Keys() []string
	Len() int
	Lock()
	NextOf([]string, int, int) int
	Row(int) Row
	RowByID(string) Row
	RowNum(string) (int, error)
	SetUpdated()
	Sort([]string) error
	Unlock()
	Updated() time.Time
}

type Base struct {
	columnNames     []string
	columns         map[string]*Column
	cursor          int
	id              string
	mutex           sync.Mutex
	name            string
	rows            []Row
	selection       map[int]struct{}
	sortKey         string
	updated         time.Time
	visibleColumns  []string
	visualSelection [3]int
}

func New() *Base {
	s := &Base{}
	s.Clear()
	return s
}

func (s *Base) Clear() {
	s.rows = make([]Row, 0)
	s.columnNames = make([]string, 0)
	s.visibleColumns = make([]string, 0)
	s.columns = make(map[string]*Column)
	s.ClearSelection()
}

func (s *Base) ID() string {
	return s.id
}

func (s *Base) SetID(id string) {
	s.id = id
}

func (s *Base) SetColumnNames(names []string) {
	s.columnNames = make([]string, len(names))
	copy(s.columnNames, names)
}

func (s *Base) ColumnNames() []string {
	names := make([]string, 0, len(s.columns))
	for key := range s.columns {
		names = append(names, key)
	}
	return names
}

func (s *Base) SetVisibleColumns(names []string) {
	s.visibleColumns = make([]string, len(names))
	copy(s.visibleColumns, names)
}

func (s *Base) VisibleColumns() []string {
	return s.visibleColumns
}

func (s *Base) Columns(names []string) []Column {
	cols := make([]Column, len(names))
	for i, name := range names {
		if col, ok := s.columns[name]; ok {
			cols[i] = *col
		}
	}
	return cols
}

func (s *Base) Add(row Row) {
	s.rows = append(s.rows, row)
	for k, v := range row.Fields() {
		if s.columns[k] == nil {
			s.columns[k] = &Column{}
		}
		s.columns[k].Add(v)
	}
}

func (s *Base) All() []Row {
	rows := make([]Row, len(s.rows))
	for i := 0; i < len(rows); i++ {
		rows[i] = s.rows[i]
	}
	return rows
}

func (s *Base) Row(n int) Row {
	if !s.InRange(n) {
		return nil
	}
	return s.rows[n]
}

func (s *Base) RowNum(id string) (int, error) {
	for n, row := range s.rows {
		if row.ID() == id {
			return n, nil
		}
	}
	return 0, fmt.Errorf("not found")
}

func (s *Base) RowByID(id string) Row {
	rown, err := s.RowNum(id)
	if err != nil {
		return nil
	}
	return s.Row(rown)
}

func (s *Base) Keys() []string {
	keys := make([]string, s.Len())
	for i := range s.rows {
		keys[i] = s.rows[i].ID()
	}
	return keys
}

func (s *Base) Len() int {
	return len(s.rows)
}

// Implements sort.Interface
func (s *Base) Less(i, j int) bool {
	return s.rows[i].Fields()[s.sortKey] < s.rows[j].Fields()[s.sortKey]
}

// Implements sort.Interface
func (s *Base) Swap(i, j int) {
	row := s.rows[i]
	s.rows[i] = s.rows[j]
	s.rows[j] = row
}

// Sort first sorts unstable, then stable, by all columns provided.
// Retains cursor position.
func (s *Base) Sort(cols []string) error {
	if s.Len() < 2 {
		return nil
	}

	// Obtain row under cursor
	cursorRow := s.CursorRow()

	fn := sort.Sort
	for _, key := range cols {
		s.sortKey = key
		fn(s)
		fn = sort.Stable
	}

	// Restore cursor position to row previously selected
	rowNum, err := s.RowNum(cursorRow.ID())
	if err != nil {
		// panics here because the row with this id must also be found in the sorted list,
		// otherwise this is a bug.
		panic(err)
	}

	s.SetCursor(rowNum)

	return nil
}

// InRange returns true if the provided index is within list range, false otherwise.
func (s *Base) InRange(index int) bool {
	return index >= 0 && index < s.Len()
}

func (s *Base) Lock() {
	s.mutex.Lock()
}

func (s *Base) Unlock() {
	s.mutex.Unlock()
}

func (s *Base) Name() string {
	return s.name
}

func (s *Base) SetName(name string) {
	s.name = name
}

// Updated returns the timestamp of when this songlist was last updated.
func (s *Base) Updated() time.Time {
	return s.updated
}

// SetUpdated sets the update timestamp of the songlist.
func (s *Base) SetUpdated() {
	s.updated = time.Now()
}

func (s *Base) InsertList(source List, position int) error {
	sourceRows := source.All()
	if s.Len() == 0 {
		s.rows = sourceRows
	} else if position == 0 {
		s.rows = append(sourceRows, s.rows...)
	} else if position == s.Len() {
		s.rows = append(s.rows, sourceRows...)
	} else if s.InRange(position) {
		s.rows = append(s.rows[:position], append(sourceRows, s.rows[position:]...)...)
	} else {
		return fmt.Errorf("out of range")
	}

	s.SetUpdated()

	return nil
}

// NextOf searches forwards or backwards for rows having different values in the specified tags.
// The index of the next song is returned.
func (s *Base) NextOf(tags []string, index int, direction int) int {
	offset := func(i int) int {
		if direction > 0 || i == 0 {
			return 0
		}
		return 1
	}

	ln := s.Len()
	index -= offset(index)
	row := s.Row(index)

LOOP:
	for ; index < ln && index >= 0; index += direction {
		for _, tag := range tags {
			if row.Fields()[tag] != s.rows[index].Fields()[tag] {
				break LOOP
			}
		}
	}

	return index + offset(index)
}

func (s *Base) Remove(index int) error {
	row := s.Row(index)
	if row == nil {
		return fmt.Errorf("out of bounds")
	}

	for k, v := range row.Fields() {
		s.columns[k].Remove(v)
	}

	if index+1 == s.Len() {
		s.rows = s.rows[:index]
	} else {
		s.rows = append(s.rows[:index], s.rows[index+1:]...)
	}

	return nil
}

// RemoveIndices removes a selection of songs from the songlist, having the
// index defined by the int slice parameter.
func (s *Base) RemoveIndices(indices []int) error {
	// Ensure that indices are removed in reverse order
	sort.Sort(sort.Reverse(sort.IntSlice(indices)))
	for _, i := range indices {
		if err := s.Remove(i); err != nil {
			return err
		}
	}
	return nil
}
