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

// A container that allows us to render content with paragraph semantics
type Paragraph struct {
	Items []Node // how can we only allow content types here so people don't do awkward things
}

func (p Paragraph) Type() ContentType { return ParagraphType }

func (p Paragraph) Children() []Node { return p.Items }

func (p Paragraph) Identifer() string { return "" }

func (p Paragraph) Write(content Node) {
	p.Items = append(p.Items, content)
}

// A single section has a name and 1..N items of content
type Section struct {
	Name    string
	Content []Node
}

func (s Section) Children() []Node { return s.Content }

func (s Section) Type() ContentType { return SectionType }

func (s Section) Identifer() string { return s.Name }

func (s Section) AddSection(section Section) {
	// TODO: Ordering is important
	s.Content = append(s.Content, section)
}

func (s Section) WriteSection(name string) *Section {
	section := Section{Name: name}

	s.Content = append(s.Content, &section)

	return &section
}

// A document contains many sections
type Document struct {
	Name    string
	Content []Node
}

func (d Document) Type() ContentType { return DocumentType }

func (d Document) Children() []Node { return d.Content }

func (d Document) Identifer() string { return d.Name }

func (d Document) AddSection(section Section) {
	// TODO: Ordering is important
	d.Content = append(d.Content, section)
}

func (d Document) WriteSection(name string) *Section {
	section := Section{Name: name}

	d.Content = append(d.Content, &section)

	return &section
}
