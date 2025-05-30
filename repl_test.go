package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input		string
		expected	[]string
	}{
		{
			input: "  hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input: "Hello World",
			expected: []string{"hello", "world"},
		},
		{
			input: "HELLO WORLD   ",
			expected: []string{"hello", "world"},
		},

	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("len actual %v != len c.expected %v", actual, c.expected)
			return
		}

		for i := range actual {
			word := actual[i]
			expected := c.expected[i]
			if word != expected {
				t.Errorf("actual word %v does not match expected word %v", word, expected)
			}
		}

	}

}
