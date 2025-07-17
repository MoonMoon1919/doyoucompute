package doyoucompute

type Table struct {
	Headers []string
	Items   []Node // being explicit about TableRowType here would be better but good for now
}

func (t Table) Type() ContentType { return TableType }

func (t Table) Children() []Node { return t.Items }

func (t Table) Identifer() string { return "" }

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
	Items      []Node
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
	Items []Node // how can we only allow content types here so people don't do awkward things
}

func NewParagraph() *Paragraph {
	return &Paragraph{
		Items: make([]Node, 0),
	}
}

func (p Paragraph) Type() ContentType { return ParagraphType }

func (p Paragraph) Children() []Node { return p.Items }

func (p Paragraph) Identifer() string { return "" }

func (p *Paragraph) Next(content Node) *Paragraph {
	p.Items = append(p.Items, content)
	return p
}

// A single section has a name and 1..N items of content
type Section struct {
	Name    string
	Content []Node
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

func (s *Section) AddTable(headers []string, rows []Node) {
	table := Table{Headers: headers, Items: rows}

	s.Content = append(s.Content, table)
}

func (s *Section) AddList(listType ListTypeE, items []Node) {
	list := List{TypeOfList: listType, Items: items}

	s.Content = append(s.Content, list)
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

func (s *Section) AddRemoteContent(value Node) {
	s.Content = append(s.Content, value)
}

// A document contains many sections
type Document struct {
	Name    string
	Content []Node
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
