package topbar

import (
	"github.com/ambientsound/visp/api"
)

// Time draws the current song's length.
type Time struct {
	api api.API
}

// NewTime returns Time.
func NewTime(a api.API, param string) Fragment {
	return &Time{a}
}

// Text implements Fragment.
func (w *Time) Text() (string, string) {
	playerStatus := w.api.PlayerStatus()
	return playerStatus.TrackRow.Fields()["time"], `time`
}
