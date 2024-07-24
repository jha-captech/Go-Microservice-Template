package testutil

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMustStructsToRows(t *testing.T) {
	type TestStruct struct {
		ID   int
		Name string
	}

	tests := map[string]struct {
		input        []TestStruct
		expectedRows *sqlmock.Rows
		willPanic    bool
	}{
		"valid slice of structs": {
			input: []TestStruct{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			expectedRows: sqlmock.NewRows([]string{"id", "name"}).AddRows([][]driver.Value{
				{1, "Alice"},
				{2, "Bob"},
			}...),
			willPanic: false,
		},
		"empty slice": {
			input:        []TestStruct{},
			expectedRows: sqlmock.NewRows([]string{}),
			willPanic:    true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.willPanic {
				assert.Panics(t, func() { MustStructsToRows(tc.input) })
				return
			}

			rows := MustStructsToRows(tc.input)
			assert.Equal(t, tc.expectedRows, rows)
		})
	}
}
