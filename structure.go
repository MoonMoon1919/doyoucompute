package doyoucompute

type Table struct {
	Headers []string
	Items   []Node
}

func (t Table) Type() ContentType { return TableType }

func (t Table) Children() []Node { return t.Items }

func (t Table) Identifer() string { return "" }

// A container that allows us to render content with list semantics (optionally ordered)
type List struct {
	Items   []Node
	ordered bool
}

func (l List) Type() ContentType { return ListType }

func (l List) Children() []Node { return l.Items }

func (l List) Identifer() string { return "" }

// A container that allows us to render content with paragraph semantics
type Paragraph struct {
	Items []Node
}

func (p Paragraph) Type() ContentType { return ParagraphType }

func (p Paragraph) Children() []Node { return p.Items }

func (p Paragraph) Identifer() string { return "" }

// A single section has a name and 1..N items of content
type Section struct {
	Name    string
	Content []Node
}

func (s Section) Children() []Node { return s.Content }

func (s Section) Type() ContentType { return SectionType }

func (s Section) Identifer() string { return s.Name }

// A document contains many sections
type Document struct {
	Name    string
	Content []Node
}

func (d Document) Type() ContentType { return DocumentType }

func (d Document) Children() []Node { return d.Content }

func (d Document) Identifer() string { return d.Name }
