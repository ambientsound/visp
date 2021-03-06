package commands_test

import (
	"testing"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/options"
	"github.com/stretchr/testify/mock"

	"github.com/ambientsound/visp/commands"
	"github.com/stretchr/testify/assert"
)

var setTests = []commands.Test{
	// Valid forms
	{``, true, testSetInit, nil, nil},
	{`foo=bar`, true, testSetInit, testFooSet(`foo`, `bar`, true), []string{}},
	{`foo="bar baz"`, true, testSetInit, testFooSet(`foo`, `bar baz`, true), []string{}},
	{`foo=${}|;#`, true, testSetInit, testFooSet(`foo`, `${}|;`, true), []string{}},
	{`foo=x bar=x baz=x int=4 invbool`, true, testSetInit, testMultiSet, []string{}},
	{`foo=y foo`, true, testSetInit, testFooSet(`foo`, `y`, true), []string{`foo`}},
	{`baz=`, true, testSetInit, testFooSet(`baz`, ``, true), []string{`="foobar"`, `=`}},
	{`bool`, true, testSetInit, nil, []string{`bool`}},

	// Invalid forms
	{`nonexist=foo`, true, testSetInit, testFooSet(`nonexist`, ``, false), []string{}},
	{`$=""`, false, testSetInit, nil, []string{}},
}

func TestSet(t *testing.T) {
	commands.TestVerb(t, "set", setTests)
}

func testSetInit(test *commands.TestData) {
	test.MockAPI.On("Changed", api.ChangeOption, mock.Anything).Return()
	options.Set("foo", "")
	options.Set("bar", "")
	options.Set("baz", "foobar")
	options.Set("int", 0)
	options.Set("bool", false)
}

func testFooSet(key, check string, ok bool) func(*commands.TestData) {
	return func(test *commands.TestData) {
		err := test.Cmd.Exec()
		assert.Equal(test.T, ok, err == nil, "Expected OK=%s", ok)
		if err != nil {
			return
		}
		val := options.GetString(key)
		assert.Equal(test.T, check, val)
	}
}

func testMultiSet(test *commands.TestData) {
	err := test.Cmd.Exec()
	assert.Nil(test.T, err)
	assert.Equal(test.T, "x", options.GetString("foo"))
	assert.Equal(test.T, "x", options.GetString("bar"))
	assert.Equal(test.T, "x", options.GetString("baz"))
	assert.Equal(test.T, 4, options.GetInt("int"))
	assert.Equal(test.T, true, options.GetBool("bool"))
}
