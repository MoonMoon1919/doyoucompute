package doyoucompute

// Node represents the base interface for all content elements in the documentation system.
// It provides type identification for polymorphic handling of different content types.
type Node interface {
	// Type returns the ContentType that identifies what kind of content this node represents
	Type() ContentType
}

// Contenter represents content elements that can be materialized into their final form.
// This interface is implemented by leaf nodes that contain actual content data
// (text, code, links, etc.) that needs to be processed and rendered.
type Contenter interface {
	Node
	// Materialize converts the content into a MaterializedContent with rendered output
	// and associated metadata. Returns an error if materialization fails.
	Materialize() (MaterializedContent, error)
}

// Structurer represents container elements that organize and structure the document.
// This interface is implemented by nodes that can contain other nodes and provide
// hierarchical organization (sections, documents, tables, lists, etc.).
type Structurer interface {
	Node
	// Identifier returns a name for this structural element
	Identifier() string
	// Children returns all child nodes contained within this structural element
	Children() []Node
}
