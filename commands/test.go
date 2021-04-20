package commands

import (
	"strings"
	"testing"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestData contains data needed for a single Command table test.
type TestData struct {
	T       *testing.T
	Cmd     Command
	Api     api.API
	MockAPI *api.MockAPI
	Test    Test
}

// Test is a structure for test data, and can be used to conveniently
// test Command instances.
type Test struct {

	// The input data for the command, as seen on the command line.
	Input string

	// True if the command should parse and execute properly, false otherwise.
	Success bool

	// An initialization function for tests.
	Init func(data *TestData)

	// A callback function to call for every test, allowing customization of tests.
	Callback func(data *TestData)

	// A slice of tab completion candidates to expect.
	TabComplete []string
}

// TestVerb runs table tests for Command implementations.
func TestVerb(t *testing.T, verb string, tests []Test) {
	for n, test := range tests {
		mockAPI := &api.MockAPI{}

		data := &TestData{
			T:       t,
			Api:     mockAPI,
			MockAPI: mockAPI,
			Cmd:     New(verb, mockAPI),
			Test:    test,
		}

		require.NotNil(t, data.Cmd, "Command '%s' is not implemented; it must be added to the `commands.Verb` variable.", verb)

		if data.Test.Init != nil {
			t.Logf("### Initializing data for verb test '%s' number %d", test.Input, n+1)
			data.Test.Init(data)
		}

		t.Logf("### Test %d: '%s'", n+1, test.Input)
		TestCommand(data)
	}
}

// TestCommand runs a single test a for Command implementation.
func TestCommand(data *TestData) {
	reader := strings.NewReader(data.Test.Input)
	scanner := lexer.NewScanner(reader)

	// Parse command
	data.Cmd.SetScanner(scanner)
	err := data.Cmd.Parse()

	// Test success
	if data.Test.Success {
		assert.Nil(data.T, err, "Expected success when parsing '%s'", data.Test.Input)
	} else {
		assert.NotNil(data.T, err, "Expected error when parsing '%s'", data.Test.Input)
	}

	// Test tab completes if not skipped by nil
	if data.Test.TabComplete != nil {
		completes := data.Cmd.TabComplete()
		assert.Equal(data.T, data.Test.TabComplete, completes)
	}

	// Test callback function
	if data.Test.Callback != nil {
		data.Test.Callback(data)
	}
}
