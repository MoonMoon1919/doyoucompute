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
	TableRowType
	BlockQuoteType
	ExecutableType
	RemoteType
	ListType
	TableType
	ParagraphType
	SectionType
	DocumentType
)

type CodeBlockExecType int // this name is awful

const (
	Static CodeBlockExecType = iota + 1
	Exec
)

type MaterializedContent struct {
	Type     ContentType
	Content  string
	Metadata map[string]interface{}
}

type Header struct {
	Content string
}

func (h Header) Type() ContentType { return HeaderType }

func (h Header) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:     h.Type(),
		Content:  h.Content,
		Metadata: map[string]interface{}{},
	}, nil
}

type Text string

func (t Text) Type() ContentType { return TextType }

func (t Text) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:     t.Type(),
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
		Type:    l.Type(),
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
		Type:     c.Type(),
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
		Type:    c.Type(),
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
		Type:     b.Type(),
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

func (e Executable) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:    e.Type(),
		Content: strings.Join(e.Cmd, " "),
		Metadata: map[string]interface{}{
			"Shell":      e.Shell,
			"Executable": true,
			"Command":    e.Cmd,
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
		Type:     r.Type(),
		Content:  string(content),
		Metadata: map[string]interface{}{},
	}, nil
}

// TableRow
type TableRow struct {
	Values []string
}

func (t TableRow) Type() ContentType { return TableRowType }

func (t TableRow) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:    t.Type(),
		Content: "",
		Metadata: map[string]interface{}{
			"Items": t.Values,
		},
	}, nil
}
