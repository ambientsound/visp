package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/options"
	"github.com/ambientsound/visp/spotify/aggregator"
	"github.com/google/uuid"

	"github.com/ambientsound/visp/api"
)

var (
	tagMaps = map[string]string{
		"albumArtist": "artist",
	}
)

// Isolate searches for songs that have similar tags as the selection.
type Isolate struct {
	command
	api  api.API
	tags []string
	list list.List
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
	cmd.list = cmd.api.List()
	cmd.tags, err = cmd.ParseTags(cmd.list.ColumnNames())
	return err
}

// Exec implements Command.
func (cmd *Isolate) Exec() error {
	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	selection := cmd.list.Selection()
	if selection.Len() != 1 {
		return fmt.Errorf("isolate operates on exactly one track")
	}

	row := selection.Row(0)
	if row.Kind() != list.DataTypeTrack {
		return fmt.Errorf("isolate needs a row of type '%s', not '%s'", list.DataTypeTrack, row.Kind())
	}

	queries := make([]string, len(cmd.tags))
	for i, tag := range cmd.tags {
		val := strconv.Quote(row.Fields()[tag])
		if v, ok := tagMaps[tag]; ok {
			tag = v
		}
		queries[i] = fmt.Sprintf("%s:%s", tag, val)
	}

	query := strings.Join(queries, " AND ")
	log.Debugf("isolate search: %s", query)
	result, err := spotify_aggregator.Search(*client, query, options.GetInt(options.Limit))

	if err != nil {
		return err
	}

	if result.Len() == 0 {
		return fmt.Errorf("no results found when isolating by %s", strings.Join(cmd.tags, ", "))
	}

	// Post-processing: sort in default order
	sort := options.GetString(options.SortTracklists)

	err = result.Sort(strings.Split(sort, ","))
	if err != nil {
		log.Errorf("error sorting: %s", err)
	}

	result.SetVisibleColumns(cmd.list.VisibleColumns())
	result.SetID(uuid.New().String())
	_ = result.SetCursorByID(row.ID())

	// Figure out a clever name
	parts := make([]string, len(cmd.tags))
	for i, tag := range cmd.tags {
		parts[i] = tag + ":" + row.Fields()[tag]
	}
	result.SetName(strings.Join(parts, ", "))

	// Activate results
	cmd.api.SetList(result)

	return nil
}
