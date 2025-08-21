package main

import (
	"testing"
)

func TestCleanInput(t *testing.T){
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		// add more cases here
	}


	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected){
			t.Errorf("length of the result array does not match the expected. Expecting: %v, Got: %v", len(c.expected), len(actual))
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("resulting elements do not match with the expected elements, Expecting: %v, Got: %v", expectedWord,word) 

			}
		}

	}}
