package doyoucompute

import (
	"io"
	"strings"
)

// ContentType represents the different types of content elements that can be
// processed and rendered in runnable documentation. Each type corresponds to
// a specific markdown or documentation construct that may require different
// handling during parsing and execution.
type ContentType int

const (
	// HeaderType represents markdown headers (# ## ### etc.)
	HeaderType ContentType = iota + 1

	// LinkType represents hyperlinks and references
	LinkType

	// TextType represents plain text content
	TextType

	// CodeType represents inline code snippets
	CodeType

	// CodeBlockType represents fenced code blocks that may contain
	// examples
	CodeBlockType

	// TableRowType represents individual rows within a table
	TableRowType

	// BlockQuoteType represents quoted text blocks (> quoted text)
	BlockQuoteType

	// ExecutableType represents code or commands that can be executed
	// as part of the runnable documentation
	ExecutableType

	// RemoteType represents content that is fetched from remote sources
	RemoteType

	// CommentType represents comment blocks in the documentation
	CommentType

	// ListType represents ordered and unordered lists
	ListType

	// TableType represents table structures
	TableType

	// ParagraphType represents standard paragraph content
	ParagraphType

	// SectionType represents logical sections or divisions of content
	SectionType

	// DocumentType represents the root document container
	DocumentType

	// FrontmatterType represents YAML/TOML frontmatter metadata
	// typically found at the beginning of markdown documents
	FrontmatterType
)

// CodeBlockExecType represents how a code block should be processed during
// documentation generation - either as static display content or as executable code
type CodeBlockExecType int // todo: this name is awful

const (
	// Static indicates the code block should be displayed as-is without execution
	Static CodeBlockExecType = iota + 1

	// Exec indicates the code block contains executable code that should be run
	// when processing the documentation as a script
	Exec
)

// MaterializedContent represents a processed content element with its type,
// rendered content, and associated metadata.
type MaterializedContent struct {
	// Type specifies what kind of content this represents (header, code block, etc.)
	Type ContentType

	// Content contains the final rendered or processed content as a string
	Content string

	// Metadata holds additional key-value pairs associated with this content,
	// such as styling information
	Metadata map[string]interface{}
}

// MARK: Header
// Header represents a header element containing text content
type Header struct {
	// Content holds the text content of the header
	Content string
}

// Type returns the ContentType for this header element
func (h Header) Type() ContentType { return HeaderType }

// Materialize converts the header into a MaterializedContent with its content
// and an empty metadata map. Headers are processed as-is without transformation
func (h Header) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:     h.Type(),
		Content:  h.Content,
		Metadata: map[string]interface{}{},
	}, nil
}

// MARK: Text
// Text represents plain text content in documentation. It is defined as a string type
// that implements content materialization behavior.
type Text string

// Type returns the ContentType for this text element.
func (t Text) Type() ContentType { return TextType }

// Materialize converts the text into a MaterializedContent with its string content
// and an empty metadata map. Plain text is processed as-is without transformation.
func (t Text) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:     t.Type(),
		Content:  string(t),
		Metadata: map[string]interface{}{},
	}, nil
}

// MARK: Link
// Link represents a hyperlink element with display text and a target URL.
type Link struct {
	// Text holds the display text for the link
	Text string
	// Url holds the target URL that the link points to
	Url string
}

// Type returns the ContentType for this link element.
func (l Link) Type() ContentType { return LinkType }

// Materialize converts the link into a MaterializedContent with the display text
// as content and the URL stored in metadata under the "Url" key.
func (l Link) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:    l.Type(),
		Content: l.Text,
		Metadata: map[string]interface{}{
			"Url": l.Url,
		},
	}, nil
}

// MARK: Code
// Code represents inline code content in documentation. It is defined as a string type
// for short code snippets that appear within text (e.g., `variable` or `function()`).
type Code string

// Type returns the ContentType for this inline code element.
func (c Code) Type() ContentType {
	return CodeType
}

// Materialize converts the inline code into a MaterializedContent with its string content
// and an empty metadata map. Inline code is processed as-is without transformation.
func (c Code) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:     c.Type(),
		Content:  string(c),
		Metadata: map[string]interface{}{},
	}, nil
}

// MARK: Codeblock

// A codeblock is a NON-EXECUTABLE block of code
// Useful for examples/payloads etc
type CodeBlock struct {
	BlockType string
	Cmd       []string
}

func (c CodeBlock) Type() ContentType { return CodeBlockType }

func (c CodeBlock) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:    c.Type(),
		Content: strings.Join(c.Cmd, " "),
		Metadata: map[string]interface{}{
			"BlockType": c.BlockType,
		},
	}, nil
}

// MARK: BlockQuote
type BlockQuote string

func (b BlockQuote) Type() ContentType { return BlockQuoteType }

func (b BlockQuote) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:     b.Type(),
		Content:  string(b),
		Metadata: map[string]interface{}{},
	}, nil
}

// MARK: Executable

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

// MARK: Remote

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

// MARK: TableRow

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

// MARK: Comment
type Comment string

func (c Comment) Type() ContentType { return CommentType }

func (c Comment) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:     c.Type(),
		Content:  string(c),
		Metadata: map[string]interface{}{},
	}, nil
}
