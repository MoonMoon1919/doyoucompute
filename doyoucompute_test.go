package doyoucompute

import (
	"strings"
	"testing"
)

func TestMarkdownRender(t *testing.T) {
	tests := []struct {
		name         string
		renderer     Markdown
		document     Document
		errorMessage string
		expected     string
	}{
		{
			name: "Passing",
			document: Document{
				Name: "MyDoc",
				Content: []Node{
					Section{
						Name: "INTRO",
						Content: []Node{
							Paragraph{
								Items: []Node{
									Text("This is an introduction."),
									Text("And another sentence here."),
								},
							},
							Remote{
								Reader: strings.NewReader("hey im some remote content"),
							},
						},
					},
				},
			},
			expected: "# MyDoc\n\n## INTRO\n\nThis is an introduction. And another sentence here. \n\nhey im some remote content\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			content, err := tc.renderer.Render(tc.document)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if errMsg != tc.errorMessage {
				t.Errorf("Expected error %s, got %s", tc.errorMessage, errMsg)
			}

			if content != tc.expected {
				t.Errorf("Expected content %s, got %s", tc.expected, content)
			}
		})
	}
}
