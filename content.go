package doyoucompute

import (
	"io"
)

type ContentType int

const (
	HeaderType ContentType = iota + 1
	LinkType
	TextType
	CodeType
	CodeBlockType
	BlockQuoteType
	ExecutableType
	RemoteType
	ListType
	ParagraphType
	SectionType
	DocumentType
)

type MaterializedContent struct {
	Type     ContentType
	Metadata map[string]interface{}
}

type Header struct {
	Content string
	Level   int
}

func (h Header) Type() ContentType {
	return HeaderType
}

func (h Header) Materialize() (MaterializedContent, error) {
	// headerLevel := strings.Repeat("#", h.Level)
	// return headerLevel + " " + h.Content, nil

	return MaterializedContent{}, nil
}

type Text string

func (t Text) Type() ContentType {
	return TextType
}

func (t Text) Materialize() (MaterializedContent, error) {
	return MaterializedContent{}, nil
}

type Link string

func (t Link) Type() ContentType {
	return LinkType
}

func (l Link) Materialize() (MaterializedContent, error) {
	return MaterializedContent{}, nil
}

type Code string

func (c Code) Type() ContentType {
	return CodeType
}

func (c Code) Materialize() (MaterializedContent, error) {
	return MaterializedContent{}, nil
}

// A codeblock is a NON-EXECUTABLE block of code
// Useful for examples/payloads etc
type CodeBlock struct {
	Shell string
	Cmd   []string
}

func (c CodeBlock) Type() ContentType {
	return CodeBlockType
}

func (c CodeBlock) Materialize() (MaterializedContent, error) {
	// cmd := strings.Join(c.Cmd, " ")
	// leadingText := strings.Join([]string{"```", c.Shell}, "")
	// return strings.Join([]string{leadingText, cmd, "```"}, "\n"), nil
	return MaterializedContent{}, nil
}

type BlockQuote string

func (b BlockQuote) Type() ContentType {
	return BlockQuoteType
}

func (b BlockQuote) Materialize() (MaterializedContent, error) {
	return MaterializedContent{}, nil
}

// Script running
// An executable code block
type Executable struct {
	Shell string
	Cmd   []string
}

func (e Executable) Type() ContentType {
	return ExecutableType
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

func (r Remote) Type() ContentType {
	return RemoteType
}

func (r Remote) Materialize() (MaterializedContent, error) {
	// content, err := io.ReadAll(r.Reader)
	// if err != nil {
	// 	return "", err
	// }

	// return string(content), nil
	return MaterializedContent{}, nil
}
