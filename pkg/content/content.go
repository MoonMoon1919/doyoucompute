package content

import (
	"io"
)

// Content mapped to markdown
type Header struct {
	Content string
	Level   int
}

func (h Header) Materialize() (MaterializedContent, error) {
	// headerLevel := strings.Repeat("#", h.Level)
	// return headerLevel + " " + h.Content, nil

	return MaterializedContent{}, nil
}

type Paragraph string

func (f Paragraph) Materialize() (MaterializedContent, error) {
	// return string(f), nil
	return MaterializedContent{}, nil
}

type List struct {
	items   []string
	ordered bool
}

func (l List) Materialize() (MaterializedContent, error) {
	return MaterializedContent{}, nil
}

type Link string

func (l Link) Materialize() (MaterializedContent, error) {
	return MaterializedContent{}, nil
}

type Code string

func (c Code) Materialize() (MaterializedContent, error) {
	return MaterializedContent{}, nil
}

// A codeblock is a NON-EXECUTABLE block of code
// Useful for examples/payloads etc
type CodeBlock struct {
	Shell string
	Cmd   []string
}

func (c CodeBlock) Materialize() (MaterializedContent, error) {
	// cmd := strings.Join(c.Cmd, " ")
	// leadingText := strings.Join([]string{"```", c.Shell}, "")
	// return strings.Join([]string{leadingText, cmd, "```"}, "\n"), nil
	return MaterializedContent{}, nil
}

type BlockQuote string

func (b BlockQuote) Materialize() (MaterializedContent, error) {
	return MaterializedContent{}, nil
}

// Script running
// An executable code block
type Executable struct {
	Shell string
	Cmd   []string
}

func (c Executable) Materialize() (MaterializedContent, error) {
	// cmd := strings.Join(c.Cmd, " ")
	// leadingText := strings.Join([]string{"```", c.Shell}, "")
	// return strings.Join([]string{leadingText, cmd, "```"}, "\n"), nil

	return MaterializedContent{}, nil
}

// Content Sources
type Remote struct { // e.g., from local file in docs folder, from GitHub.. etc
	Reader io.Reader
}

func (r Remote) Materialize() (MaterializedContent, error) {
	// content, err := io.ReadAll(r.Reader)
	// if err != nil {
	// 	return "", err
	// }

	// return string(content), nil
	return MaterializedContent{}, nil
}
