// Code generated by mockery 2.10.0. DO NOT EDIT.

package api

import (
	clipboard "github.com/ambientsound/visp/clipboard"
	db "github.com/ambientsound/visp/db"

	keys "github.com/ambientsound/visp/input/keys"

	list "github.com/ambientsound/visp/list"

	mock "github.com/stretchr/testify/mock"

	multibar "github.com/ambientsound/visp/multibar"

	oauth2 "golang.org/x/oauth2"

	player "github.com/ambientsound/visp/player"

	spotify "github.com/zmb3/spotify"

	spotify_library "github.com/ambientsound/visp/spotify/library"

	style "github.com/ambientsound/visp/style"
)

// MockAPI is an autogenerated mock type for the API type
type MockAPI struct {
	mock.Mock
}

// Authenticate provides a mock function with given fields: token
func (_m *MockAPI) Authenticate(token *oauth2.Token) error {
	ret := _m.Called(token)

	var r0 error
	if rf, ok := ret.Get(0).(func(*oauth2.Token) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Changed provides a mock function with given fields: typ, data
func (_m *MockAPI) Changed(typ ChangeType, data interface{}) {
	_m.Called(typ, data)
}

// Clipboards provides a mock function with given fields:
func (_m *MockAPI) Clipboards() *clipboard.List {
	ret := _m.Called()

	var r0 *clipboard.List
	if rf, ok := ret.Get(0).(func() *clipboard.List); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clipboard.List)
		}
	}

	return r0
}

// Db provides a mock function with given fields:
func (_m *MockAPI) Db() *db.List {
	ret := _m.Called()

	var r0 *db.List
	if rf, ok := ret.Get(0).(func() *db.List); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*db.List)
		}
	}

	return r0
}

// Exec provides a mock function with given fields: _a0
func (_m *MockAPI) Exec(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// History provides a mock function with given fields:
func (_m *MockAPI) History() list.List {
	ret := _m.Called()

	var r0 list.List
	if rf, ok := ret.Get(0).(func() list.List); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(list.List)
		}
	}

	return r0
}

// Library provides a mock function with given fields:
func (_m *MockAPI) Library() *spotify_library.List {
	ret := _m.Called()

	var r0 *spotify_library.List
	if rf, ok := ret.Get(0).(func() *spotify_library.List); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*spotify_library.List)
		}
	}

	return r0
}

// List provides a mock function with given fields:
func (_m *MockAPI) List() list.List {
	ret := _m.Called()

	var r0 list.List
	if rf, ok := ret.Get(0).(func() list.List); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(list.List)
		}
	}

	return r0
}

// Multibar provides a mock function with given fields:
func (_m *MockAPI) Multibar() *multibar.Multibar {
	ret := _m.Called()

	var r0 *multibar.Multibar
	if rf, ok := ret.Get(0).(func() *multibar.Multibar); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*multibar.Multibar)
		}
	}

	return r0
}

// PlayerStatus provides a mock function with given fields:
func (_m *MockAPI) PlayerStatus() player.State {
	ret := _m.Called()

	var r0 player.State
	if rf, ok := ret.Get(0).(func() player.State); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(player.State)
	}

	return r0
}

// Quit provides a mock function with given fields:
func (_m *MockAPI) Quit() {
	_m.Called()
}

// Sequencer provides a mock function with given fields:
func (_m *MockAPI) Sequencer() *keys.Sequencer {
	ret := _m.Called()

	var r0 *keys.Sequencer
	if rf, ok := ret.Get(0).(func() *keys.Sequencer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*keys.Sequencer)
		}
	}

	return r0
}

// SetList provides a mock function with given fields: _a0
func (_m *MockAPI) SetList(_a0 list.List) {
	_m.Called(_a0)
}

// Spotify provides a mock function with given fields:
func (_m *MockAPI) Spotify() (*spotify.Client, error) {
	ret := _m.Called()

	var r0 *spotify.Client
	if rf, ok := ret.Get(0).(func() *spotify.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*spotify.Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Styles provides a mock function with given fields:
func (_m *MockAPI) Styles() style.Stylesheet {
	ret := _m.Called()

	var r0 style.Stylesheet
	if rf, ok := ret.Get(0).(func() style.Stylesheet); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(style.Stylesheet)
		}
	}

	return r0
}

// UI provides a mock function with given fields:
func (_m *MockAPI) UI() UI {
	ret := _m.Called()

	var r0 UI
	if rf, ok := ret.Get(0).(func() UI); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(UI)
		}
	}

	return r0
}
