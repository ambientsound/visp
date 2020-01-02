package api

import (
	"github.com/ambientsound/pms/list"
	"github.com/ambientsound/pms/songlist"
	"github.com/spf13/viper"
)

type Window int

const (
	WindowLogs Window = iota
	WindowMusic
	WindowPlaylists
)

type Collection interface {
	Activate(songlist.Songlist)
	ActivateIndex(int) error
	Add(songlist.Songlist)
	Current() songlist.Songlist
	Index() (int, error)
	Last() songlist.Songlist
	Len() int
	Remove(int) error
	ValidIndex(int) bool
}

type TableWidget interface {
	GetVisibleBoundaries() (int, int)
	List() list.List
	PositionReadout() string
	ScrollViewport(int, bool)
	SetColumns([]string)
	SetList(list.List)
	Size() (int, int)
}

type UI interface {
	ActivateWindow(Window)
	Refresh()
	TableWidget() TableWidget
}

type Options interface {
	AllKeys() []string
	Set(string, interface{})
	Get(string) interface{}
	GetString(string) string
	GetInt(string) int
	GetBool(string) bool
}

var _ Options = viper.GetViper()
