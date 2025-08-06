package doyoucompute

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

// TaskStatus represents the outcome of executing a command or task.
type TaskStatus int

const (
	// COMPLETED indicates the task executed successfully without errors
	COMPLETED TaskStatus = iota + 1
	// FAILED indicates the task execution failed or encountered an error
	FAILED
)

// TaskResult contains the outcome and details of executing a single command,
// including context information and any errors that occurred.
type TaskResult struct {
	// SectionName identifies which document section the task originated from
	SectionName string
	// Command contains the full command string that was executed
	Command string
	// Status indicates whether the task completed successfully or failed
	Status TaskStatus
	// Error holds any error that occurred during task execution (nil if successful)
	Error error
}

// Runner defines the interface for executing command plans and returning results.
// This abstraction allows for different execution strategies (local, remote, mock, etc.).
type Runner interface {
	// Run executes the provided command plan and returns the result with status and error information
	Run(plan CommandPlan) TaskResult
}

// TaskRunner implements the Runner interface for executing commands locally
// using the operating system's command execution facilities.
type TaskRunner struct{}

// NewTaskRunner creates a new TaskRunner instance for local command execution.
func NewTaskRunner() TaskRunner {
	return TaskRunner{}
}

// Run executes a command plan locally using exec.Command, streaming output to
// stdout/stderr in real-time. Returns a TaskResult with execution status and any errors.
func (t TaskRunner) Run(plan CommandPlan) TaskResult {
	result := TaskResult{
		SectionName: plan.Context.Name,
		Command:     strings.Join(plan.Args, " "),
	}

	if len(plan.Args) == 0 {
		result.Error = errors.New("no command specified")
		result.Status = FAILED
		return result
	}

	cmd := exec.Command(plan.Args[0], plan.Args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		result.Error = err
		result.Status = FAILED
	} else {
		result.Status = COMPLETED
	}

	return result
}

// RunExecutionPlan executes a sequence of command plans using the provided runner,
// logging each command before execution and returning results for all commands.
// Commands are executed sequentially in the order they appear in the plan.
func RunExecutionPlan(plans []CommandPlan, runner Runner) []TaskResult {
	results := make([]TaskResult, len(plans))

	for idx, commandPlan := range plans {
		log.Printf("[Section: %s] - Running command: '%s'", commandPlan.Context.Name, strings.Join(commandPlan.Args, " "))

		results[idx] = runner.Run(commandPlan)
	}

	return results
}
