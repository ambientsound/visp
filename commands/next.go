package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
)

// Next switches to the next song in MPD's queue.
type Next struct {
	api API
}

func NewNext(api API) Command {
	return &Next{
		api: api,
	}
}

func (cmd *Next) Execute(t lexer.Token) error {
	switch t.Class {
	case lexer.TokenEnd:
		client := cmd.api.MpdClient()
		if client == nil {
			return fmt.Errorf("Unable to play next song: cannot communicate with MPD")
		}
		return client.Next()

	default:
		return fmt.Errorf("Unknown input '%s', expected END", t.String())
	}

	return nil
}