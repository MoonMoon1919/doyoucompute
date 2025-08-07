package doyoucompute

import (
	"context"
	"errors"
	"fmt"
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
type TaskRunner struct {
	config ExecutionConfig
}

// NewTaskRunner creates a new TaskRunner instance for local command execution.
func NewTaskRunner(config ExecutionConfig) TaskRunner {
	return TaskRunner{
		config: config,
	}
}

func validateEnvironment(requiredEnvVars []string) error {
	var missing []string

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			missing = append(missing, envVar)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("required environment variables not set: %v", missing)
	}

	return nil
}

// Run executes a command plan locally using exec.Command, streaming output to
// stdout/stderr in real-time. Returns a TaskResult with execution status and any errors.
func (t TaskRunner) Run(plan CommandPlan) TaskResult {
	result := TaskResult{
		SectionName: plan.Context.Name,
		Command:     strings.Join(plan.Args, " "),
	}

	if err := ValidateCommandPlan(plan, t.config); err != nil {
		result.Error = fmt.Errorf("security validation failed: %w", err)
		result.Status = FAILED
		return result
	}

	// Check required environment variables
	if err := validateEnvironment(plan.Environment); err != nil {
		result.Error = fmt.Errorf("environment validation failed: %w", err)
		result.Status = FAILED
		return result
	}

	if len(plan.Args) == 0 {
		result.Error = errors.New("no command specified")
		result.Status = FAILED
		return result
	}

	ctx := context.Background()
	if t.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, t.config.Timeout)
		defer cancel()
	}

	var cmd *exec.Cmd

	if plan.Shell == "sh" || plan.Shell == "bash" {
		// We currently only support variable expansion for bash/sh
		cmd = exec.CommandContext(ctx, plan.Shell, "-c", strings.Join(plan.Args, " "))
	} else {
		cmd = exec.CommandContext(ctx, plan.Args[0], plan.Args[1:]...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("[Section: %s] - Running command: '%s'", plan.Context.Name, strings.Join(plan.Args, " "))
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
		results[idx] = runner.Run(commandPlan)
	}

	return results
}
