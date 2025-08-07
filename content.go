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
// Headers are processed as-is without transformation
func (h Header) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:     h.Type(),
		Content:  h.Content,
		Metadata: map[string]interface{}{},
	}, nil
}

// MARK: Text

// Text represents plain text content in documentation.
type Text string

// Type returns the ContentType for this text element.
func (t Text) Type() ContentType { return TextType }

// Materialize converts the text into a MaterializedContent with its string content.
// Plain text is processed as-is without transformation.
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
// Inline code is processed as-is without transformation.
func (c Code) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:     c.Type(),
		Content:  string(c),
		Metadata: map[string]interface{}{},
	}, nil
}

// MARK: Codeblock

// CodeBlock represents a non-executable code block used for displaying examples,
// payloads, or other code snippets that should not be run during documentation generation.
type CodeBlock struct {
	// BlockType specifies the language or type of the code block (e.g., "json", "bash", "go")
	BlockType string
	// Cmd contains the lines or components of the code block content
	Cmd []string
}

// Type returns the ContentType for this code block element.
func (c CodeBlock) Type() ContentType { return CodeBlockType }

// Materialize converts the code block into a MaterializedContent by joining
// the Cmd slice with spaces as the content and storing the BlockType in metadata.
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

// BlockQuote represents quoted text content in documentation (e.g., > quoted text).
// It is defined as a string type for text that should be rendered as a block quote.
type BlockQuote string

// Type returns the ContentType for this block quote element.
func (b BlockQuote) Type() ContentType { return BlockQuoteType }

// Materialize converts the block quote into a MaterializedContent with its string content
// Block quotes are processed as-is without transformation.
func (b BlockQuote) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:     b.Type(),
		Content:  string(b),
		Metadata: map[string]interface{}{},
	}, nil
}

// MARK: Executable

// Executable represents a code block that can be executed while running documentation as a script.
// This is the core component that enables "runnable documentation" by storing commands
// that will be executed, while representing them as code blocks in rendered output.
type Executable struct {
	// Shell specifies the shell or interpreter to use for execution (e.g., "bash", "sh", "python")
	Shell string
	// Cmd contains the command and arguments to be executed
	Cmd []string
	// Environment variables that must be set for the command to be run
	Environment []string
}

// Type returns the ContentType for this executable element.
func (e Executable) Type() ContentType { return ExecutableType }

// Materialize converts the executable into a MaterializedContent with the joined command
// as content and execution metadata including the shell, executable flag, and original command.
func (e Executable) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:    e.Type(),
		Content: strings.Join(e.Cmd, " "),
		Metadata: map[string]interface{}{
			"Shell":       e.Shell,
			"Command":     e.Cmd,
			"Environment": e.Environment,
		},
	}, nil
}

// MARK: Remote

// Remote represents content that is sourced from external locations such as local files
// in a docs folder, GitHub repositories, or other remote sources.
type Remote struct {
	// Reader provides access to the remote content data
	Reader io.Reader
}

// Type returns the ContentType for this remote content element.
func (r Remote) Type() ContentType { return RemoteType }

// Materialize reads all content from the Reader and converts it into a MaterializedContent.
// Returns an error if the content cannot be read from the remote source.
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

// TableRow represents a single row in a table with multiple column values.
type TableRow struct {
	// Values contains the content for each column in this table row
	Values []string
}

// Type returns the ContentType for this table row element.
func (t TableRow) Type() ContentType { return TableRowType }

// Materialize converts the table row into a MaterializedContent with empty content
// and the row values stored in metadata under the "Items" key.
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

// Comment represents comment content in documentation that provides additional context
// or notes. It is defined as a string type for text that may be rendered differently
// from regular content.
type Comment string

// Type returns the ContentType for this comment element.
func (c Comment) Type() ContentType { return CommentType }

// Materialize converts the comment into a MaterializedContent with its string content
// Comments are processed as-is without transformation.
func (c Comment) Materialize() (MaterializedContent, error) {
	return MaterializedContent{
		Type:     c.Type(),
		Content:  string(c),
		Metadata: map[string]interface{}{},
	}, nil
}
