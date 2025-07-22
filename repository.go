package doyoucompute

import "os"

type FileRepository struct{}

func (f FileRepository) Load(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	return LoadFile(file)
}

func (f FileRepository) Save(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	return WriteFile(file, content)
}
