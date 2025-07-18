package doyoucompute

type Table struct {
	Headers []string
	Items   []TableRow
}

func (t Table) Type() ContentType { return TableType }

func (t Table) Children() []Node {
	nodes := make([]Node, len(t.Items))

	for idx, row := range t.Items {
		nodes[idx] = row
	}

	return nodes
}

func (t Table) Identifer() string { return "" }

func (t *Table) AddRow(row TableRow) {
	t.Items = append(t.Items, row)
}

// A container that allows us to render content with list semantics (optionally ordered)
type ListTypeE int

const (
	BULLET ListTypeE = iota + 1
	NUMBERED
)

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

type List struct {
	Items      []Node // Preferably Text, Link, Code
	TypeOfList ListTypeE
}

func (l List) Type() ContentType { return ListType }

func (l List) Children() []Node { return l.Items }

func (l List) Identifer() string { return "" }

func (l *List) Push(item Node) {
	l.Items = append([]Node{item}, l.Items...)
}

func (l *List) Append(item Node) {
	l.Items = append(l.Items, item)
}

// A container that allows us to render content with paragraph semantics
type Paragraph struct {
	Items []Node
}

func NewParagraph() *Paragraph {
	return &Paragraph{
		Items: make([]Node, 0),
	}
}

func (p Paragraph) Type() ContentType { return ParagraphType }

func (p Paragraph) Children() []Node { return p.Items }

func (p Paragraph) Identifer() string { return "" }

func (p *Paragraph) Text(val string) *Paragraph {
	p.Items = append(p.Items, Text(val))

	return p
}

func (p *Paragraph) Code(val string) *Paragraph {
	p.Items = append(p.Items, Code(val))

	return p
}

func (p *Paragraph) Link(text, url string) *Paragraph {
	p.Items = append(p.Items, Link{Text: text, Url: url})

	return p
}

// A single section has a name and 1..N items of content
type Section struct {
	Name    string
	Content []Node // Preferably Paragraph, Table, List, Remote, CodeBlock, Executable, BlockQuote, etc
}

func NewSection(name string) Section {
	return Section{
		Name:    name,
		Content: make([]Node, 0),
	}
}

func (s Section) Children() []Node { return s.Content }

func (s Section) Type() ContentType { return SectionType }

func (s Section) Identifer() string { return s.Name }

func (s *Section) AddIntro(content *Paragraph) {
	s.Content = append([]Node{content}, s.Content...)
}

func (s *Section) AddSection(section Section) {
	s.Content = append(s.Content, section)
}

func (s *Section) AddParagraph(paragraph Paragraph) {
	s.Content = append(s.Content, paragraph)
}

func (s *Section) AddTable(headers []string, rows []TableRow) {
	table := Table{Headers: headers, Items: rows}

	s.Content = append(s.Content, table)
}

func (s *Section) NewTable(headers []string) *Table {
	table := Table{Headers: headers, Items: make([]TableRow, 0)}

	s.Content = append(s.Content, &table)

	return &table
}

func (s *Section) AddList(listType ListTypeE, items []Node) {
	list := List{TypeOfList: listType, Items: items}

	s.Content = append(s.Content, list)
}

func (s *Section) NewList(listType ListTypeE) *List {
	list := List{TypeOfList: listType}

	s.Content = append(s.Content, &list)

	return &list
}

func (s *Section) AddCodeBlock(blockType string, cmd []string, executable bool) {
	var newContent Node

	if executable {
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

func (s *Section) AddBlockQuote(value string) {
	s.Content = append(s.Content, BlockQuote(value))
}

func (s *Section) AddRemoteContent(remote Remote) {
	s.Content = append(s.Content, remote)
}

// A document contains many sections
type Document struct {
	Name    string
	Content []Node // Preferably Paragraph and Section
}

func NewDocument(name string) Document {
	return Document{
		Name:    name,
		Content: make([]Node, 0),
	}
}

func (d Document) Type() ContentType { return DocumentType }

func (d Document) Children() []Node { return d.Content }

func (d Document) Identifer() string { return d.Name }

func (d *Document) AddIntro(content *Paragraph) {
	d.Content = append([]Node{content}, d.Content...)
}

func (d *Document) AddSection(section Section) {
	d.Content = append(d.Content, section)
}
