package doyoucompute

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

type TaskStatus int

const (
	COMPLETED TaskStatus = iota + 1
	FAILED
)

type TaskResult struct {
	SectionName string
	Command     string
	Status      TaskStatus
	Error       error
}

type Runner interface {
	Run(plan CommandPlan) TaskResult
}

type TaskRunner struct{}

func NewTaskRunner() TaskRunner {
	return TaskRunner{}
}

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

func RunExecutionPlan(plans []CommandPlan, runner Runner) []TaskResult {
	results := make([]TaskResult, len(plans))

	for idx, commandPlan := range plans {
		log.Printf("[Section: %s] - Running command: '%s'", commandPlan.Context.Name, strings.Join(commandPlan.Args, " "))

		results[idx] = runner.Run(commandPlan)
	}

	return results
}
