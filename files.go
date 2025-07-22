package doyoucompute

import "io"

func WriteFile(writer io.Writer, document Document, renderer Renderer[string]) error {
	content, err := renderer.Render(document)
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(content))
	return err
}
