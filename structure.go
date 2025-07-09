package doyoucompute

// A container that allows us to render content with list semantics (optionally ordered)
type List struct {
	Items   []Node
	ordered bool
}

func (l List) Type() ContentType {
	return ListType
}

func (l List) Children() []Node { return l.Items }

// A container that allows us to render content with paragraph semantics
type Paragraph struct {
	Items []Node
}

func (p Paragraph) Type() ContentType {
	return ParagraphType
}

func (p Paragraph) Children() []Node { return p.Items }

// A single section has a name and 1..N items of content
type Section struct {
	Name    string
	Content []Node
}

func (s Section) Children() []Node { return s.Content }

func (s Section) Type() ContentType {
	return SectionType
}

// A document contains many sections
type Document struct {
	Name     string
	Sections []Section
}

func (d Document) Type() ContentType {
	return DocumentType
}
