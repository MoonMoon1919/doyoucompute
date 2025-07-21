package doyoucompute

import (
	"log"
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
	Output      string
	Status      TaskStatus
	Error       error
}

type TaskRunner func(plan CommandPlan) TaskResult

func RunTask(plan CommandPlan) TaskResult {
	cmd := exec.Command(plan.Args[0], plan.Args[1:]...)

	output, err := cmd.Output()
	if err != nil {
		return TaskResult{
			SectionName: plan.Context.Name,
			Command:     strings.Join(plan.Args, " "),
			Output:      string(output),
			Status:      FAILED,
			Error:       err,
		}
	}

	return TaskResult{
		SectionName: plan.Context.Name,
		Command:     strings.Join(plan.Args, " "),
		Output:      string(output),
		Status:      COMPLETED,
		Error:       nil,
	}
}

func RunExecutionPlan(plans []CommandPlan, runner TaskRunner) []TaskResult {
	results := make([]TaskResult, len(plans))

	for idx, commandPlan := range plans {
		log.Printf("[Section: %s] - Running command: '%s'", commandPlan.Context.Name, strings.Join(commandPlan.Args, " "))

		results[idx] = runner(commandPlan)
	}

	return results
}
