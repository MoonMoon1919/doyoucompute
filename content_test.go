package doyoucompute

import (
	"reflect"
	"strings"
	"testing"
)

func testMaterialize(
	t *testing.T,
	setup func() Contenter,
	errorMessage string,
	comparisonFunc func(m MaterializedContent, t *testing.T),
) {
	itemUnderTest := setup()

	materializedContent, err := itemUnderTest.Materialize()

	checkErrors(errorMessage, err, t)
	if errorMessage != "" {
		return
	}

	comparisonFunc(materializedContent, t)
}

func TestHeaderMaterialize(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		errorMessage string
	}{
		{
			name:         "Passing",
			content:      "doot doot doot doot",
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMaterialize(
				t,
				func() Contenter {
					return Header{Content: tc.content}
				},
				tc.errorMessage,
				func(m MaterializedContent, t *testing.T) {
					if m.Type != HeaderType {
						t.Errorf("Expected Type to be %d, got %d", HeaderType, m.Type)
					}

					if m.Content != tc.content {
						t.Errorf("Expected Content to be %s, got %s", tc.content, m.Content)
					}
				},
			)
		})
	}
}

func TestTextMaterialize(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		errorMessage string
	}{
		{
			name:         "Passing",
			content:      "Lemons",
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMaterialize(
				t,
				func() Contenter {
					return Text(tc.content)
				},
				tc.errorMessage,
				func(m MaterializedContent, t *testing.T) {
					if m.Type != TextType {
						t.Errorf("Expected Type to be %d, got %d", TextType, m.Type)
					}

					if m.Content != tc.content {
						t.Errorf("Expected Content to be %s, got %s", tc.content, m.Content)
					}
				},
			)
		})
	}
}

func TestLinkMaterialize(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		url          string
		errorMessage string
	}{
		{
			name:         "Passing",
			content:      "Example",
			url:          "https://example.com",
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMaterialize(
				t,
				func() Contenter {
					return Link{Text: tc.content, Url: tc.url}
				},
				tc.errorMessage,
				func(m MaterializedContent, t *testing.T) {
					if m.Type != LinkType {
						t.Errorf("Expected Type to be %d, got %d", LinkType, m.Type)
					}

					if m.Content != tc.content {
						t.Errorf("Expected Content to be %s, got %s", tc.content, m.Content)
					}

					if val, ok := m.Metadata["Url"]; ok {
						if val != tc.url {
							t.Errorf("Expected Url to be %s, got %s", tc.url, val)
						}
					} else {
						t.Errorf("Did not find url")
					}
				},
			)
		})
	}
}

func TestCodeMaterialize(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		errorMessage string
	}{
		{
			name:         "Passing",
			content:      "npm i",
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMaterialize(
				t,
				func() Contenter {
					return Code(tc.content)
				},
				tc.errorMessage,
				func(m MaterializedContent, t *testing.T) {
					if m.Type != CodeType {
						t.Errorf("Expected Type to be %d, got %d", CodeType, m.Type)
					}

					if m.Content != tc.content {
						t.Errorf("Expected Content to be %s, got %s", tc.content, m.Content)
					}
				},
			)
		})
	}
}

func TestCodeBlockMaterialize(t *testing.T) {
	tests := []struct {
		name         string
		blockType    string
		content      []string
		errorMessage string
	}{
		{
			name:         "Passing",
			blockType:    "sh",
			content:      []string{"go", "vet"},
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMaterialize(
				t,
				func() Contenter {
					return CodeBlock{BlockType: tc.blockType, Cmd: tc.content}
				},
				tc.errorMessage,
				func(m MaterializedContent, t *testing.T) {
					if m.Type != CodeBlockType {
						t.Errorf("Expected Type to be %d, got %d", CodeBlockType, m.Type)
					}

					if m.Content != strings.Join(tc.content, " ") {
						t.Errorf("Expected content to be %s, got %s", strings.Join(tc.content, " "), m.Content)
					}

					if val, ok := m.Metadata["BlockType"]; ok {
						if val != tc.blockType {
							t.Errorf("Expected BlockType to be %s, got %s", tc.blockType, val)
						}
					} else {
						t.Errorf("Did not find BlockType")
					}
				},
			)
		})
	}
}

func TestExecutableMaterialize(t *testing.T) {
	tests := []struct {
		name         string
		blockType    string
		content      []string
		errorMessage string
	}{
		{
			name:         "Passing",
			blockType:    "sh",
			content:      []string{"go", "vet"},
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMaterialize(
				t,
				func() Contenter {
					return Executable{Shell: tc.blockType, Cmd: tc.content}
				},
				tc.errorMessage,
				func(m MaterializedContent, t *testing.T) {
					if m.Type != ExecutableType {
						t.Errorf("Expected Type to be %d, got %d", ExecutableType, m.Type)
					}

					if m.Content != strings.Join(tc.content, " ") {
						t.Errorf("Expected content to be %s, got %s", strings.Join(tc.content, " "), m.Content)
					}

					if val, ok := m.Metadata["Shell"]; ok {
						if val != tc.blockType {
							t.Errorf("Expected Shell to be %s, got %s", tc.blockType, val)
						}
					} else {
						t.Errorf("Did not find Shell")
					}

					if val, ok := m.Metadata["Command"]; ok {
						if !reflect.DeepEqual(val, tc.content) {
							t.Errorf("Expected Command to be %s got %s", tc.content, val)
						}
					} else {
						t.Errorf("Did not find Command")
					}
				},
			)
		})
	}
}

func TestRemoteMaterialize(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		errorMessage string
	}{
		{
			name:         "Passing",
			content:      "doot",
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMaterialize(
				t,
				func() Contenter {
					return Remote{Reader: strings.NewReader(tc.content)}
				},
				tc.errorMessage,
				func(m MaterializedContent, t *testing.T) {
					if m.Type != RemoteType {
						t.Errorf("Expected Type to be %d, got %d", RemoteType, m.Type)
					}

					if m.Content != tc.content {
						t.Errorf("Expected Content to be %s, got %s", tc.content, m.Content)
					}
				},
			)
		})
	}
}

func TestTableRowMaterialize(t *testing.T) {
	tests := []struct {
		name         string
		content      []string
		errorMessage string
	}{
		{
			name:         "Passing",
			content:      []string{"if", "i", "know", "only", "one", "thing"},
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMaterialize(
				t,
				func() Contenter {
					return TableRow{Values: tc.content}
				},
				tc.errorMessage,
				func(m MaterializedContent, t *testing.T) {
					if m.Type != TableRowType {
						t.Errorf("Expected Type to be %d, got %d", TableRowType, m.Type)
					}

					if m.Content != "" {
						t.Errorf("Expected Content to be an empty string")
					}

					if val, ok := m.Metadata["Items"]; ok {
						if !reflect.DeepEqual(val, tc.content) {
							t.Errorf("Expected Items to be %s, got %s", tc.content, val)
						}
					} else {
						t.Errorf("Did not find Items")
					}
				},
			)
		})
	}
}

func TestCommmentMaterialize(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		errorMessage string
	}{
		{
			name: "Passing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMaterialize(
				t,
				func() Contenter {
					return Comment(tc.content)
				},
				tc.errorMessage,
				func(m MaterializedContent, t *testing.T) {
					if m.Type != CommentType {
						t.Errorf("Expected Type to be %d, got %d", CommentType, m.Type)
					}

					if m.Content != tc.content {
						t.Errorf("Expected content to be %s, got %s", tc.content, m.Content)
					}
				},
			)
		})
	}
}
