package doyoucompute

import (
	"fmt"
	"reflect"
	"testing"
)

func testOperation[T Structurer, R any](
	t *testing.T,
	setup func() *T,
	operation func(*T) (R, error),
	errorMessage string,
	comparisonFunc func(R, *T, *testing.T),
) {
	structUnderTest := setup()

	res, err := operation(structUnderTest)

	checkErrors(errorMessage, err, t)
	if errorMessage != "" {
		return
	}

	comparisonFunc(res, structUnderTest, t)
}

func TestTableAddRow(t *testing.T) {
	tests := []struct {
		name         string
		numItems     int
		errorMessage string
	}{
		{
			name:         "Pass-HasRows",
			numItems:     10,
			errorMessage: "",
		},
		{
			name:         "Pass-NoRows",
			numItems:     0,
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Table {
					table := NewTable([]string{"cool", "header"}, []TableRow{})

					return table
				},
				func(table *Table) ([]TableRow, error) {
					for idx := range tc.numItems {
						table.AddRow(TableRow{Values: []string{"sweet", fmt.Sprintf("%d", idx)}})
					}

					return table.Items, nil
				},
				tc.errorMessage,
				func(rows []TableRow, table *Table, t *testing.T) {
					if len(table.Items) != tc.numItems {
						t.Errorf("Expected %d items, found %d", tc.numItems, len(table.Items))
					}
				},
			)
		})
	}
}

func TestTableChildren(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func() *Table
		validationFunc func([]Node, *Table, *testing.T)
		errorMessage   string
	}{
		{
			name: "Pass-HasChildren",
			setupFunc: func() *Table {
				table := NewTable([]string{"some", "cool", "headers"}, []TableRow{})

				table.AddRow(TableRow{
					Values: []string{"a", "cool", "value"},
				})

				return table
			},
			validationFunc: func(result []Node, table *Table, t *testing.T) {
				expected := make([]Node, len(table.Items))
				for idx, row := range table.Items {
					expected[idx] = row
				}

				if !reflect.DeepEqual(result, expected) {
					t.Errorf("Expected result %v, got %v", expected, result)
				}
			},
		},
		{
			name: "Pass-NoChildren",
			setupFunc: func() *Table {
				return NewTable([]string{"some", "cool", "headers"}, []TableRow{})
			},
			validationFunc: func(result []Node, table *Table, t *testing.T) {
				if len(result) != 0 {
					t.Errorf("expected empty result")
				}

				expected := make([]Node, len(table.Items))
				for idx, row := range table.Items {
					expected[idx] = row
				}

				if !reflect.DeepEqual(result, expected) {
					t.Errorf("Expected result %v, got %v", expected, result)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				tc.setupFunc,
				func(t *Table) ([]Node, error) {
					return t.Children(), nil
				},
				tc.errorMessage,
				tc.validationFunc,
			)
		})
	}
}

func TestTableType(t *testing.T) {
	table := NewTable([]string{"cool", "table"}, []TableRow{{Values: []string{"sweet", "value"}}})

	if table.Type() != TableType {
		t.Errorf("Expected Type() to return %d, got %d", TableType, table.Type())
	}
}

func TestTableIdentifier(t *testing.T) {
	table := NewTable([]string{"cool", "table"}, []TableRow{{Values: []string{"sweet", "value"}}})

	if table.Identifer() != "" {
		t.Errorf("Expected no table identifiers")
	}
}
