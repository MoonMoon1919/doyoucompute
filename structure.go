package doyoucompute

import (
	"errors"
)

// MARK: Frontmatter

// Frontmatter represents YAML/TOML metadata typically found at the beginning of
// markdown documents, containing key-value configuration and metadata pairs.
type Frontmatter struct {
	// Data holds the parsed frontmatter key-value pairs
	Data map[string]interface{}
}

// NewFrontmatter creates a new Frontmatter instance with the provided data map.
func NewFrontmatter(data map[string]interface{}) *Frontmatter {
	return &Frontmatter{
		Data: data,
	}
}

// Empty returns true if the frontmatter contains no data.
func (f *Frontmatter) Empty() bool {
	return len(f.Data) == 0
}

// MARK: Table

// Table represents a complete table structure with headers and multiple rows of data.
type Table struct {
	// Headers contains the column header names for the table
	Headers []string
	// Items contains all the rows of data in the table
	Items []TableRow
}

// NewTable creates a new Table instance with the specified headers and rows.
func NewTable(headers []string, items []TableRow) *Table {
	return &Table{
		Headers: headers,
		Items:   items,
	}
}

// Type returns the ContentType for this table element.
func (t Table) Type() ContentType { return TableType }

// Children returns all table rows as Node interfaces, allowing the table
// to be treated as a parent node in a document tree structure.
func (t Table) Children() []Node {
	nodes := make([]Node, len(t.Items))

	for idx, row := range t.Items {
		nodes[idx] = row
	}

	return nodes
}

// Identifier returns an empty string as tables do not have specific identifiers.
func (t Table) Identifier() string { return "" }

// AddRow appends a new row to the table with the provided column values.
// Returns an error if the number of values exceeds the number of headers.
func (t *Table) AddRow(row ...string) error {
	if len(row) > len(t.Headers) {
		return errors.New("Row length exceeds number of headers")
	}

	t.Items = append(t.Items, TableRow{Values: row})

	return nil
}

// MARK: List

// ListTypeE represents the different types of lists that can be rendered.
type ListTypeE int

const (
	// BULLET represents an unordered list with bullet points
	BULLET ListTypeE = iota + 1
	// NUMBERED represents an ordered list with numbers
	NUMBERED
)

// Prefix returns the string prefix used for rendering this list type.
// Returns "-" for bullet lists and "1." for numbered lists.
func (l ListTypeE) Prefix() string {
	switch l {
	case BULLET:
		return "-"
	case NUMBERED:
		return "1."
	}

	// Default to bulleted list
	return "-"
}

// List represents a container for rendering content as either bulleted or numbered lists.
type List struct {
	// Items contains the text content for each list item
	Items []Text
	// TypeOfList specifies whether this is a bulleted or numbered list
	TypeOfList ListTypeE
}

// NewList creates a new List instance of the specified type with an empty items slice.
func NewList(typeOfList ListTypeE) *List {
	return &List{
		TypeOfList: typeOfList,
		Items:      make([]Text, 0),
	}
}

// Type returns the ContentType for this list element.
func (l List) Type() ContentType { return ListType }

// Children returns all list items as Node interfaces, allowing the list
// to be treated as a parent node in a document tree structure.
func (l List) Children() []Node {
	nodes := make([]Node, len(l.Items))

	for idx, item := range l.Items {
		nodes[idx] = item
	}

	return nodes
}

// Identifier returns an empty string as lists do not have specific identifiers.
func (l List) Identifier() string { return "" }

// Push adds a new item to the beginning of the list.
func (l *List) Push(val string) {
	l.Items = append([]Text{Text(val)}, l.Items...)
}

// Append adds a new item to the end of the list.
func (l *List) Append(val string) {
	l.Items = append(l.Items, Text(val))
}

// MARK: Container

// Paragraph represents a container for mixed content elements that should be
// rendered together as a cohesive paragraph block.
type Paragraph struct {
	// Items contains the various content elements (text, code, links) within this paragraph
	Items []Node
}

// NewParagraph creates a new Paragraph instance with an empty items slice.
func NewParagraph() *Paragraph {
	return &Paragraph{
		Items: make([]Node, 0),
	}
}

// Type returns the ContentType for this paragraph element.
func (p Paragraph) Type() ContentType { return ParagraphType }

// Children returns all items within the paragraph as Node interfaces.
func (p Paragraph) Children() []Node { return p.Items }

// Identifier returns an empty string as paragraphs do not have specific identifiers.
func (p Paragraph) Identifier() string { return "" }

// Text adds a text element to the paragraph and returns the paragraph for method chaining.
func (p *Paragraph) Text(val string) *Paragraph {
	p.Items = append(p.Items, Text(val))

	return p
}

// Code adds an inline code element to the paragraph and returns the paragraph for method chaining.
func (p *Paragraph) Code(val string) *Paragraph {
	p.Items = append(p.Items, Code(val))

	return p
}

// Link adds a hyperlink element to the paragraph and returns the paragraph for method chaining.
func (p *Paragraph) Link(text, url string) *Paragraph {
	p.Items = append(p.Items, Link{Text: text, Url: url})

	return p
}

// MARK: Section

// Section represents a named container that holds various types of content elements,
// allowing for hierarchical document organization.
type Section struct {
	// Name is the identifier/title for this section
	Name string
	// Content holds all the content elements within this section
	Content []Node
}

// NewSection creates a new Section with the specified name and empty content.
func NewSection(name string) Section {
	return Section{
		Name:    name,
		Content: make([]Node, 0),
	}
}

// Children returns all content within the section as Node interfaces.
func (s Section) Children() []Node { return s.Content }

// Type returns the ContentType for this section element.
func (s Section) Type() ContentType { return SectionType }

// Identifier returns the section name as its identifier.
func (s Section) Identifier() string { return s.Name }

// AddIntro prepends a paragraph to the beginning of the section content.
func (s *Section) AddIntro(content *Paragraph) {
	s.Content = append([]Node{content}, s.Content...)
}

// WriteIntro creates a new paragraph at the beginning of the section and returns it for editing.
func (s *Section) WriteIntro() *Paragraph {
	paragraph := NewParagraph()

	s.Content = append([]Node{paragraph}, s.Content...)

	return paragraph
}

// AddSection appends an existing section as a subsection.
func (s *Section) AddSection(section Section) {
	s.Content = append(s.Content, section)
}

// CreateSection creates a new subsection with the given name and returns it for editing.
func (s *Section) CreateSection(name string) *Section {
	section := NewSection(name)

	s.Content = append(s.Content, &section)

	return &section
}

// AddParagraph appends an existing paragraph to the section.
func (s *Section) AddParagraph(paragraph Paragraph) {
	s.Content = append(s.Content, paragraph)
}

// WriteParagraph creates a new paragraph in the section and returns it for editing.
func (s *Section) WriteParagraph() *Paragraph {
	paragraph := NewParagraph()

	s.Content = append(s.Content, paragraph)

	return paragraph
}

// AddTable creates and adds a table with the specified headers and rows.
func (s *Section) AddTable(headers []string, rows []TableRow) {
	table := Table{Headers: headers, Items: rows}

	s.Content = append(s.Content, table)
}

// CreateTable creates a new table with the given headers and returns it for editing.
func (s *Section) CreateTable(headers []string) *Table {
	table := Table{Headers: headers, Items: make([]TableRow, 0)}

	s.Content = append(s.Content, &table)

	return &table
}

// AddList creates and adds a list of the specified type with the given items.
func (s *Section) AddList(listType ListTypeE, items []Text) {
	list := List{TypeOfList: listType, Items: items}

	s.Content = append(s.Content, list)
}

// CreateList creates a new list of the specified type and returns it for editing.
func (s *Section) CreateList(listType ListTypeE) *List {
	list := List{TypeOfList: listType}

	s.Content = append(s.Content, &list)

	return &list
}

// WriteCodeBlock adds either an executable or non-executable code block based on the executable parameter.
// If executable is Exec, creates an Executable; otherwise creates a CodeBlock.
func (s *Section) WriteCodeBlock(blockType string, cmd []string, executable CodeBlockExecType) {
	var newContent Node

	if executable == Exec {
		newContent = Executable{
			Shell: blockType,
			Cmd:   cmd,
		}
	} else {
		newContent = CodeBlock{
			BlockType: blockType,
			Cmd:       cmd,
		}
	}

	s.Content = append(s.Content, newContent)
}

// WriteBlockQuote adds a block quote with the specified content to the section.
func (s *Section) WriteBlockQuote(value string) {
	s.Content = append(s.Content, BlockQuote(value))
}

// WriteRemoteContent adds remote content to the section.
func (s *Section) WriteRemoteContent(remote Remote) {
	s.Content = append(s.Content, remote)
}

// WriteComment adds a comment to the section.
func (s *Section) WriteComment(value string) {
	s.Content = append(s.Content, Comment(value))
}

// MARK: Document

// Document represents the top-level container for a complete document with optional
// frontmatter metadata and structured content organized into sections and paragraphs.
type Document struct {
	// Name is the identifier/title for this document
	Name string
	// Frontmatter contains optional metadata for the document
	Frontmatter Frontmatter
	// Content holds all the content elements within this document
	Content []Node
}

// NewDocument creates a new Document with the specified name and empty content.
func NewDocument(name string) Document {
	return Document{
		Name:    name,
		Content: make([]Node, 0),
	}
}

// Type returns the ContentType for this document element.
func (d Document) Type() ContentType { return DocumentType }

// Children returns all content within the document as Node interfaces.
func (d Document) Children() []Node { return d.Content }

// Identifier returns the document name as its identifier.
func (d Document) Identifier() string { return d.Name }

// AddIntro prepends a paragraph to the beginning of the document content.
func (d *Document) AddIntro(content *Paragraph) {
	d.Content = append([]Node{content}, d.Content...)
}

// AddFrontmatter sets the frontmatter metadata for the document.
func (d *Document) AddFrontmatter(f Frontmatter) {
	d.Frontmatter = f
}

// HasFrontmatter returns true if the document has frontmatter data.
func (d *Document) HasFrontmatter() bool {
	if d.Frontmatter.Data != nil {
		return true
	}

	if !d.Frontmatter.Empty() {
		return true
	}

	return false
}

// WriteIntro creates a new paragraph at the beginning of the document and returns it for editing.
func (d *Document) WriteIntro() *Paragraph {
	paragraph := NewParagraph()

	d.Content = append([]Node{paragraph}, d.Content...)

	return paragraph
}

// AddSection appends an existing section to the document.
func (d *Document) AddSection(section Section) {
	d.Content = append(d.Content, section)
}

// CreateSection creates a new section with the given name and returns it for editing.
func (d *Document) CreateSection(name string) *Section {
	s := NewSection(name)

	d.Content = append(d.Content, &s)

	return &s
}
