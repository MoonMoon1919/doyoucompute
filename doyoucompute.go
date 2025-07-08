package doyoucompute

import (
	"github.com/MoonMoon1919/doyoucompute/pkg/content"
)

/*
This module contains the core domain logic for doyoucompute

All render methods should return all generated content,
rather than a string formatted as markdown.

Later, we can add a markdown formatter and script running formatter
*/

// A single section has a name and 1..N items of content
type Section struct {
	Name    string
	Content []content.Materializer
}

func (i Section) Materialize() (content.MaterializedContent, error) {
	// var joinedString string

	// name := content.Header{Content: i.Name, Level: 1}
	// nameStr, _ := name.Render()
	// joinedString = joinedString + nameStr

	// for _, item := range i.Content {
	// 	content, _ := item.Render()

	// 	joinedString = joinedString + "\n\n" + content
	// }

	// return joinedString, nil
	return content.MaterializedContent{}, nil
}

// A document contains all the things
type Document struct {
	Name     string
	Sections []Section
}
