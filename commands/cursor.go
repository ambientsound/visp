package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/ambientsound/visp/list"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
)

// Cursor moves the cursor in a songlist widget. It can take human-readable
// parameters such as 'up' and 'down', and it also accepts relative positions
// if a number is given.
type Cursor struct {
	command
	api             api.API
	absolute        int
	current         bool
	finished        bool
	list            list.List
	nextOfDirection int
	nextOfTags      []string
	relative        int
}

// NewCursor returns Cursor.
func NewCursor(api api.API) Command {
	return &Cursor{
		api: api,
	}
}

// Parse parses cursor movement.
func (cmd *Cursor) Parse() error {
	tableWidget := cmd.api.UI().TableWidget()
	cmd.list = cmd.api.List()

	tok, lit := cmd.ScanIgnoreWhitespace()
	cmd.setTabCompleteVerbs(lit)

	switch tok {
	// In case of a number, scan the actual number and return
	case lexer.TokenMinus, lexer.TokenPlus:
		cmd.setTabCompleteEmpty()
		cmd.Unscan()
		_, lit, absolute, err := cmd.ParseInt()
		if err != nil {
			return err
		}
		if absolute {
			cmd.absolute = lit
		} else {
			cmd.relative = lit
		}
		return cmd.ParseEnd()

	case lexer.TokenIdentifier:
	default:
		return fmt.Errorf("unexpected '%v', expected number or identifier", lit)
	}

	switch lit {
	case "up":
		cmd.relative = -1
	case "down":
		cmd.relative = 1
	case "home":
		cmd.absolute = 0
	case "end":
		cmd.absolute = cmd.list.Len() - 1
	case "high":
		ymin, _ := tableWidget.GetVisibleBoundaries()
		cmd.absolute = ymin
	case "middle":
		ymin, ymax := tableWidget.GetVisibleBoundaries()
		cmd.absolute = (ymin + ymax) / 2
	case "low":
		_, ymax := tableWidget.GetVisibleBoundaries()
		cmd.absolute = ymax
	case "current":
		cmd.current = true
	case "random":
		cmd.absolute = cmd.random()
	case "nextOf":
		cmd.nextOfDirection = 1
		return cmd.parseNextOf()
	case "prevOf":
		cmd.nextOfDirection = -1
		return cmd.parseNextOf()
	default:
		i, err := strconv.Atoi(lit)
		if err != nil {
			return fmt.Errorf("cursor command '%s' not recognized, and is not a number", lit)
		}
		cmd.relative = i
	}

	cmd.setTabCompleteEmpty()

	return cmd.ParseEnd()
}

// Exec implements Command
func (cmd *Cursor) Exec() error {
	switch {
	case cmd.nextOfDirection != 0:
		cmd.absolute = cmd.runNextOf()

	case cmd.current:
		track := cmd.api.PlayerStatus().Item
		if track == nil {
			return fmt.Errorf("no track is currently playing")
		}

		tl := cmd.api.List()

		err := tl.SetCursorByID(track.ID.String())
		if err != nil {
			return fmt.Errorf("currently playing track is not in this list")
		}

		return nil
	}

	switch {
	case cmd.relative != 0:
		cmd.list.MoveCursor(cmd.relative)
	default:
		cmd.list.SetCursor(cmd.absolute)
	}

	return nil
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *Cursor) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"current",
		"down",
		"end",
		"high",
		"home",
		"low",
		"middle",
		"nextOf",
		"prevOf",
		"random",
		"up",
	})
}

// random returns a random list index in the songlist.
func (cmd *Cursor) random() int {
	ln := cmd.list.Len()
	if ln == 0 {
		return cmd.absolute
	}
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return r.Int() % ln
}

// parseNextOf assigns the nextOf tags and directions, or returns an error if
// no tags are specified.
func (cmd *Cursor) parseNextOf() error {
	var err error
	cmd.nextOfTags, err = cmd.ParseTags(cmd.list.ColumnNames())
	return err
}

// runNextOf finds the next song with different tags.
func (cmd *Cursor) runNextOf() int {
	return cmd.list.NextOf(cmd.nextOfTags, cmd.list.Cursor(), cmd.nextOfDirection)
}
