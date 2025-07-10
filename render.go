package doyoucompute

import (
	"errors"
	"strings"
)

type Renderer[T any] interface {
	Render(node Node) (T, error)
}

type Markdown struct{}

func (m Markdown) writeHeader(builder *strings.Builder, content string, level int) {
	builder.WriteString(strings.Repeat("#", level))
	builder.WriteString(" ")
	builder.WriteString(content)
	builder.WriteString("\n\n")
}

func (m Markdown) renderStructureNode(structureNode Structurer) (string, error) {
	var documentContent strings.Builder

	if structureNode.Type() == DocumentType || structureNode.Type() == SectionType {
		var headerLevel int
		switch structureNode.Type() {
		case DocumentType:
			headerLevel = 1
		case SectionType:
			headerLevel = 2
		}

		m.writeHeader(&documentContent, structureNode.Identifer(), headerLevel)
	}

	for _, leaf := range structureNode.Children() {
		leafContent, err := m.Render(leaf)
		if err != nil {
			return "", err
		}

		documentContent.WriteString(leafContent)

		if structureNode.Type() == ParagraphType {
			documentContent.WriteString(" ")
		}

	}

	documentContent.WriteString("\n\n")

	return documentContent.String(), nil
}

func (m Markdown) renderHeader(header Header) (string, error) {
	content, err := header.Materialize()
	if err != nil {
		return "", nil
	}

	headerLevel := content.Metadata["Level"].(int)

	var headerContent strings.Builder

	m.writeHeader(&headerContent, content.Content, headerLevel)

	return headerContent.String(), nil
}

func (m Markdown) renderContent(contentNode Contenter) (string, error) {
	switch contentNode.Type() {
	case HeaderType:
		return m.renderHeader(contentNode.(Header))
	case LinkType:
		return "", nil
	case TextType:
		content, err := contentNode.Materialize()
		if err != nil {
			return "", err
		}

		return content.Content, nil
	case CodeType:
		return "", nil
	case CodeBlockType:
		return "", nil
	case BlockQuoteType:
		return "", nil
	case ExecutableType:
		content, err := contentNode.Materialize()
		if err != nil {
			return "", err
		}

		shell := content.Metadata["Shell"].(string)

		var builder strings.Builder

		builder.WriteString("```")
		builder.WriteString(shell)
		builder.WriteString("\n")
		builder.WriteString(content.Content)
		builder.WriteString("\n")
		builder.WriteString("```")
		builder.WriteString("\n")

		return builder.String(), nil
	case RemoteType:
		content, err := contentNode.Materialize()
		if err != nil {
			return "", err
		}

		var builder strings.Builder
		builder.WriteString(content.Content)
		builder.WriteString("\n")

		return builder.String(), nil
	}

	return "", errors.New("unknown content node type")
}

func (m Markdown) Render(node Node) (string, error) {
	switch node.Type() {
	case DocumentType, SectionType, ParagraphType, ListType:
		return m.renderStructureNode(node.(Structurer))
	default: // let the content check through an error for invalid type
		return m.renderContent(node.(Contenter))
	}
}
