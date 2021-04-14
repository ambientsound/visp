package tabcomplete_test

import (
	"testing"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/commands"
	"github.com/ambientsound/visp/tabcomplete"
	"github.com/stretchr/testify/assert"
)

var tabCompleteTests = []struct {
	input       string
	success     bool
	completions []string
}{
	{"", true, commands.Keys()},
	{"s", true, []string{
		"se",
		"seek",
		"select",
		"set",
		"show",
		"single",
		"sort",
		"stop",
		"style",
	}},
	{"set", true, []string{}},
	{"add ", true, []string{}},
	{"cursor nextOf", true, []string{}},
	{"foobarbaz", false, []string{}},
	{"foobarbaz ", false, []string{}},
	{"$var", false, []string{}},
	{"{foo", false, []string{}},
	{"# bar", false, []string{}},
}

func TestTabComplete(t *testing.T) {
	for n, test := range tabCompleteTests {

		api := api.NewTestAPI()

		t.Logf("### Test %d: '%s'", n+1, test.input)

		clen := len(test.completions)
		tabComplete := tabcomplete.New(test.input, api)
		sentences := make([]string, clen)
		i := 0

		for i < len(sentences) {
			sentence, err := tabComplete.Scan()
			if test.success {
				assert.Nil(t, err, "Expected success when parsing '%s'", test.input)
			} else {
				assert.NotNil(t, err, "Expected error when parsing '%s'", test.input)
			}
			sentences[i] = sentence
			i++
			if i == clen {
				break
			}
		}

		assert.Equal(t, test.completions, sentences)
		assert.Equal(t, clen, tabComplete.Len())
	}
}
