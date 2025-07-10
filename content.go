package doyoucompute

import (
	"io"
	"strings"
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
	Content  string
	Metadata map[string]interface{}
}

type Header struct {
	Content string
	Level   int
}

func (h Header) Type() ContentType { return HeaderType }

func (h Header) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:    HeaderType,
		Content: h.Content,
		Metadata: map[string]interface{}{
			"Level": h.Level,
		},
	}, nil
}

type Text string

func (t Text) Type() ContentType { return TextType }

func (t Text) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:     TextType,
		Content:  string(t),
		Metadata: map[string]interface{}{},
	}, nil
}

type Link struct {
	Text string
	Url  string
}

func (l Link) Type() ContentType { return LinkType }

func (l Link) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:    LinkType,
		Content: l.Text,
		Metadata: map[string]interface{}{
			"Url": l.Url,
		},
	}, nil
}

type Code string

func (c Code) Type() ContentType {
	return CodeType
}

func (c Code) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:     CodeType,
		Content:  string(c),
		Metadata: map[string]interface{}{},
	}, nil
}

// A codeblock is a NON-EXECUTABLE block of code
// Useful for examples/payloads etc
type CodeBlock struct {
	BlockType string
	Cmd       []string
}

func (c CodeBlock) Type() ContentType { return CodeBlockType }

func (c CodeBlock) Materialize() (MaterializedContent, error) {
	// leadingText := strings.Join([]string{"```", c.Shell}, "")
	// return strings.Join([]string{leadingText, cmd, "```"}, "\n"), nil
	return MaterializedContent{
		Type:    CodeBlockType,
		Content: strings.Join(c.Cmd, " "),
		Metadata: map[string]interface{}{
			"BlockType": c.BlockType,
		},
	}, nil
}

type BlockQuote string

func (b BlockQuote) Type() ContentType { return BlockQuoteType }

func (b BlockQuote) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:     BlockQuoteType,
		Content:  string(b),
		Metadata: map[string]interface{}{},
	}, nil
}

// Script running
// An executable code block
type Executable struct {
	Shell string
	Cmd   []string
}

func (e Executable) Type() ContentType { return ExecutableType }

func (c Executable) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:    ExecutableType,
		Content: strings.Join(c.Cmd, " "),
		Metadata: map[string]interface{}{
			"Shell":      c.Shell,
			"Executable": true,
		},
	}, nil
}

// Content Sources
type Remote struct { // e.g., from local file in docs folder, from GitHub.. etc
	Reader io.Reader
}

func (r Remote) Type() ContentType { return RemoteType }

func (r Remote) Materialize() (MaterializedContent, error) {
	content, err := io.ReadAll(r.Reader)
	if err != nil {
		return MaterializedContent{}, err
	}

	return MaterializedContent{
		Type:     RemoteType,
		Content:  string(content),
		Metadata: map[string]interface{}{},
	}, nil
}
