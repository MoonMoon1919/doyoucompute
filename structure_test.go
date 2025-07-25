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

// MARK: Table
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

// MARK: List
func TestListChildren(t *testing.T) {
	tests := []struct {
		name         string
		numChildren  int
		errorMessage string
	}{
		{
			name:        "Pass-SomeChildren",
			numChildren: 10,
		},
		{
			name:        "Pass-NoChildren",
			numChildren: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *List {
					list := NewList(BULLET)

					for idx := range tc.numChildren {
						list.Append(fmt.Sprintf("%d", idx))
					}

					return list
				},
				func(t *List) ([]Node, error) {
					return t.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, list *List, t *testing.T) {
					if len(list.Items) != tc.numChildren {
						t.Errorf("Expected %d children, found %d", tc.numChildren, len(list.Items))
					}

					expected := make([]Node, len(list.Items))
					for idx, row := range list.Items {
						expected[idx] = row
					}

					if !reflect.DeepEqual(result, expected) {
						t.Errorf("Expected result %v, got %v", expected, result)
					}
				},
			)
		})
	}
}

func TestListPush(t *testing.T) {
	tests := []struct {
		name           string
		startingLength int
		newItem        string
		errorMessage   string
	}{
		{
			name:           "Pass-ExistingItems",
			startingLength: 10,
			newItem:        "new stuff",
			errorMessage:   "",
		},
		{
			name:           "Pass-NoItems",
			startingLength: 0,
			newItem:        "new stuff",
			errorMessage:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *List {
					initialItems := make([]Text, tc.startingLength)

					for idx := range tc.startingLength {
						initialItems[idx] = Text(fmt.Sprintf("item-%d", idx))
					}

					list := List{TypeOfList: BULLET, Items: initialItems}

					return &list
				},
				func(list *List) (Text, error) {
					list.Push(tc.newItem)

					return list.Items[0], nil
				},
				tc.errorMessage,
				func(firstItem Text, list *List, t *testing.T) {
					if len(list.Items) != tc.startingLength+1 {
						t.Errorf("Expected to have %d items, found %d", tc.startingLength+1, len(list.Items))
					}

					if firstItem != Text(tc.newItem) {
						t.Errorf("Expected first item to be %v, found %v", Text(tc.newItem), firstItem)
					}
				},
			)
		})
	}
}

func TestListAppend(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Pass",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}

func TestListIdentifer(t *testing.T) {
	list := NewList(BULLET)

	if list.Identifer() != "" {
		t.Error("Expected list identifier to be empty")
	}
}

func TestListType(t *testing.T) {
	list := NewList(NUMBERED)

	if list.Type() != ListType {
		t.Errorf("Expected List.Type() to be %d, got %d", ListType, list.Type())
	}
}
