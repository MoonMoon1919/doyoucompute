package doyoucompute

import (
	"fmt"
	"reflect"
	"strings"
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
		header       []string
		row          []string
		errorMessage string
	}{
		{
			name:         "Pass-HasRows",
			numItems:     10,
			header:       []string{"dude", "sweet"},
			row:          []string{"mine", "says"},
			errorMessage: "",
		},
		{
			name:         "Pass-NoRows",
			numItems:     0,
			header:       []string{"dude", "sweet"},
			row:          []string{"mine", "says"},
			errorMessage: "",
		},
		{
			name:         "Fail-RowTooLong",
			numItems:     0,
			header:       []string{"dude", "sweet"},
			row:          []string{"what", "does", "mine", "say"},
			errorMessage: "Row length exceeds number of headers",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Table {
					table := NewTable(tc.header, []TableRow{})

					return table
				},
				func(table *Table) ([]TableRow, error) {
					// Add the header as a row because it will
					// have the correct number of columns
					for range tc.numItems {
						table.AddRow(tc.header...)
					}

					if err := table.AddRow(tc.row...); err != nil {
						return table.Items, err
					}

					return table.Items, nil
				},
				tc.errorMessage,
				func(rows []TableRow, table *Table, t *testing.T) {
					if len(table.Items) != tc.numItems+1 {
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

				table.AddRow("a", "cool", "value")

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

	if table.Identifier() != "" {
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
				func(l *List) ([]Node, error) {
					return l.Children(), nil
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
					list.Append(tc.newItem)

					return list.Items[len(list.Items)-1], nil
				},
				tc.errorMessage,
				func(firstItem Text, list *List, t *testing.T) {
					if len(list.Items) != tc.startingLength+1 {
						t.Errorf("Expected to have %d items, found %d", tc.startingLength+1, len(list.Items))
					}

					if firstItem != Text(tc.newItem) {
						t.Errorf("Expected last item to be %v, found %v", Text(tc.newItem), firstItem)
					}
				},
			)
		})
	}
}

func TestListIdentifier(t *testing.T) {
	list := NewList(BULLET)

	if list.Identifier() != "" {
		t.Error("Expected list identifier to be empty")
	}
}

func TestListType(t *testing.T) {
	list := NewList(NUMBERED)

	if list.Type() != ListType {
		t.Errorf("Expected List.Type() to be %d, got %d", ListType, list.Type())
	}
}

// MARK: Paragraph
func TestParagraphChildren(t *testing.T) {
	tests := []struct {
		name         string
		numChildren  int
		errorMessage string
	}{
		{
			name:         "Pass-SomeChildren",
			numChildren:  10,
			errorMessage: "",
		},
		{
			name:         "Pass-NoChildren",
			numChildren:  0,
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Paragraph {
					paragraph := NewParagraph()

					for idx := range tc.numChildren {
						paragraph.Text(fmt.Sprintf("Text item %d", idx))
					}

					return paragraph
				},
				func(p *Paragraph) ([]Node, error) {
					return p.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, paragraph *Paragraph, t *testing.T) {
					if len(paragraph.Items) != tc.numChildren {
						t.Errorf("Expected %d children, found %d", tc.numChildren, len(paragraph.Items))
					}

					if !reflect.DeepEqual(result, paragraph.Items) {
						t.Errorf("Expected result %v, got %v", paragraph.Items, result)
					}
				},
			)
		})
	}
}

func TestParagraphText(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Pass-NoExistingItems",
			body:          "Sik text",
			existingItems: 0,
			errorMessage:  "",
		},
		{
			name:          "Pass-SomeItems",
			body:          "Sik text",
			existingItems: 10,
			errorMessage:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Paragraph {
					para := NewParagraph()

					for idx := range tc.existingItems {
						if idx%2 == 0 {
							para.Text(fmt.Sprintf("Text idx %d", idx))
						} else {
							para.Code(fmt.Sprintf("Code idx %d", idx))
						}
					}

					return para
				},
				func(p *Paragraph) ([]Node, error) {
					p = p.Text(tc.body)

					return p.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, paragraph *Paragraph, t *testing.T) {
					if len(paragraph.Items) < 1 {
						t.Errorf("Expected at least 1 child, got %d", len(paragraph.Children()))
					}

					lastItem := result[len(paragraph.Children())-1]

					if lastItem.Type() != TextType {
						t.Errorf("Expected error type to be %d got %d", TextType, lastItem.Type())
					}

					materializedContent, err := lastItem.(Text).Materialize()
					if err != nil {
						t.Errorf("Got unexpected error materializing content %s", err.Error())
					}
					if materializedContent.Content != tc.body {
						t.Errorf("Expected content %s, got %s", materializedContent.Content, tc.body)
					}
				},
			)
		})
	}
}
func TestParagraphCode(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Pass-NoExistingItems",
			body:          "Sik text",
			existingItems: 0,
			errorMessage:  "",
		},
		{
			name:          "Pass-SomeItems",
			body:          "Sik text",
			existingItems: 10,
			errorMessage:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Paragraph {
					para := NewParagraph()

					for idx := range tc.existingItems {
						if idx%2 == 0 {
							para.Text(fmt.Sprintf("Text idx %d", idx))
						} else {
							para.Code(fmt.Sprintf("Code idx %d", idx))
						}
					}

					return para
				},
				func(p *Paragraph) ([]Node, error) {
					p = p.Code(tc.body)

					return p.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, paragraph *Paragraph, t *testing.T) {
					if len(paragraph.Items) < 1 {
						t.Errorf("Expected at least 1 child, got %d", len(paragraph.Children()))
					}

					lastItem := result[len(paragraph.Children())-1]

					if lastItem.Type() != CodeType {
						t.Errorf("Expected error type to be %d got %d", CodeType, lastItem.Type())
					}

					materializedContent, err := lastItem.(Code).Materialize()
					if err != nil {
						t.Errorf("Got unexpected error materializing content %s", err.Error())
					}
					if materializedContent.Content != tc.body {
						t.Errorf("Expected content %s, got %s", materializedContent.Content, tc.body)
					}
				},
			)
		})
	}
}
func TestParagraphLink(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		link          string
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Pass-NoExistingItems",
			body:          "Sik text",
			link:          "https://example.com",
			existingItems: 0,
			errorMessage:  "",
		},
		{
			name:          "Pass-SomeItems",
			body:          "Sik text",
			link:          "https://google.com",
			existingItems: 10,
			errorMessage:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Paragraph {
					para := NewParagraph()

					for idx := range tc.existingItems {
						if idx%2 == 0 {
							para.Text(fmt.Sprintf("Text idx %d", idx))
						} else {
							para.Code(fmt.Sprintf("Code idx %d", idx))
						}
					}

					return para
				},
				func(p *Paragraph) ([]Node, error) {
					p = p.Link(tc.body, tc.link)

					return p.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, paragraph *Paragraph, t *testing.T) {
					if len(paragraph.Items) < 1 {
						t.Errorf("Expected at least 1 child, got %d", len(paragraph.Children()))
					}

					lastItem := result[len(paragraph.Children())-1]

					if lastItem.Type() != LinkType {
						t.Errorf("Expected error type to be %d got %d", LinkType, lastItem.Type())
					}

					materializedContent, err := lastItem.(Link).Materialize()
					if err != nil {
						t.Errorf("Got unexpected error materializing content %s", err.Error())
					}
					if materializedContent.Content != tc.body {
						t.Errorf("Expected content %s, got %s", tc.body, materializedContent.Content)
					}
					if materializedContent.Metadata["Url"] != tc.link {
						t.Errorf("Expected content %s, got %s", tc.link, materializedContent.Metadata["Url"])
					}
				},
			)
		})
	}
}

// MARK: Section
func TestSectionChildren(t *testing.T) {
	tests := []struct {
		name         string
		sectionName  string
		numChildren  int
		errorMessage string
	}{
		{
			name:         "Pass-SomeChildren",
			sectionName:  "Cool Section",
			numChildren:  10,
			errorMessage: "",
		},
		{
			name:         "Pass-NoChildren",
			sectionName:  "Very Cool Section",
			numChildren:  0,
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Section {
					section := NewSection(tc.sectionName)

					for idx := range tc.numChildren {
						section.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &section
				},
				func(p *Section) ([]Node, error) {
					return p.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, paragraph *Section, t *testing.T) {
					if len(paragraph.Content) != tc.numChildren {
						t.Errorf("Expected %d children, found %d", tc.numChildren, len(paragraph.Content))
					}

					if !reflect.DeepEqual(result, paragraph.Content) {
						t.Errorf("Expected result %v, got %v", paragraph.Content, result)
					}
				},
			)
		})
	}
}

func TestSectionAddIntro(t *testing.T) {
	tests := []struct {
		name          string
		introContent  string
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			introContent:  "Cool intro",
			errorMessage:  "",
			existingItems: 0,
		},
		{
			name:          "Passing-SomeItems",
			introContent:  "Cool intro",
			errorMessage:  "",
			existingItems: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Section {
					section := NewSection("test")

					for idx := range tc.existingItems {
						section.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &section
				},
				func(s *Section) ([]Node, error) {
					introText := NewParagraph().Text(tc.introContent)

					s.AddIntro(introText)

					return s.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, section *Section, t *testing.T) {
					if len(section.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(section.Content))
					}

					// Get the first item to ensure it's a paragraph
					firstItem := result[0]
					if firstItem.Type() != ParagraphType {
						t.Errorf("Expected error type to be %d got %d", ParagraphType, firstItem.Type())
					}

					// Check the content
					content := firstItem.(*Paragraph).Children()

					if content[0] != Text(tc.introContent) {
						t.Errorf("Expected content %v, got %v", Text(tc.introContent), content[0])
					}
				},
			)
		})
	}
}

func TestSectionWriteIntro(t *testing.T) {
	tests := []struct {
		name          string
		introContent  string
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			introContent:  "Cool intro",
			errorMessage:  "",
			existingItems: 0,
		},
		{
			name:          "Passing-SomeItems",
			introContent:  "Cool intro",
			errorMessage:  "",
			existingItems: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Section {
					section := NewSection("test")

					for idx := range tc.existingItems {
						section.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &section
				},
				func(s *Section) ([]Node, error) {

					s.WriteIntro().Text(tc.introContent)

					return s.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, section *Section, t *testing.T) {
					if len(section.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(section.Content))
					}

					// Get the first item to ensure it's a paragraph
					firstItem := result[0]
					if firstItem.Type() != ParagraphType {
						t.Errorf("Expected error type to be %d got %d", ParagraphType, firstItem.Type())
					}

					// Check the content
					content := firstItem.(*Paragraph).Children()

					if content[0] != Text(tc.introContent) {
						t.Errorf("Expected content %v, got %v", Text(tc.introContent), content[0])
					}
				},
			)
		})
	}
}

func TestSectionAddSection(t *testing.T) {
	tests := []struct {
		name          string
		sectionName   string
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			sectionName:   "Cool section",
			errorMessage:  "",
			existingItems: 0,
		},
		{
			name:          "Passing-SomeItems",
			sectionName:   "Cool section",
			errorMessage:  "",
			existingItems: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Section {
					section := NewSection("test")

					for idx := range tc.existingItems {
						section.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &section
				},
				func(s *Section) ([]Node, error) {
					s.AddSection(NewSection(tc.sectionName))

					return s.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, section *Section, t *testing.T) {
					if len(section.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(section.Content))
					}

					// Get the last item to ensure it's a section
					lastItem := result[len(section.Children())-1]
					if lastItem.Type() != SectionType {
						t.Errorf("Expected error type to be %d got %d", SectionType, lastItem.Type())
					}

					if lastItem.(Section).Name != tc.sectionName {
						t.Errorf("Expected last section to have name %s, got %s", tc.sectionName, lastItem.(Section).Name)
					}
				},
			)
		})
	}
}

func TestSectionCreateSection(t *testing.T) {
	tests := []struct {
		name          string
		sectionName   string
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			sectionName:   "Cool section",
			errorMessage:  "",
			existingItems: 0,
		},
		{
			name:          "Passing-SomeItems",
			sectionName:   "Cool section",
			errorMessage:  "",
			existingItems: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Section {
					section := NewSection("test")

					for idx := range tc.existingItems {
						section.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &section
				},
				func(s *Section) ([]Node, error) {
					s.CreateSection(tc.sectionName)

					return s.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, section *Section, t *testing.T) {
					if len(section.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(section.Content))
					}

					// Get the last item to ensure it's a section
					lastItem := result[len(section.Children())-1]
					if lastItem.Type() != SectionType {
						t.Errorf("Expected error type to be %d got %d", SectionType, lastItem.Type())
					}

					if lastItem.(*Section).Name != tc.sectionName {
						t.Errorf("Expected last section to have name %s, got %s", tc.sectionName, lastItem.(Section).Name)
					}
				},
			)
		})
	}
}

func TestSectionWriteParagraph(t *testing.T) {
	tests := []struct {
		name             string
		paragraphContent string
		existingItems    int
		errorMessage     string
	}{
		{
			name:             "Passing-NoItems",
			paragraphContent: "Cool information",
			errorMessage:     "",
			existingItems:    0,
		},
		{
			name:             "Passing-SomeItems",
			paragraphContent: "Cool information",
			errorMessage:     "",
			existingItems:    10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Section {
					section := NewSection("test")

					for idx := range tc.existingItems {
						section.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &section
				},
				func(s *Section) ([]Node, error) {
					s.WriteParagraph().Text(tc.paragraphContent)

					return s.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, section *Section, t *testing.T) {
					if len(section.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(section.Content))
					}

					// Get the last item to ensure it's a section
					lastItem := result[len(section.Children())-1]
					if lastItem.Type() != ParagraphType {
						t.Errorf("Expected error type to be %d got %d", ParagraphType, lastItem.Type())
					}

					// Check if the first item in the section is a paragraph
					firstItem := lastItem.(*Paragraph).Children()[0]
					if firstItem.Type() != TextType {
						t.Errorf("Expected error type to be %d got %d", TextType, firstItem.Type())
					}

					materializedContent, _ := firstItem.(Text).Materialize()
					if materializedContent.Content != tc.paragraphContent {
						t.Errorf("Expected content %v, got %v", tc.paragraphContent, materializedContent.Content)
					}
				},
			)
		})
	}
}

func TestSectionAddTable(t *testing.T) {
	tests := []struct {
		name          string
		headers       []string
		rows          []TableRow
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			headers:       []string{"cool", "headers"},
			rows:          []TableRow{{Values: []string{"cool", "value"}}},
			errorMessage:  "",
			existingItems: 0,
		},
		{
			name:          "Passing-SomeItems",
			headers:       []string{"cool", "headers"},
			rows:          []TableRow{{Values: []string{"cool", "value"}}},
			errorMessage:  "",
			existingItems: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Section {
					section := NewSection("test")

					for idx := range tc.existingItems {
						section.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &section
				},
				func(s *Section) ([]Node, error) {
					s.AddTable(tc.headers, tc.rows)

					return s.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, section *Section, t *testing.T) {
					if len(section.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(section.Content))
					}

					// Get the last item to ensure it's a section
					lastItem := result[len(section.Children())-1]
					if lastItem.Type() != TableType {
						t.Errorf("Expected error type to be %d got %d", TableType, lastItem.Type())
					}
				},
			)
		})
	}
}

func TestSectionCreateTable(t *testing.T) {
	tests := []struct {
		name          string
		headers       []string
		rows          [][]string
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			headers:       []string{"cool", "headers"},
			rows:          [][]string{{"cool", "value"}},
			errorMessage:  "",
			existingItems: 0,
		},
		{
			name:          "Passing-SomeItems",
			headers:       []string{"cool", "headers"},
			rows:          [][]string{{"cool", "value"}},
			errorMessage:  "",
			existingItems: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Section {
					section := NewSection("test")

					for idx := range tc.existingItems {
						section.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &section
				},
				func(s *Section) ([]Node, error) {
					table := s.CreateTable(tc.headers)

					for _, row := range tc.rows {
						table.AddRow(row...)
					}

					return s.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, section *Section, t *testing.T) {
					if len(section.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(section.Content))
					}

					// Get the last item to ensure it's a section
					lastItem := result[len(section.Children())-1]
					if lastItem.Type() != TableType {
						t.Errorf("Expected error type to be %d got %d", TableType, lastItem.Type())
					}
				},
			)
		})
	}
}

func TestSectionAddList(t *testing.T) {
	tests := []struct {
		name          string
		listItems     []Text
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			listItems:     []Text{"Cool item"},
			existingItems: 0,
			errorMessage:  "",
		},
		{
			name:          "Passing-SomeItems",
			listItems:     []Text{"Cool item"},
			existingItems: 10,
			errorMessage:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Section {
					section := NewSection("test")

					for idx := range tc.existingItems {
						section.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &section
				},
				func(s *Section) ([]Node, error) {
					s.AddList(BULLET, tc.listItems)

					return s.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, section *Section, t *testing.T) {
					if len(section.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(section.Content))
					}

					// Get the last item to ensure it's a section
					lastItem := result[len(section.Children())-1]
					if lastItem.Type() != ListType {
						t.Errorf("Expected error type to be %d got %d", ListType, lastItem.Type())
					}
				},
			)
		})
	}
}

func TestSectionCreateList(t *testing.T) {
	tests := []struct {
		name          string
		listItems     []Text
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			listItems:     []Text{"Cool item"},
			existingItems: 0,
			errorMessage:  "",
		},
		{
			name:          "Passing-SomeItems",
			listItems:     []Text{"Cool item"},
			existingItems: 10,
			errorMessage:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Section {
					section := NewSection("test")

					for idx := range tc.existingItems {
						section.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &section
				},
				func(s *Section) ([]Node, error) {
					list := s.CreateList(BULLET)

					for _, item := range tc.listItems {
						list.Append(string(item))
					}

					return s.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, section *Section, t *testing.T) {
					if len(section.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(section.Content))
					}

					// Get the last item to ensure it's a section
					lastItem := result[len(section.Children())-1]
					if lastItem.Type() != ListType {
						t.Errorf("Expected error type to be %d got %d", ListType, lastItem.Type())
					}

					if len(lastItem.(*List).Children()) != len(tc.listItems) {
						t.Errorf("Expected list to have %d children, found %d", len(tc.listItems), len(lastItem.(*List).Children()))
					}
				},
			)
		})
	}
}

func TestSectionWriteCodeBlock(t *testing.T) {
	tests := []struct {
		name          string
		shell         string
		content       []string
		executable    CodeBlockExecType
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			shell:         "sh",
			content:       []string{"echo", "hello", "world"},
			executable:    Exec,
			existingItems: 0,
			errorMessage:  "",
		},
		{
			name:          "Passing-SomeItems",
			shell:         "sh",
			content:       []string{"echo", "hello", "world"},
			executable:    Static,
			existingItems: 10,
			errorMessage:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Section {
					section := NewSection("test")

					for idx := range tc.existingItems {
						section.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &section
				},
				func(s *Section) ([]Node, error) {
					s.WriteCodeBlock(tc.shell, tc.content, tc.executable)

					return s.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, section *Section, t *testing.T) {
					if len(section.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(section.Content))
					}

					// Get the last item to ensure it's a section
					lastItem := result[len(section.Children())-1]

					var blockType ContentType
					if tc.executable == Exec {
						blockType = ExecutableType
					} else {
						blockType = CodeBlockType
					}

					if lastItem.Type() != blockType {
						t.Errorf("Expected error type to be %d got %d", blockType, lastItem.Type())
					}
				},
			)
		})
	}
}

func TestSectionWriteBlockQuote(t *testing.T) {
	tests := []struct {
		name          string
		quote         string
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			quote:         "Cool quote",
			existingItems: 0,
			errorMessage:  "",
		},
		{
			name:          "Passing-SomeItems",
			quote:         "Cool quote",
			existingItems: 10,
			errorMessage:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Section {
					section := NewSection("test")

					for idx := range tc.existingItems {
						section.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &section
				},
				func(s *Section) ([]Node, error) {
					s.WriteBlockQuote(tc.quote)

					return s.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, section *Section, t *testing.T) {
					if len(section.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(section.Content))
					}

					// Get the last item to ensure it's a section
					lastItem := result[len(section.Children())-1]
					if lastItem.Type() != BlockQuoteType {
						t.Errorf("Expected error type to be %d got %d", BlockQuoteType, lastItem.Type())
					}
				},
			)
		})
	}
}

func TestSectionWriteRemoteContent(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			content:       "Doot doot",
			existingItems: 0,
			errorMessage:  "",
		},
		{
			name:          "Passing-SomeItems",
			content:       "Doot doot",
			existingItems: 10,
			errorMessage:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Section {
					section := NewSection("test")

					for idx := range tc.existingItems {
						section.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &section
				},
				func(s *Section) ([]Node, error) {
					s.WriteRemoteContent(Remote{Reader: strings.NewReader(tc.content)})

					return s.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, section *Section, t *testing.T) {
					if len(section.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(section.Content))
					}

					// Get the last item to ensure it's a section
					lastItem := result[len(section.Children())-1]
					if lastItem.Type() != RemoteType {
						t.Errorf("Expected error type to be %d got %d", RemoteType, lastItem.Type())
					}

					materialized, err := lastItem.(Remote).Materialize()
					if err != nil {
						t.Errorf("Unexpected error %s", err.Error())
					}

					if materialized.Content != tc.content {
						t.Errorf("Got content %s, expected %s", materialized.Content, tc.content)
					}
				},
			)
		})
	}
}

func TestSectionWriteComment(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		existingItems int
		errorMessage  string
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Section {
					section := NewSection("test")

					for idx := range tc.existingItems {
						section.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &section
				},
				func(s *Section) ([]Node, error) {
					s.WriteComment(tc.content)

					return s.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, section *Section, t *testing.T) {
					if len(section.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(section.Content))
					}

					// Get the last item to ensure it's a section
					lastItem := result[len(section.Children())-1]
					if lastItem.Type() != CommentType {
						t.Errorf("Expected error type to be %d got %d", CommentType, lastItem.Type())
					}

					materialized, err := lastItem.(Comment).Materialize()
					if err != nil {
						t.Errorf("Unexpected error %s", err.Error())
					}

					if materialized.Content != tc.content {
						t.Errorf("Got content %s, expected %s", materialized.Content, tc.content)
					}
				},
			)
		})
	}
}

// MARK: Document
func TestDocumentAddIntro(t *testing.T) {
	tests := []struct {
		name          string
		introContent  string
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			introContent:  "Cool intro",
			errorMessage:  "",
			existingItems: 0,
		},
		{
			name:          "Passing-SomeItems",
			introContent:  "Cool intro",
			errorMessage:  "",
			existingItems: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Document {
					document, _ := NewDocument("test")

					for idx := range tc.existingItems {
						document.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &document
				},
				func(d *Document) ([]Node, error) {
					introText := NewParagraph().Text(tc.introContent)

					d.AddIntro(introText)

					return d.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, document *Document, t *testing.T) {
					if len(document.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(document.Content))
					}

					// Get the first item to ensure it's a paragraph
					firstItem := result[0]
					if firstItem.Type() != ParagraphType {
						t.Errorf("Expected error type to be %d got %d", ParagraphType, firstItem.Type())
					}

					// Check the content
					content := firstItem.(*Paragraph).Children()

					if content[0] != Text(tc.introContent) {
						t.Errorf("Expected content %v, got %v", Text(tc.introContent), content[0])
					}
				},
			)
		})
	}
}

func TestDocumentWriteIntro(t *testing.T) {
	tests := []struct {
		name          string
		introContent  string
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			introContent:  "Cool intro",
			errorMessage:  "",
			existingItems: 0,
		},
		{
			name:          "Passing-SomeItems",
			introContent:  "Cool intro",
			errorMessage:  "",
			existingItems: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Document {
					document, _ := NewDocument("test")

					for idx := range tc.existingItems {
						document.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &document
				},
				func(d *Document) ([]Node, error) {

					d.WriteIntro().Text(tc.introContent)

					return d.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, document *Document, t *testing.T) {
					if len(document.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(document.Content))
					}

					// Get the first item to ensure it's a paragraph
					firstItem := result[0]
					if firstItem.Type() != ParagraphType {
						t.Errorf("Expected error type to be %d got %d", ParagraphType, firstItem.Type())
					}

					// Check the content
					content := firstItem.(*Paragraph).Children()

					if content[0] != Text(tc.introContent) {
						t.Errorf("Expected content %v, got %v", Text(tc.introContent), content[0])
					}
				},
			)
		})
	}
}

func TestDocumentAddSection(t *testing.T) {
	tests := []struct {
		name          string
		sectionName   string
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			sectionName:   "Cool section",
			errorMessage:  "",
			existingItems: 0,
		},
		{
			name:          "Passing-SomeItems",
			sectionName:   "Cool section",
			errorMessage:  "",
			existingItems: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Document {
					document, _ := NewDocument("test")

					for idx := range tc.existingItems {
						document.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &document
				},
				func(d *Document) ([]Node, error) {
					d.AddSection(NewSection(tc.sectionName))

					return d.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, document *Document, t *testing.T) {
					if len(document.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(document.Content))
					}

					// Get the last item to ensure it's a section
					lastItem := result[len(document.Children())-1]
					if lastItem.Type() != SectionType {
						t.Errorf("Expected error type to be %d got %d", SectionType, lastItem.Type())
					}

					if lastItem.(Section).Name != tc.sectionName {
						t.Errorf("Expected last section to have name %s, got %s", tc.sectionName, lastItem.(Section).Name)
					}
				},
			)
		})
	}
}

func TestDocumentCreateSection(t *testing.T) {
	tests := []struct {
		name          string
		sectionName   string
		existingItems int
		errorMessage  string
	}{
		{
			name:          "Passing-NoItems",
			sectionName:   "Cool section",
			errorMessage:  "",
			existingItems: 0,
		},
		{
			name:          "Passing-SomeItems",
			sectionName:   "Cool section",
			errorMessage:  "",
			existingItems: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Document {
					document, _ := NewDocument("test")

					for idx := range tc.existingItems {
						document.AddSection(NewSection(fmt.Sprintf("Section %d", idx)))
					}

					return &document
				},
				func(d *Document) ([]Node, error) {
					d.CreateSection(tc.sectionName)

					return d.Children(), nil
				},
				tc.errorMessage,
				func(result []Node, document *Document, t *testing.T) {
					if len(document.Content) != tc.existingItems+1 {
						t.Errorf("Expected %d children, found %d", tc.existingItems+1, len(document.Content))
					}

					// Get the last item to ensure it's a section
					lastItem := result[len(document.Children())-1]
					if lastItem.Type() != SectionType {
						t.Errorf("Expected error type to be %d got %d", SectionType, lastItem.Type())
					}

					if lastItem.(*Section).Name != tc.sectionName {
						t.Errorf("Expected last section to have name %s, got %s", tc.sectionName, lastItem.(Section).Name)
					}
				},
			)
		})
	}
}

func TestDocumentAddFrontmatter(t *testing.T) {
	tests := []struct {
		name         string
		content      Frontmatter
		errorMessage string
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOperation(
				t,
				func() *Document {
					document, _ := NewDocument("test")

					return &document
				},
				func(d *Document) (Frontmatter, error) {
					d.AddFrontmatter(tc.content)

					return d.Frontmatter, nil
				},
				tc.errorMessage,
				func(result Frontmatter, document *Document, t *testing.T) {
					if document.Frontmatter.Data == nil {
						t.Errorf("Frontmatter is empty, expected to have content")
					}
				},
			)
		})
	}
}
