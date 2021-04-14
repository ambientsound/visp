package commands

import (
	"fmt"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
	spotify_proxyclient "github.com/ambientsound/visp/spotify/proxyclient"
)

// Auth runs OAuth2 authentication flow against Spotify.
type Auth struct {
	command
	api   api.API
	token string
}

func NewAuth(api api.API) Command {
	return &Auth{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Auth) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()
	if tok == lexer.TokenIdentifier {
		cmd.token = lit
	} else {
		return fmt.Errorf("unexpected '%s'; expected token from web page", lit)
	}
	return cmd.ParseEnd()
}

func (cmd *Auth) Exec() error {
	token, err := spotify_proxyclient.DecodeTokenString(cmd.token)
	if err != nil {
		return err
	}
	return cmd.api.Authenticate(token)
}
