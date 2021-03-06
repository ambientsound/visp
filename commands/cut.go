package commands

import (
	"fmt"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
)

// Cut removes songs from songlists.
type Cut struct {
	command
	api  api.API
	list list.List
}

// NewCut returns Cut.
func NewCut(api api.API) Command {
	return &Cut{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Cut) Parse() error {
	cmd.list = cmd.api.List()
	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Cut) Exec() error {

	selection := cmd.list.Selection()
	indices := cmd.list.SelectionIndices()
	ln := len(indices)

	if ln == 0 {
		return fmt.Errorf("no tracks selected")
	}

	// Remove songs from list
	index := indices[0]
	err := cmd.list.RemoveIndices(indices)

	cmd.api.Changed(api.ChangeList, cmd.list)

	if err != nil {
		return err
	}

	cmd.list.ClearSelection()
	cmd.list.SetCursor(index)

	selection.SetVisibleColumns(cmd.list.ColumnNames())

	cmd.api.Clipboards().Insert(selection)

	log.Infof("%d fewer songs; stored in %s", ln, selection.Name())

	return nil
}
