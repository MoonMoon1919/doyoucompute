package content

import (
	"io"
	"strings"
)

type Contenter interface {
	Render() (string, error)
}

// Content mapped to markdown
type Header struct {
	Content string
	Level   int
}

func (h Header) Render() (string, error) { // Markdown specific..
	headerLevel := strings.Repeat("#", h.Level)

	return headerLevel + " " + h.Content, nil
}

type Paragraph string

func (f Paragraph) Render() (string, error) {
	return string(f), nil
}

type OrderedList struct {
	items []string
}

func (l OrderedList) Render() (string, error) {
	return "", nil
}

type UnorderedList struct {
	items []string
}

func (l UnorderedList) Render() (string, error) {
	return "", nil
}

type Link string

func (l Link) Render() (string, error) {
	return "", nil
}

type Code string

func (c Code) Render() (string, error) {
	return "", nil
}

// A codeblock is a NON-EXECUTABLE block of code
// Useful for examples/payloads etc
type CodeBlock struct {
	Shell string
	Cmd   []string
}

func (c CodeBlock) Render() (string, error) {
	cmd := strings.Join(c.Cmd, " ")

	leadingText := strings.Join([]string{"```", c.Shell}, "")

	return strings.Join([]string{leadingText, cmd, "```"}, "\n"), nil
}

type BlockQuote string

func (b BlockQuote) Render() (string, error) {
	return "", nil
}

// Script running
// An executable code block
type Executable struct {
	Shell string
	Cmd   []string
}

func (c Executable) Render() (string, error) {
	cmd := strings.Join(c.Cmd, " ")

	leadingText := strings.Join([]string{"```", c.Shell}, "")

	return strings.Join([]string{leadingText, cmd, "```"}, "\n"), nil
}

// Content Sources
type Remote struct { // e.g., from local file in docs folder, from GitHub.. etc
	Reader io.Reader
}

func (r Remote) Render() (string, error) {
	content, err := io.ReadAll(r.Reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
