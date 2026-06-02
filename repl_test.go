package repl_test

import (
	"fmt"
	"testing"
	"github.com/Fladiem/pokedexcli/repl"
)
	

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		input: "  hello  world  "
		expected: []string{"hello", "world"}
	}, {
		input: "I am a Banana"
		expected: []string{"i", "am", "a", "banana"}
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf()
			return 
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf()
				return
			}
		}
	}
}