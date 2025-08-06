package doyoucompute

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// Repository provides abstraction for file system operations, allowing the service
// to load and save content without being tied to specific storage implementations.
type Repository interface {
	// Load reads content from the specified file path and returns it as a string.
	// Returns an error if the file cannot be read or does not exist.
	Load(path string) (string, error)
	// Save writes the provided content to the specified file path.
	// Returns an error if the file cannot be written.
	Save(path string, content string) error
}

// Service orchestrates the core functionality of the runnable documentation system,
// coordinating between content rendering, script execution, and file operations.
type Service struct {
	repository        Repository
	taskRunner        Runner
	fileRenderer      Renderer[string]
	executionRenderer Renderer[[]CommandPlan]
}

// ALL_SECTIONS is a constant used to indicate that all sections should be processed
// when no specific section name is provided.
const ALL_SECTIONS = ""

// NewService creates a new Service instance with the provided dependencies for
// repository access, task execution, and content rendering.
func NewService(repo Repository, runner Runner, fileRenderer Renderer[string], executionRenderer Renderer[[]CommandPlan]) Service {
	return Service{
		repository:        repo,
		taskRunner:        runner,
		fileRenderer:      fileRenderer,
		executionRenderer: executionRenderer,
	}
}

// RenderFile generates the final content for a document and saves it to the specified output path.
// Returns an error if rendering fails or the file cannot be saved.
func (s Service) RenderFile(document *Document, outpath string) error {
	content, err := s.fileRenderer.Render(document)
	if err != nil {
		return err
	}

	return s.repository.Save(outpath, content)
}

// ComparisonResult contains the results of comparing a document's rendered content
// with an existing file, including match status and content hashes.
type ComparisonResult struct {
	// Matches indicates whether the document content matches the existing file
	Matches bool
	// DocumentHash is the MD5 hash of the rendered document content
	DocumentHash string
	// FileHash is the MD5 hash of the existing file content
	FileHash string
}

// CompareFile renders a document and compares its content with an existing file,
// returning detailed comparison results including MD5 hashes for verification.
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

// PlanScriptExecution analyzes a document and creates an execution plan for all executable
// content blocks. If sectionName is provided, only executable blocks from that section
// are included. Use ALL_SECTIONS constant to include all sections.
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
		return []CommandPlan{}, fmt.Errorf("no executable blocks found for section '%s'", sectionName)
	}

	return commands, nil
}

// ExecuteScript creates an execution plan for the specified document section and runs
// all executable blocks, returning the results of each executed command.
func (s Service) ExecuteScript(document *Document, sectionName string) ([]TaskResult, error) {
	executionPlan, err := s.PlanScriptExecution(document, sectionName)
	if err != nil {
		return []TaskResult{}, err
	}

	results := RunExecutionPlan(executionPlan, s.taskRunner)

	return results, nil
}
