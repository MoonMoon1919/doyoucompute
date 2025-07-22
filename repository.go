package doyoucompute

import "os"

type Repository interface {
	Save(path string, document Document, renderer Renderer[string]) error
}

type FileRepository struct{}

func (f FileRepository) Save(path string, document Document, renderer Renderer[string]) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	return WriteFile(file, document, renderer)
}
