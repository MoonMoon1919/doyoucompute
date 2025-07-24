package doyoucompute

import (
	"errors"
	"strings"
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

type MockTaskRunner struct {
	expectations map[string]TaskResult
}

func (m MockTaskRunner) Run(plan CommandPlan) TaskResult {
	key := strings.Join(plan.Args, " ")
	if result, exists := m.expectations[key]; exists {
		return result
	}

	return TaskResult{
		SectionName: plan.Context.Name,
		Command:     key,
		Status:      COMPLETED,
		Error:       nil,
	}
}

func newService(taskRunnerExpectations map[string]TaskResult) Service {
	return NewService(
		NewFakeFileRepo(),
		MockTaskRunner{expectations: taskRunnerExpectations},
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
			svc:          newService(map[string]TaskResult{}),
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
		matches      bool
	}{
		{
			name:         "Passing",
			document:     newDocument(),
			svc:          newService(map[string]TaskResult{}),
			outpath:      "test.md",
			errorMessage: "",
			matches:      true,
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

			if comparisonResult.Matches != tc.matches {
				t.Errorf("expected comparison to be %v, Document Hash %s, File Hash %s", tc.matches, comparisonResult.DocumentHash, comparisonResult.FileHash)
			}
		})
	}
}

func TestPlanScriptExecution(t *testing.T) {
	tests := []struct {
		name         string
		document     Document
		svc          Service
		errorMessage string
		expected     []CommandPlan
	}{
		{
			name:         "Passing",
			document:     newDocument(),
			svc:          newService(map[string]TaskResult{}),
			errorMessage: "",
			expected: []CommandPlan{
				{
					Shell: "bash",
					Args:  []string{"echo", "hello", "world"},
					Context: SectionInfo{
						Name:  "INTRO",
						Level: 2,
					},
				},
				{
					Shell: "bash",
					Args:  []string{"go", "get"},
					Context: SectionInfo{
						Name:  "Quick Start",
						Level: 3,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// When
			commands, err := tc.svc.PlanScriptExecution(&tc.document, ALL_SECTIONS)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			// Then
			if errMsg != tc.errorMessage {
				t.Errorf("Expected error %s, got %s", tc.errorMessage, errMsg)
			}

			for idx, expected := range tc.expected {
				found := commands[idx]

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

func TestExecuteScript(t *testing.T) {
	tests := []struct {
		name              string
		document          Document
		taskRunnerResults map[string]TaskResult
		errorMessage      string
	}{
		{
			name:     "Passing",
			document: newDocument(),
			taskRunnerResults: map[string]TaskResult{
				"echo hello world": {
					SectionName: "INTRO",
					Command:     "echo hello world",
					Status:      COMPLETED,
					Error:       nil,
				},
				"go get": {
					SectionName: "Quick Start",
					Command:     "go get",
					Status:      COMPLETED,
					Error:       nil,
				},
			},
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			svc := newService(tc.taskRunnerResults)

			var expected []TaskResult
			for _, result := range tc.taskRunnerResults {
				expected = append(expected, result)
			}

			expectedNumTasks := len(expected)

			// When
			results, err := svc.ExecuteScript(&tc.document, ALL_SECTIONS)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			// Then
			if errMsg != tc.errorMessage {
				t.Errorf("Expected error %s, got %s", tc.errorMessage, errMsg)
			}

			if len(results) != expectedNumTasks {
				t.Errorf("Expected %d results, got %d", expectedNumTasks, len(results))
			}

			for idx, result := range results {
				expected := expected[idx]

				if expected != result {
					t.Errorf("Expected TaskResult %v, got %v", expected, result)
				}
			}
		})
	}
}
