package doyoucompute

import "io"

// LoadFile reads all content from the provided io.Reader and returns it as a string.
// This utility function abstracts file reading operations and can work with any
// io.Reader implementation (files, network streams, etc.).
func LoadFile(reader io.Reader) (string, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// WriteFile writes the provided content string to the given io.Writer.
// This utility function abstracts file writing operations and can work with any
// io.Writer implementation (files, network streams, buffers, etc.).
func WriteFile(writer io.Writer, content string) error {
	_, err := writer.Write([]byte(content))
	return err
}
