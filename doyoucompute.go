package doyoucompute

/*
This module contains the core domain logic for doyoucompute

All render methods should return all generated content,
rather than a string formatted as markdown.

Later, we can add a markdown formatter and script running formatter
*/

type Node interface {
	Type() ContentType
}

type Contenter interface {
	Node
	Materialize() (MaterializedContent, error)
}

type Structurer interface {
	Node
	Identifier() string
	Children() []Node
}
