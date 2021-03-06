package spotify_aggregator

import (
	"context"

	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/options"
	spotify_albums "github.com/ambientsound/visp/spotify/albums"
	"github.com/ambientsound/visp/spotify/library"
	"github.com/ambientsound/visp/spotify/playlists"
	"github.com/ambientsound/visp/spotify/tracklist"
	"github.com/zmb3/spotify/v2"
)

func Search(client spotify.Client, query string, limit int) (list.List, error) {
	results, err := client.Search(
		context.TODO(),
		query,
		spotify.SearchTypeTrack,
		spotify.Limit(limit),
	)
	if err != nil {
		return nil, err
	}

	lst, err := spotify_tracklist.NewFromFullTrackPage(client, results.Tracks)
	if err != nil {
		return nil, err
	}

	lst.Sort(options.GetList(options.SortSearch))
	lst.SetCursor(0)

	return lst, nil
}

func FeaturedPlaylists(client spotify.Client, limit int) (*spotify_playlists.List, error) {
	message, playlists, err := client.FeaturedPlaylists(
		context.TODO(),
		spotify.Limit(limit),
	)
	if err != nil {
		return nil, err
	}

	lst, err := spotify_playlists.New(client, playlists)
	if err != nil {
		return nil, err
	}

	lst.SetName(message)
	lst.SetID(spotify_library.FeaturedPlaylists)
	lst.SetVisibleColumns(options.GetList(options.ColumnsPlaylists))
	lst.Sort(options.GetList(options.SortPlaylists))
	lst.SetCursor(0)

	return lst, nil
}

func ListWithID(client spotify.Client, id string, limit int) (list.List, error) {
	sid := spotify.ID(id)

	playlist, err := client.GetPlaylist(context.TODO(), sid)
	if err != nil {
		return nil, err
	}

	tracks, err := client.GetPlaylistTracks(context.TODO(), sid, spotify.Limit(limit))
	if err != nil {
		return nil, err
	}

	lst, err := spotify_tracklist.NewFromPlaylistTrackPage(client, tracks)
	if err != nil {
		return nil, err
	}

	lst.SetName(playlist.Name)
	lst.SetID(id)
	lst.SetURI(playlist.URI)
	lst.SetRemote(true)
	lst.SetSyncedToRemote()
	lst.SetVisibleColumns(options.GetList(options.ColumnsTracklists))

	// don't sort list with ID's, their order are significant.

	return lst, nil
}

func MyPrivatePlaylists(client spotify.Client, limit int) (*spotify_playlists.List, error) {
	playlists, err := client.CurrentUsersPlaylists(context.TODO(), spotify.Limit(limit))
	if err != nil {
		return nil, err
	}

	lst, err := spotify_playlists.New(client, playlists)
	if err != nil {
		return nil, err
	}

	lst.SetName("My playlists")
	lst.SetID(spotify_library.MyPlaylists)
	lst.SetVisibleColumns(options.GetList(options.ColumnsPlaylists))
	lst.Sort(options.GetList(options.SortPlaylists))
	lst.SetCursor(0)

	return lst, nil
}

func MyTracks(client spotify.Client, limit int) (list.List, error) {
	tracks, err := client.CurrentUsersTracks(context.TODO(), spotify.Limit(limit))
	if err != nil {
		return nil, err
	}

	lst, err := spotify_tracklist.NewFromSavedTrackPage(client, tracks)
	if err != nil {
		return nil, err
	}

	lst.SetName("Saved tracks")
	lst.SetID(spotify_library.MyTracks)
	lst.SetVisibleColumns(options.GetList(options.ColumnsTracklists))
	lst.Sort(options.GetList(options.SortTracklists))
	lst.SetCursor(0)

	return lst, nil
}

func MyAlbums(client spotify.Client) (*spotify_albums.List, error) {
	albums, err := client.CurrentUsersAlbums(context.TODO())
	if err != nil {
		return nil, err
	}

	lst, err := spotify_albums.NewFromSavedAlbumPage(client, albums)
	if err != nil {
		return nil, err
	}

	lst.SetName("Saved albums")
	lst.SetID(spotify_library.MyAlbums)
	lst.SetVisibleColumns(options.GetList(options.ColumnsAlbums))
	lst.Sort(options.GetList(options.SortAlbums))
	lst.SetCursor(0)

	return lst, nil
}

func TopTracks(client spotify.Client, limit int) (list.List, error) {
	tracks, err := client.CurrentUsersTopTracks(context.TODO(), spotify.Limit(limit))
	if err != nil {
		return nil, err
	}

	lst, err := spotify_tracklist.NewFromFullTrackPage(client, tracks)
	if err != nil {
		return nil, err
	}

	lst.SetName("Top tracks")
	lst.SetID(spotify_library.TopTracks)
	lst.SetVisibleColumns(options.GetList(options.ColumnsTracklists))

	// the order of this list is significant, but it's probably best to sort it for better UX.
	lst.Sort(options.GetList(options.SortTracklists))
	lst.SetCursor(0)

	return lst, nil
}

func NewReleases(client spotify.Client) (list.List, error) {
	albums, err := client.NewReleases(context.TODO())
	if err != nil {
		return nil, err
	}

	lst, err := spotify_tracklist.NewFromSimpleAlbumPage(client, albums)
	if err != nil {
		return nil, err
	}

	lst.SetName("New releases")
	lst.SetID(spotify_library.NewReleases)
	lst.SetVisibleColumns(options.GetList(options.ColumnsTracklists))
	lst.Sort(options.GetList(options.SortTracklists))
	lst.SetCursor(0)

	return lst, nil
}
