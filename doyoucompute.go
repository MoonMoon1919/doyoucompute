package doyoucompute

import (
	"io"
	"strings"
)

/*
This module contains the core domain logic for doyoucompute

All render methods should return all generated content,
rather than a string formatted as markdown.

Later, we can add a markdown formatter and script running formatter
*/

type Contenter interface {
	Render() (string, error)
}

type Header struct {
	Content string
	Level   int
}

func (h Header) Render() (string, error) {
	headerLevel := strings.Repeat("#", h.Level)

	return headerLevel + " " + h.Content, nil
}

// Content types
type FreeText string

func (f FreeText) Render() (string, error) {
	return string(f), nil
}

type Remote struct { // e.g., from local file in docs folder, from GitHub.. etc
	reader io.Reader
}

func (r Remote) Render() (string, error) {
	content, err := io.ReadAll(r.reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// A single section has a name and 1..N items of content
type Section struct {
	Name    string
	Content []Contenter
}

func (i Section) Render() (string, error) {
	var joinedString string

	name := Header{Content: i.Name, Level: 1}
	nameStr, _ := name.Render()
	joinedString = joinedString + nameStr

	for _, item := range i.Content {
		content, _ := item.Render()

		joinedString = joinedString + "\n\n" + content
	}

	return joinedString, nil
}

// A document contains all the things
