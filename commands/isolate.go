package commands

import (
	"fmt"
	"github.com/ambientsound/pms/options"
	"strings"

	"github.com/ambientsound/pms/api"
)

// Isolate searches for songs that have similar tags as the selection.
type Isolate struct {
	newcommand
	api  api.API
	tags []string
}

// NewIsolate returns Isolate.
func NewIsolate(api api.API) Command {
	return &Isolate{
		api:  api,
		tags: make([]string, 0),
	}
}

// Parse implements Command.
func (cmd *Isolate) Parse() error {
	var err error
	list := cmd.api.Songlist()
	cmd.tags, err = cmd.ParseTags(list.CursorSong().TagKeys())
	return err
}

// Exec implements Command.
func (cmd *Isolate) Exec() error {
	library := cmd.api.Library()
	if library == nil {
		return fmt.Errorf("Song library is not present.")
	}

	db := cmd.api.Db()
	panel := db.Panel()
	list := cmd.api.Songlist()
	selection := list.Selection()
	song := list.CursorSong()

	if selection.Len() == 0 {
		return fmt.Errorf("Isolate needs at least one track.")
	}

	result, err := library.Isolate(selection, cmd.tags)
	if err != nil {
		return err
	}

	if result.Len() == 0 {
		return fmt.Errorf("No results found when isolating by %s", strings.Join(cmd.tags, ", "))
	}

	// Sort the new list.
	sort := cmd.api.Options().GetString(options.Sort)
	fields := strings.Split(sort, ",")
	result.Sort(fields)

	// Clear selection in the source list, and add a new list to the index.
	list.ClearSelection()
	panel.Add(result)
	panel.Activate(result)
	list.CursorToSong(song)

	return nil
}
