package doyoucompute

import "io"

func LoadFile(reader io.Reader) (string, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func WriteFile(writer io.Writer, content string) error {
	_, err := writer.Write([]byte(content))
	return err
}
