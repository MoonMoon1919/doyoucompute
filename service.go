package doyoucompute

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type Repository interface {
	Load(path string) (string, error)
	Save(path string, content string) error
}

type Service struct {
	repository        Repository
	taskRunner        TaskRunner
	fileRenderer      Renderer[string]
	executionRenderer Renderer[[]CommandPlan]
}

const ALL_SECTIONS = ""

func NewService(repo Repository, runner TaskRunner, fileRenderer Renderer[string], executionRenderer Renderer[[]CommandPlan]) *Service {
	return &Service{
		repository:        repo,
		taskRunner:        runner,
		fileRenderer:      fileRenderer,
		executionRenderer: executionRenderer,
	}
}

func (s Service) RenderFile(document *Document, outpath string) error {
	content, err := s.fileRenderer.Render(document)
	if err != nil {
		return err
	}

	return s.repository.Save(outpath, content)
}

type ComparisonResult struct {
	Matches      bool
	DocumentHash string
	FileHash     string
}

func (s Service) CompareFile(document *Document, pathToFile string) (ComparisonResult, error) {
	content, err := s.fileRenderer.Render(document)
	if err != nil {
		return ComparisonResult{}, err
	}

	loadedContent, err := s.repository.Load(pathToFile)
	if err != nil {
		return ComparisonResult{}, err
	}

	expectedHash := md5.Sum([]byte(content))
	currentHash := md5.Sum([]byte(loadedContent))

	return ComparisonResult{
		Matches:      expectedHash == currentHash,
		DocumentHash: hex.EncodeToString(expectedHash[:]),
		FileHash:     hex.EncodeToString(currentHash[:]),
	}, nil
}

func (s Service) PlanScriptExecution(document *Document, sectionName string) ([]CommandPlan, error) {
	executionPlan, err := s.executionRenderer.Render(document)
	if err != nil {
		return []CommandPlan{}, err
	}

	if sectionName == "" {
		return executionPlan, nil
	}

	var commands []CommandPlan

	for _, commandPlan := range executionPlan {
		if commandPlan.Context.Name == sectionName {
			commands = append(commands, commandPlan)
		}
	}

	if len(commands) == 0 {
		return []CommandPlan{}, fmt.Errorf("no executable blocks found for sectio '%s'", sectionName)
	}

	return commands, nil
}

func (s Service) ExecuteScript(document *Document, sectionName string) ([]TaskResult, error) {
	executionPlan, err := s.PlanScriptExecution(document, sectionName)
	if err != nil {
		return []TaskResult{}, err
	}

	results := RunExecutionPlan(executionPlan, s.taskRunner)

	return results, nil
}
