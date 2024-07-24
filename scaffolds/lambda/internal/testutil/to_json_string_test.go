package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToJSONString(t *testing.T) {
	type testCase struct {
		input       any
		expected    string
		shouldPanic bool
	}

	tests := map[string]testCase{
		"simple struct": {
			input:    struct{ Name string }{Name: "John"},
			expected: `{"Name":"John"}`,
		},
		"map": {
			input:    map[string]int{"one": 1, "two": 2},
			expected: `{"one":1,"two":2}`,
		},
		"slice": {
			input:    []string{"apple", "banana"},
			expected: `["apple","banana"]`,
		},
		"integer": {
			input:    42,
			expected: `42`,
		},
		"float": {
			input:    3.14,
			expected: `3.14`,
		},
		"boolean": {
			input:    true,
			expected: `true`,
		},
		"panic on channel": {
			input:       make(chan int),
			shouldPanic: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.shouldPanic {
				assert.Panics(t, func() { ToJSONString(tc.input) })
				return
			}

			result := ToJSONString(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
