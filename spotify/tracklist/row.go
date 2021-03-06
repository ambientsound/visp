package spotify_tracklist

import (
	"github.com/ambientsound/visp/list"
	"github.com/zmb3/spotify/v2"
)

type Row struct {
	*list.BaseRow
	track spotify.FullTrack
}

func (row *Row) Track() spotify.FullTrack {
	return row.track
}
