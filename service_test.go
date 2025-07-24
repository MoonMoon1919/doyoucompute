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

func checkErrors(expectedErrorMsg string, err error, t *testing.T) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	if errMsg != expectedErrorMsg {
		t.Errorf("expected error %s, got %s", expectedErrorMsg, errMsg)
	}
}

func testServiceOperation[T any](
	t *testing.T,
	operation func(*Service) (T, error),
	taskRunnerResults map[string]TaskResult,
	errorMessage string,
	comparisonFunc func(T, *Service, *testing.T),
) {
	svc := newService(taskRunnerResults)

	res, err := operation(&svc)

	checkErrors(errorMessage, err, t)
	if errorMessage != "" {
		return // bail out before validation for tests w/ non-null errors
	}

	comparisonFunc(res, &svc, t)
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
			outpath:      "test.md",
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testServiceOperation(
				t,
				func(s *Service) (string, error) {
					err := s.RenderFile(&tc.document, tc.outpath)

					return "", err
				},
				map[string]TaskResult{},
				tc.errorMessage,
				func(res string, svc *Service, t *testing.T) {
					comparisonResult, err := svc.CompareFile(&tc.document, tc.outpath)
					if err != nil {
						t.Errorf("Error during comparison %s", err.Error())
					}

					if !comparisonResult.Matches {
						t.Errorf("expected comparison match, Document Hash %s, File Hash %s", comparisonResult.DocumentHash, comparisonResult.FileHash)
					}
				},
			)
		})
	}
}

func TestCompareFile(t *testing.T) {
	tests := []struct {
		name         string
		document     Document
		outpath      string
		errorMessage string
		matches      bool
	}{
		{
			name:         "Passing",
			document:     newDocument(),
			outpath:      "test.md",
			errorMessage: "",
			matches:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testServiceOperation(
				t,
				func(s *Service) (ComparisonResult, error) {
					err := s.RenderFile(&tc.document, tc.outpath)
					if err != nil {
						t.Errorf("unexpected error %s rendering file", err.Error())
					}

					return s.CompareFile(&tc.document, tc.outpath)
				},
				map[string]TaskResult{},
				tc.errorMessage,
				func(cr ComparisonResult, s *Service, t *testing.T) {
					if cr.Matches != tc.matches {
						t.Errorf("expected comparison to be %v, Document Hash %s, File Hash %s", tc.matches, cr.DocumentHash, cr.FileHash)
					}
				},
			)
		})
	}
}

func TestPlanScriptExecution(t *testing.T) {
	tests := []struct {
		name         string
		document     Document
		errorMessage string
		expected     []CommandPlan
	}{
		{
			name:         "Passing",
			document:     newDocument(),
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
			testServiceOperation(
				t,
				func(s *Service) ([]CommandPlan, error) {
					return s.PlanScriptExecution(&tc.document, ALL_SECTIONS)
				},
				map[string]TaskResult{},
				tc.errorMessage,
				func(cp []CommandPlan, s *Service, t *testing.T) {
					for idx, expected := range tc.expected {
						found := cp[idx]

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
				},
			)

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
			testServiceOperation(
				t,
				func(s *Service) ([]TaskResult, error) {
					return s.ExecuteScript(&tc.document, ALL_SECTIONS)
				},
				tc.taskRunnerResults,
				tc.errorMessage,
				func(tr []TaskResult, s *Service, t *testing.T) {
					var expected []TaskResult
					for _, result := range tc.taskRunnerResults {
						expected = append(expected, result)
					}

					for idx, result := range tr {
						expected := expected[idx]

						if expected != result {
							t.Errorf("Expected TaskResult %v, got %v", expected, result)
						}
					}
				},
			)
		})
	}
}
