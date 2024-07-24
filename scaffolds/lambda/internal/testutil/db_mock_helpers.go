package testutil

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
)

// MustStructsToRows converts a slice of structs to sqlmock.Rows using reflect.
// It can also be used when only a single struct is needed by wrapping in a slice.
func MustStructsToRows[T any](slice []T) *sqlmock.Rows {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic(fmt.Sprintf("expected a slice but got %T", slice))
	}

	if v.Len() == 0 {
		panic("slice is empty")
	}

	elemType := reflect.TypeOf(slice).Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	if elemType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected a slice of structs but got a slice of %v", elemType.Kind()))
	}

	numFields := elemType.NumField()
	columns := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		colName := elemType.Field(i).Name
		colNameSnake := toSnake(colName)
		columns[i] = colNameSnake
	}

	rows := sqlmock.NewRows(columns)

	for i := 0; i < v.Len(); i++ {
		var values []driver.Value
		elem := v.Index(i)
		for j := 0; j < elem.NumField(); j++ {
			values = append(values, elem.Field(j).Interface())
		}
		rows.AddRow(values...)
	}

	return rows
}

// MustStructToEmptyRow converts a struct into an *sqlmock.Rows object with headers but no rows.
func MustStructToEmptyRow[T any](obj T) *sqlmock.Rows {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected a struct but got %T", obj))
	}

	elemType := v.Type()
	numFields := v.NumField()
	columns := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		colName := elemType.Field(i).Name
		colNameSnake := toSnake(colName)
		columns[i] = colNameSnake
	}

	return sqlmock.NewRows(columns)
}

// toSnake converts PascalCase to snake_case with special handling for abbreviations
func toSnake(camel string) (snake string) {
	var b strings.Builder
	diff := 'a' - 'A'
	l := len(camel)
	for i, v := range camel {
		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		if (i != 0 || i == l-1) &&
			((i > 0 && rune(camel[i-1]) >= 'a') || (i < l-1 && rune(camel[i+1]) >= 'a')) {
			b.WriteRune('_')
		}
		b.WriteRune(v + diff)
	}
	return b.String()
}
