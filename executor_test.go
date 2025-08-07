package doyoucompute

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// Generic test harness for task runner operations
func testTaskRunnerOperation[T any](
	t *testing.T,
	operation func() (T, error),
	errorMessage string,
	validate func(T, *testing.T),
) {
	result, err := operation()

	checkErrors(errorMessage, err, t)
	if errorMessage != "" {
		return
	}

	validate(result, t)
}

// Mock runner for testing
type MockRunner struct {
	results []TaskResult
	calls   []CommandPlan
}

func (m *MockRunner) Run(plan CommandPlan) TaskResult {
	m.calls = append(m.calls, plan)
	if len(m.results) > 0 {
		result := m.results[0]
		m.results = m.results[1:]
		return result
	}
	return TaskResult{Status: COMPLETED}
}

func TestTaskRunner_Run(t *testing.T) {
	tests := []struct {
		name           string
		plan           CommandPlan
		expectedStatus TaskStatus
		expectedError  bool
		errorMessage   string
	}{
		{
			name: "Successful command execution",
			plan: CommandPlan{
				Shell: "sh",
				Args:  []string{"echo", "hello"},
				Context: SectionInfo{
					Name: "TestSection",
				},
			},
			expectedStatus: COMPLETED,
			expectedError:  false,
			errorMessage:   "",
		},
		{
			name: "Failed command execution",
			plan: CommandPlan{
				Shell: "sh",
				Args:  []string{"nonexistentcommand", "arg1"},
				Context: SectionInfo{
					Name: "TestSection",
				},
			},
			expectedStatus: FAILED,
			expectedError:  true,
			errorMessage:   "",
		},
		{
			name: "Command with multiple args",
			plan: CommandPlan{
				Shell: "sh",
				Args:  []string{"echo", "hello", "world"},
				Context: SectionInfo{
					Name: "MultiArgSection",
				},
			},
			expectedStatus: COMPLETED,
			expectedError:  false,
			errorMessage:   "",
		},
		{
			name: "Empty command args",
			plan: CommandPlan{
				Shell: "sh",
				Args:  []string{},
				Context: SectionInfo{
					Name: "EmptySection",
				},
			},
			expectedStatus: FAILED,
			expectedError:  true,
			errorMessage:   "",
		},
		{
			name: "Command that returns non-zero exit code",
			plan: CommandPlan{
				Shell: "sh",
				Args:  []string{"sh", "-c", "exit 1"},
				Context: SectionInfo{
					Name: "FailSection",
				},
			},
			expectedStatus: FAILED,
			expectedError:  true,
			errorMessage:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := NewTaskRunner(DefaultSecureConfig())

			testTaskRunnerOperation(
				t,
				func() (TaskResult, error) {
					return runner.Run(tc.plan), nil
				},
				tc.errorMessage,
				func(result TaskResult, t *testing.T) {
					// Validate status
					if result.Status != tc.expectedStatus {
						t.Errorf("Expected status %v, got %v", tc.expectedStatus, result.Status)
					}

					// Validate error presence
					hasError := result.Error != nil
					if hasError != tc.expectedError {
						t.Errorf("Expected error presence %v, got %v", tc.expectedError, hasError)
					}

					// Validate section name
					if result.SectionName != tc.plan.Context.Name {
						t.Errorf("Expected section name %s, got %s", tc.plan.Context.Name, result.SectionName)
					}

					// Validate command string
					expectedCommand := strings.Join(tc.plan.Args, " ")
					if result.Command != expectedCommand {
						t.Errorf("Expected command %s, got %s", expectedCommand, result.Command)
					}
				},
			)
		})
	}
}

func TestRunExecutionPlan(t *testing.T) {
	tests := []struct {
		name         string
		plans        []CommandPlan
		mockResults  []TaskResult
		errorMessage string
	}{
		{
			name:         "Empty execution plan",
			plans:        []CommandPlan{},
			mockResults:  []TaskResult{},
			errorMessage: "",
		},
		{
			name: "Single command plan",
			plans: []CommandPlan{
				{
					Args:    []string{"echo", "hello"},
					Context: SectionInfo{Name: "Section1"},
				},
			},
			mockResults: []TaskResult{
				{
					SectionName: "Section1",
					Command:     "echo hello",
					Status:      COMPLETED,
					Error:       nil,
				},
			},
			errorMessage: "",
		},
		{
			name: "Multiple command plans",
			plans: []CommandPlan{
				{
					Args:    []string{"echo", "hello"},
					Context: SectionInfo{Name: "Section1"},
				},
				{
					Args:    []string{"echo", "world"},
					Context: SectionInfo{Name: "Section2"},
				},
			},
			mockResults: []TaskResult{
				{
					SectionName: "Section1",
					Command:     "echo hello",
					Status:      COMPLETED,
					Error:       nil,
				},
				{
					SectionName: "Section2",
					Command:     "echo world",
					Status:      COMPLETED,
					Error:       nil,
				},
			},
			errorMessage: "",
		},
		{
			name: "Mixed success and failure results",
			plans: []CommandPlan{
				{
					Args:    []string{"echo", "hello"},
					Context: SectionInfo{Name: "Section1"},
				},
				{
					Args:    []string{"false"},
					Context: SectionInfo{Name: "Section2"},
				},
			},
			mockResults: []TaskResult{
				{
					SectionName: "Section1",
					Command:     "echo hello",
					Status:      COMPLETED,
					Error:       nil,
				},
				{
					SectionName: "Section2",
					Command:     "false",
					Status:      FAILED,
					Error:       errors.New("exit status 1"),
				},
			},
			errorMessage: "",
		},
		{
			name: "Plans with same section name",
			plans: []CommandPlan{
				{
					Args:    []string{"echo", "first"},
					Context: SectionInfo{Name: "SameSection"},
				},
				{
					Args:    []string{"echo", "second"},
					Context: SectionInfo{Name: "SameSection"},
				},
			},
			mockResults: []TaskResult{
				{
					SectionName: "SameSection",
					Command:     "echo first",
					Status:      COMPLETED,
					Error:       nil,
				},
				{
					SectionName: "SameSection",
					Command:     "echo second",
					Status:      COMPLETED,
					Error:       nil,
				},
			},
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRunner := &MockRunner{
				results: make([]TaskResult, len(tc.mockResults)),
			}
			copy(mockRunner.results, tc.mockResults)

			testTaskRunnerOperation(
				t,
				func() ([]TaskResult, error) {
					return RunExecutionPlan(tc.plans, mockRunner), nil
				},
				tc.errorMessage,
				func(results []TaskResult, t *testing.T) {
					// Validate result count
					if len(results) != len(tc.plans) {
						t.Errorf("Expected %d results, got %d", len(tc.plans), len(results))
						return
					}

					// Validate results match expected
					for i, result := range results {
						expected := tc.mockResults[i]
						if result.SectionName != expected.SectionName {
							t.Errorf("Result %d: Expected section name %s, got %s", i, expected.SectionName, result.SectionName)
						}
						if result.Command != expected.Command {
							t.Errorf("Result %d: Expected command %s, got %s", i, expected.Command, result.Command)
						}
						if result.Status != expected.Status {
							t.Errorf("Result %d: Expected status %v, got %v", i, expected.Status, result.Status)
						}
						if (result.Error == nil) != (expected.Error == nil) {
							t.Errorf("Result %d: Expected error presence %v, got %v", i, expected.Error != nil, result.Error != nil)
						}
						if result.Error != nil && expected.Error != nil && result.Error.Error() != expected.Error.Error() {
							t.Errorf("Result %d: Expected error %s, got %s", i, expected.Error.Error(), result.Error.Error())
						}
					}

					// Validate all plans were called
					if len(mockRunner.calls) != len(tc.plans) {
						t.Errorf("Expected %d runner calls, got %d", len(tc.plans), len(mockRunner.calls))
						return
					}

					// Validate plans were called in order
					for i, expectedPlan := range tc.plans {
						if !reflect.DeepEqual(mockRunner.calls[i], expectedPlan) {
							t.Errorf("Expected plan %d to be %v, got %v", i, expectedPlan, mockRunner.calls[i])
						}
					}
				},
			)
		})
	}
}

func TestNewTaskRunner(t *testing.T) {
	tests := []struct {
		name         string
		errorMessage string
	}{
		{
			name:         "Creates new TaskRunner instance",
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testTaskRunnerOperation(
				t,
				func() (TaskRunner, error) {
					return NewTaskRunner(DefaultSecureConfig()), nil
				},
				tc.errorMessage,
				func(result TaskRunner, t *testing.T) {
					// Validate it's a valid TaskRunner (interface check)
					var _ Runner = result
				},
			)
		})
	}
}
