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
			expected: "# MyDoc\n\n## INTRO\n\nThis is an introduction. And another sentence here.\n\nhey im some remote content\n",
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

func TestExecutionPlanRender(t *testing.T) {
	tests := []struct {
		name         string
		renderer     ExecutionPlan
		document     Document
		errorMessage string
		expected     []CommandPlan
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
							Executable{
								Shell: "bash",
								Cmd:   []string{"echo", "hello", "world"},
							},
						},
					},
				},
			},
			expected: []CommandPlan{
				{
					Shell: "bash",
					Args:  []string{"echo", "hello", "world"},
					Context: SectionInfo{
						Name:  "INTRO",
						Level: 2,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.renderer.Render(tc.document)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if errMsg != tc.errorMessage {
				t.Errorf("Expected error %s, got %s", tc.errorMessage, errMsg)
			}

			for idx, expected := range tc.expected {
				found := tc.renderer.Commands[idx]

				if found.Context != expected.Context {
					t.Errorf("Expected context %v, got %v", expected.Context, found.Context)
				}

				if found.Shell != expected.Shell {
					t.Errorf("Expected context %s, got %s", expected.Shell, found.Shell)
				}

				for idx, expectedArg := range expected.Args {
					foundArg := found.Args[idx]

					if expectedArg != foundArg {
						t.Errorf("expected arg %s at index %d, found %s", expectedArg, idx, foundArg)
					}
				}
			}
		})
	}
}
