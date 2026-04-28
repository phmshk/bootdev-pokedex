package cli

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "     hello world     ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU ",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "some interesting    text   ",
			expected: []string{"some", "interesting", "text"},
		},
	}

	for _, c := range cases {
		actual := CleanInput(c.input)
		lenActual := len(actual)
		lenEx := len(c.expected)
		if lenActual != lenEx {
			t.Errorf("expected length actual (%d) be equal length expected (%d), got %d != %d", lenActual, lenEx, lenActual, lenEx)
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("expected words to be equal, got %s != %s", word, expectedWord)
			}
		}
	}
}
