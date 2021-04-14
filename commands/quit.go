package commands

import (
	"github.com/ambientsound/visp/api"
)

// Quit exits the program.
type Quit struct {
	command
	api api.API
}

func NewQuit(api api.API) Command {
	return &Quit{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Quit) Parse() error {
	return cmd.ParseEnd()
}

func (cmd *Quit) Exec() error {
	cmd.api.Quit()
	return nil
}
