package doyoucompute

import "os"

// FileRepository implements the Repository interface using the local file system
// for loading and saving content. It provides concrete file operations for the
// runnable documentation system.
type FileRepository struct{}

// NewFileRepository creates a new FileRepository instance for file system operations.
func NewFileRepository() FileRepository {
	return FileRepository{}
}

// Load opens and reads the content of a file at the specified path, returning
// the content as a string. Returns an error if the file cannot be opened or read.
func (f FileRepository) Load(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	return LoadFile(file)
}

// Save creates a new file at the specified path and writes the provided content to it.
// If the file already exists, it will be overwritten. Returns an error if the file
// cannot be created or written to.
func (f FileRepository) Save(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	return WriteFile(file, content)
}
