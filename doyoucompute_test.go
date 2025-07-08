package doyoucompute

import (
	"strings"
	"testing"

	"github.com/MoonMoon1919/doyoucompute/pkg/content"
)

func TestSectionRender(t *testing.T) {
	tests := []struct {
		name         string
		section      Section
		errorMessage string
		expected     string
	}{
		{
			name: "Passing",
			section: Section{
				Name: "INTRO",
				Content: []content.Contenter{
					content.Paragraph("This is an introduction"),
					content.Remote{
						Reader: strings.NewReader("hey im some remote content"),
					},
				},
			},
			expected: "# INTRO\n\nThis is an introduction\n\nhey im some remote content",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			content, err := tc.section.Render()

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
