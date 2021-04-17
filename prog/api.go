package prog

import (
	"context"
	"fmt"
	"time"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/clipboard"
	"github.com/ambientsound/visp/db"
	"github.com/ambientsound/visp/input/keys"
	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/multibar"
	"github.com/ambientsound/visp/options"
	"github.com/ambientsound/visp/player"
	"github.com/ambientsound/visp/spotify/library"
	"github.com/ambientsound/visp/spotify/proxyclient"
	"github.com/ambientsound/visp/spotify/tracklist"
	"github.com/ambientsound/visp/style"
	"github.com/ambientsound/visp/topbar"
	"github.com/spf13/viper"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

func (v *Visp) Authenticate(token *oauth2.Token) error {
	log.Infof("Configured Spotify access token, expires at %s", token.Expiry.Format(time.RFC1123))

	cli := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	scli := spotify.NewClient(cli)

	v.client = &scli

	next := spotify_proxyclient.TokenTTL(token)

	v.tokenRefresh = time.After(next)

	err := v.Tokencache.Write(*token)
	if err != nil {
		return fmt.Errorf("write Spotify token to file: %s", err)
	}

	return nil
}

func (v *Visp) Clipboards() *clipboard.List {
	return v.clipboards
}

func (v *Visp) Db() *db.List {
	return v.db
}

func (v *Visp) Exec(command string) error {
	log.Debugf("Run command: %s", command)
	return v.interpreter.Exec(command)
}

func (v *Visp) Library() *spotify_library.List {
	return v.library
}

func (v *Visp) List() list.List {
	return v.list
}

func (v *Visp) ListChanged() {
	// FIXME
}

func (v *Visp) Message(fmt string, a ...interface{}) {
	log.Infof(fmt, a...)
	log.Debugf("Using obsolete Message() for previous message")
}

func (v *Visp) OptionChanged(key string) {
	switch key {
	case options.LogFile:
		logFile := v.Options().GetString(options.LogFile)
		overwrite := v.Options().GetBool(options.LogOverwrite)
		if len(logFile) == 0 {
			break
		}
		err := log.Configure(logFile, overwrite)
		if err != nil {
			log.Errorf("log configuration: %s", err)
			break
		}
		log.Infof("Note: log file will be backfilled with existing log")
		log.Infof("Writing debug log to %s", logFile)

	case options.Topbar:
		config := v.Options().GetString(options.Topbar)
		matrix, err := topbar.Parse(v, config)
		if err == nil {
			v.Termui.Widgets.Topbar.SetMatrix(matrix)
			v.Termui.Resize()
		} else {
			log.Errorf("topbar configuration: %s", err)
		}
	}
}

func (v *Visp) Options() api.Options {
	return viper.GetViper()
}

func (v *Visp) PlayerStatus() player.State {
	return v.player
}

func (v *Visp) Quit() {
	v.quit <- new(interface{})
}

func (v *Visp) Sequencer() *keys.Sequencer {
	return v.sequencer
}

func (v *Visp) Multibar() *multibar.Multibar {
	return v.multibar
}

func (v *Visp) SetList(lst list.List) {
	c := v.db.Cache(lst)
	v.db.SetCursor(c)
	v.list = lst
	v.Termui.TableWidget().SetList(lst)
}

func (v *Visp) Spotify() (*spotify.Client, error) {
	if v.client == nil {
		return nil, fmt.Errorf("please authenticate with Spotify at: %s", v.Options().GetString("spotifyauthserver"))
	}
	token, err := v.client.Token()
	if err != nil {
		return nil, fmt.Errorf("unable to refresh Spotify token: %s", err)
	}
	_ = v.Tokencache.Write(*token)
	return v.client, nil
}

func (v *Visp) Styles() style.Stylesheet {
	return v.stylesheet
}

func (v *Visp) Tracklist() *spotify_tracklist.List {
	switch v := v.List().(type) {
	case *spotify_tracklist.List:
		return v
	default:
		return nil
	}
}

func (v *Visp) UI() api.UI {
	return v.Termui
}
