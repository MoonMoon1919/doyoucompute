package content

type ContentType int

const (
	HeaderType ContentType = iota + 1
	ParagraphType
	ListType
	LinkType
	CodeType
	CodeBlockType
	BlockQuoteType
	ExecutableType
	RemoteType
)

type MaterializedContent struct {
	Type     string
	Metadata map[string]interface{}
}

type Materializer interface {
	Materialize() (MaterializedContent, error)
}
