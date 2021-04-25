package widgets

import (
	"fmt"
	"math"
	"time"

	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/options"
	"github.com/ambientsound/visp/spotify/devices"
	"github.com/ambientsound/visp/spotify/tracklist"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/style"
	"github.com/ambientsound/visp/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type lineStyler func(row list.Row) (string, bool)

// Table is a tcell widget which draws a gridded table from a List instance.
type Table struct {
	api     api.API
	columns []column
	list    list.List

	view     views.View
	viewport views.ViewPort
	lastDraw time.Time

	style.Styled
	views.WidgetWatchers
}

func NewTable(a api.API) *Table {
	return &Table{
		api: a,
	}
}

func (w *Table) drawNext(v views.View, x, y, strmin, strmax int, runes []rune, style tcell.Style) int {
	strmin = utils.Min(len(runes), strmin)
	n := 0
	for n < strmin {
		v.SetContent(x, y, runes[n], nil, style)
		n++
		x++
	}
	for n < strmax {
		v.SetContent(x, y, ' ', nil, style)
		n++
		x++
	}
	return x
}

func (w *Table) List() list.List {
	return w.list
}

func (w *Table) SetList(lst list.List) {
	w.list = lst
	w.SetColumns(lst.VisibleColumns())
	log.Debugf("SetList: %v", w.list.Name())
	log.Debugf("SetColumns: %v", w.ColumnNames())
}

func (w *Table) Draw() {
	var styler lineStyler
	var specialStyler lineStyler
	var st tcell.Style

	w.SetStylesheet(w.api.Styles())

	// Make sure that the viewport matches the list size.
	w.setViewportSize()

	// Update draw time
	w.lastDraw = time.Now()

	// a, b, c, d := w.viewport.GetPhysical()
	// log.Debugf("Drawing table widget on viewport: %#v", w.viewport)
	// log.Debugf("Visible phys coordinates: (%v,%v) (%v,%v)", a, b, c, d)
	_, ymin, xmax, ymax := w.viewport.GetVisible()
	x, y := 0, 0
	xmax += 1
	cursor := false

	_, isTracklist := w.list.(*spotify_tracklist.List)
	_, isDevicelist := w.list.(*spotify_devices.List)

	// Special line styling based on list type
	switch {
	case isTracklist && w.api.PlayerStatus().Item != nil:
		trackID := w.api.PlayerStatus().Item.ID.String()
		specialStyler = func(row list.Row) (string, bool) {
			return `currentSong`, trackID == row.ID()
		}
	case isDevicelist:
		id := w.api.PlayerStatus().Device.ID
		deviceID := id.String()
		specialStyler = func(row list.Row) (string, bool) {
			return `currentDevice`, deviceID == row.ID()
		}
	default:
		specialStyler = func(row list.Row) (string, bool) {
			return ``, false
		}
	}

	// Generic line styling.
	styler = func(row list.Row) (string, bool) {
		st, special := specialStyler(row)
		switch {
		case cursor:
			return `cursor`, true
		case special:
			return st, special
		case w.list.Selected(y):
			return `selection`, true
		default:
			return `default`, false
		}
	}

	w.drawHeaders()

	for y = ymin; y <= ymax; y++ {
		row := w.list.Row(y)
		if row == nil {
			panic(fmt.Sprintf("nil row: %d", y))
		}

		x = 0
		cursor = y == w.list.Cursor()
		styleName, lineStyled := styler(row)
		st = w.Style(styleName)

		// Draw each column separately
		for _, col := range w.columns {

			runes := []rune(row.Fields()[col.key])
			if !lineStyled {
				st = w.Style(col.key)
			}

			strmin := col.width - col.rightPadding

			x = w.drawNext(&w.viewport, x, y, strmin, col.width, runes, st)
		}
	}
}

func (w *Table) drawHeaders() {
	x := 0
	st := w.Style("header")
	for _, col := range w.columns {
		runes := []rune(col.title)
		strmin := col.width - col.rightPadding
		x = w.drawNext(w.view, x, 0, strmin, col.width, runes, st)
	}
}

func (w *Table) GetVisibleBoundaries() (ymin, ymax int) {
	_, ymin, _, ymax = w.viewport.GetVisible()
	return
}

// Width returns the widget width.
func (w *Table) Width() int {
	_, _, xmax, _ := w.viewport.GetVisible()
	return xmax
}

// Height returns the widget height.
func (w *Table) Height() int {
	_, ymin, _, ymax := w.viewport.GetVisible()
	return ymax - ymin
}

func (w *Table) setViewportSize() {
	x, y := w.Size()
	w.viewport.SetContentSize(x, w.list.Len(), true)
	w.viewport.SetSize(x, utils.Min(y-1, w.list.Len()))
	w.validateViewport()
}

// validateViewport moves the visible viewport so that the cursor is made visible.
// If the 'center' option is enabled, the viewport is centered on the cursor.
func (w *Table) validateViewport() {
	cursor := w.list.Cursor()

	// Make the cursor visible
	if !options.GetBool(options.Center) {
		w.viewport.MakeVisible(0, cursor)
		return
	}

	// If 'center' is on, make the cursor centered.
	half := w.Height() / 2
	min := utils.Max(0, cursor-half)
	max := utils.Min(w.list.Len()-1, cursor+half)
	w.viewport.MakeVisible(0, min)
	w.viewport.MakeVisible(0, max)
}

func (w *Table) Resize() {
	x, y := w.Size()
	w.viewport.Resize(0, 1, x, y-1)
	w.SetColumns(w.ColumnNames())
}

func (w *Table) HandleEvent(ev tcell.Event) bool {
	return false
}

func (w *Table) SetView(v views.View) {
	w.view = v
	w.viewport = *views.NewViewPort(v, 0, 0, -1, -1)
}

func (w *Table) Size() (int, int) {
	return w.view.Size()
}

func (w *Table) Name() string {
	return w.list.Name()
}

// PositionReadout returns a combination of PositionLongReadout() and PositionShortReadout().
// FIXME: move this into a positionreadout fragment
func (w *Table) PositionReadout() string {
	return fmt.Sprintf("%s    %s", w.PositionLongReadout(), w.PositionShortReadout())
}

// PositionLongReadout returns a formatted string containing the visible song
// range as well as the total number of songs.
// FIXME: move this into a positionreadout fragment
func (w *Table) PositionLongReadout() string {
	ymin, ymax := w.GetVisibleBoundaries()
	return fmt.Sprintf("%d,%d-%d/%d", w.list.Cursor()+1, ymin+1, ymax+1, w.list.Len())
}

// PositionShortReadout returns a percentage indicator on how far the songlist is scrolled.
// FIXME: move this into a positionreadout fragment
func (w *Table) PositionShortReadout() string {
	ymin, ymax := w.GetVisibleBoundaries()
	if ymin == 0 && ymax+1 == w.list.Len() {
		return `All`
	}
	if ymin == 0 {
		return `Top`
	}
	if ymax+1 == w.list.Len() {
		return `Bot`
	}
	fraction := float64(float64(ymin) / float64(w.list.Len()))
	percent := int(math.Floor(fraction * 100))
	return fmt.Sprintf("%2d%%", percent)
}

// ColumnNames returns a list of the visible columns
func (w *Table) ColumnNames() []string {
	keys := make([]string, len(w.columns))
	for i := range w.columns {
		keys[i] = w.columns[i].key
	}
	return keys
}

// SetColumns sets which columns that should be visible, and adjusts the sizes so they
// fit as close as possible to the median size of the content displayed.
func (w *Table) SetColumns(tags []string) {
	totalWidth, _ := w.Size()
	usedWidth := 0

	cols := w.list.Columns(tags)
	w.columns = make([]column, len(tags))

	expandColumns := options.GetList(options.ExpandColumns)
	expand := make(map[string]bool)
	for _, col := range expandColumns {
		expand[col] = true
	}

	fullHeaderColumns := options.GetList(options.FullHeaderColumns)
	fullHeader := make(map[string]bool)
	for _, col := range fullHeaderColumns {
		fullHeader[col] = true
	}

	for i, key := range tags {
		w.columns[i].col = cols[i]
		w.columns[i].key = key
		w.columns[i].title = ColumnTitle(key)
		if expand[key] {
			// auto-expanded columns start at their median size
			w.columns[i].width = cols[i].Median()
		} else if fullHeader[key] {
			// non-expanded columns with full headers start at maximum size plus one character for padding.
			// header titles are included in the maximum size calculation.
			w.columns[i].width = utils.Max(cols[i].Max(), len(w.columns[i].title)) + 1
		} else {
			// non-expanded columns start at maximum size plus one character for padding.
			// Allow at least three characters for zero-width columns.
			w.columns[i].width = utils.Max(cols[i].Max(), 4) + 1
		}
		w.columns[i].rightPadding = 1
		usedWidth += w.columns[i].width
	}

	if len(tags) == 0 {
		return
	}

	// create a list of tags that should auto-expand
	poolSize := 0
	for _, tag := range tags {
		if expand[tag] {
			poolSize++
		}
	}

	// if no columns are marked as auto-expand, skip expanison algorithm
	if poolSize > 0 {
		// expand columns to maximum length as long as there is space left
		saturated := make([]bool, len(tags))

	outer:
		for {
			for i, tag := range tags {
				if usedWidth > totalWidth {
					break outer
				}
				if !expand[tag] {
					continue
				}
				if poolSize > 0 && saturated[i] {
					continue
				}
				col := w.columns[i]
				if poolSize > 0 && col.width > col.col.Max() {
					// log.Debugf("saturating column %s at width %d", tags[i], col.width)
					saturated[i] = true
					poolSize--
					continue
				}
				w.columns[i].width++
				// log.Debugf("increase column %s to width %d", tags[i], w.columns[i].width)
				usedWidth++
			}
		}
	}

	// Set column names, preferably to their maximum size, but truncate as needed.
	for i := range tags {
		col := w.columns[i]
		if len(col.title) >= col.width && col.width > 2 {
			w.columns[i].title = col.title[:col.width-2] + "."
		}
	}
}

// ScrollViewport scrolls the viewport by delta rows, as far as possible.
// If movecursor is false, the cursor is kept pointing at the same song where
// possible. If true, the cursor is moved delta rows.
func (w *Table) ScrollViewport(delta int, movecursor bool) {
	// Do nothing if delta is zero
	if delta == 0 {
		return
	}

	if delta < 0 {
		w.viewport.ScrollUp(-delta)
	} else {
		w.viewport.ScrollDown(delta)
	}

	if movecursor {
		w.list.MoveCursor(delta)
	}

	w.validateCursor()
}

// validateCursor ensures the cursor is within the allowable area without moving
// the viewport.
func (w *Table) validateCursor() {
	ymin, ymax := w.GetVisibleBoundaries()
	cursor := w.list.Cursor()

	if options.GetBool(options.Center) {
		// When 'center' is on, move cursor to the centre of the viewport
		target := cursor
		lowerbound := (ymin + ymax) / 2
		upperbound := lowerbound
		if ymin <= 0 {
			// We are scrolled to the top, so the cursor is allowed to go above
			// the middle of the viewport
			lowerbound = 0
		}
		if ymax >= w.list.Len()-1 {
			// We are scrolled to the bottom, so the cursor is allowed to go
			// below the middle of the viewport
			upperbound = w.list.Len() - 1
		}
		if target < lowerbound {
			target = lowerbound
		}
		if target > upperbound {
			target = upperbound
		}
		w.list.SetCursor(target)
	} else {
		// When 'center' is off, move cursor into the viewport
		if cursor < ymin {
			w.list.SetCursor(ymin)
		} else if cursor > ymax {
			w.list.SetCursor(ymax)
		}
	}
}
