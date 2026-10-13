package repl

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "Hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Pika Someguy      hi lol",
			expected: []string{"pika", "someguy", "hi", "lol"},
		},
	}

	for idx, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Fatalf("test[%d] - Length of both arrays is different, expected=%d but got=%d", idx, len(c.expected), len(actual))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Fatalf("test[%d] - Word no. %d don't match, expected=%q but got=%q", idx, i, expectedWord, word)
			}

		}
	}
}
