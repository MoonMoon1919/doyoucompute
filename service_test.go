package doyoucompute

import (
	"errors"
	"testing"
)

type FakeFileRepo struct {
	files map[string]string
}

func NewFakeFileRepo() *FakeFileRepo {
	return &FakeFileRepo{
		files: map[string]string{},
	}
}

func (f *FakeFileRepo) Load(path string) (string, error) {
	file, ok := f.files[path]

	if !ok {
		return "", errors.New("file not found")
	}

	return file, nil
}

func (f *FakeFileRepo) Save(path string, content string) error {
	f.files[path] = content

	return nil
}

func FakeRunner(plan CommandPlan) TaskResult {
	return TaskResult{}
}

func newService() Service {
	return NewService(
		NewFakeFileRepo(),
		FakeRunner,
		NewMarkdownRenderer(),
		NewExecutionRenderer(),
	)
}

func newDocument() Document {
	return Document{
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
					Executable{
						Shell: "bash",
						Cmd:   []string{"echo", "hello", "world"},
					},
					Section{
						Name: "Quick Start",
						Content: []Node{
							Text("Install dependencies"),
							Executable{
								Shell: "bash",
								Cmd:   []string{"go", "get"},
							},
						},
					},
				},
			},
		},
	}
}

func TestRenderFile(t *testing.T) {
	tests := []struct {
		name         string
		document     Document
		svc          Service
		outpath      string
		errorMessage string
	}{
		{
			name:         "Passing",
			document:     newDocument(),
			svc:          newService(),
			outpath:      "test.md",
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// When
			err := tc.svc.RenderFile(&tc.document, tc.outpath)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			// Then
			if errMsg != tc.errorMessage {
				t.Errorf("Expected error %s, got %s", tc.errorMessage, errMsg)
			}

			comparisonResult, err := tc.svc.CompareFile(&tc.document, tc.outpath)
			if err != nil {
				t.Errorf("Error during comparison %s", err.Error())
			}

			if !comparisonResult.Matches {
				t.Errorf("expected comparison match, Document Hash %s, File Hash %s", comparisonResult.DocumentHash, comparisonResult.FileHash)
			}
		})
	}
}

func TestCompareFile(t *testing.T) {
	tests := []struct {
		name         string
		document     Document
		svc          Service
		outpath      string
		errorMessage string
	}{
		{
			name:         "Passing",
			document:     newDocument(),
			svc:          newService(),
			outpath:      "test.md",
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			err := tc.svc.RenderFile(&tc.document, tc.outpath)
			if err != nil {
				t.Errorf("unexpected error %s rendering file", err.Error())
			}

			// When
			comparisonResult, err := tc.svc.CompareFile(&tc.document, tc.outpath)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			// Then
			if errMsg != tc.errorMessage {
				t.Errorf("Expected error %s, got %s", tc.errorMessage, errMsg)
			}

			if !comparisonResult.Matches {
				t.Errorf("expected comparison match, Document Hash %s, File Hash %s", comparisonResult.DocumentHash, comparisonResult.FileHash)
			}
		})
	}
}

func TestPlanScriptExecution(t *testing.T) {
	tests := []struct {
		name string
		svc  Service
	}{
		{
			name: "Passing",
			svc:  newService(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given

			// When

			// Then
		})
	}
}

func TestExecuteScript(t *testing.T) {
	tests := []struct {
		name string
		svc  Service
	}{
		{
			name: "Passing",
			svc:  newService(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given

			// When

			// Then
		})
	}
}
