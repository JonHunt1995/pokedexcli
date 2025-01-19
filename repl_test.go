package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " ",
			expected: []string{},
		},
		{
			input:    "cLocKing",
			expected: []string{"clocking"},
		},
		{
			input:    "cool  beans",
			expected: []string{"cool", "beans"},
		},
		{
			input:    "",
			expected: []string{},
		},

		// add more cases here
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("length of actual is %v, which differs from length of expected %v", len(actual), len(c.expected))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("%q didn't match %q", word, expectedWord)
			}
		}
	}
}
