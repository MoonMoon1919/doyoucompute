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

func (m Markdown) renderChildren(children []Node) ([]string, error) {
	results := make([]string, len(children))

	for idx, leaf := range children {
		leafContent, err := m.Render(leaf)
		if err != nil {
			return []string{}, err
		}

		results[idx] = leafContent
	}

	return results, nil
}

func (m Markdown) renderParagraph(p Structurer) (string, error) {
	childContent, err := m.renderChildren(p.Children())
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	builder.WriteString(strings.Join(childContent, " "))
	builder.WriteString("\n\n")

	return builder.String(), nil
}

func (m Markdown) renderHeaderedPortion(s Structurer, headerLevel int) (string, error) {
	childContent, err := m.renderChildren(s.Children())
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	m.writeHeader(&builder, s.Identifer(), headerLevel)
	builder.WriteString(strings.Join(childContent, ""))

	return builder.String(), nil
}

func (m Markdown) renderStructureNode(structureNode Structurer) (string, error) {
	switch structureNode.Type() {
	case DocumentType:
		return m.renderHeaderedPortion(structureNode, 1)
	case SectionType:
		return m.renderHeaderedPortion(structureNode, 2)
	case ParagraphType:
		return m.renderParagraph(structureNode)
	case ListType:
		// TODO: Add support for List type here..
		return "", errors.New("not implemented")
	}

	return "", errors.New("unhandled structure node type")
}

func (m Markdown) renderHeader(content MaterializedContent) (string, error) {
	headerLevel := content.Metadata["Level"].(int)

	var headerContent strings.Builder

	m.writeHeader(&headerContent, content.Content, headerLevel)

	return headerContent.String(), nil
}

func (m Markdown) renderLink(content MaterializedContent) (string, error) {
	url := content.Metadata["Url"].(string)

	var builder strings.Builder

	builder.WriteString("[")
	builder.WriteString(content.Content)
	builder.WriteString("]")
	builder.WriteString("(")
	builder.WriteString(url)
	builder.WriteString(")")

	return builder.String(), nil
}

func (m Markdown) renderText(content MaterializedContent) (string, error) {
	return content.Content, nil
}

func (m Markdown) renderCode(content MaterializedContent) (string, error) {
	var builder strings.Builder

	builder.WriteString("`")
	builder.WriteString(content.Content)
	builder.WriteString("`")

	return builder.String(), nil
}

func (m Markdown) renderBlockofCode(typeHint string, content string, builder *strings.Builder) {
	builder.WriteString("```")
	builder.WriteString(typeHint)
	builder.WriteString("\n")
	builder.WriteString(content)
	builder.WriteString("\n")
	builder.WriteString("```")
	builder.WriteString("\n\n")
}

func (m Markdown) renderCodeBlock(content MaterializedContent) (string, error) {
	shell := content.Metadata["BlockType"].(string)

	var builder strings.Builder

	m.renderBlockofCode(shell, content.Content, &builder)

	return builder.String(), nil
}

func (m Markdown) renderBlockQuote(content MaterializedContent) (string, error) {
	var builder strings.Builder

	builder.WriteString("> ")
	builder.WriteString(content.Content)
	builder.WriteString("\n\n")

	return builder.String(), nil
}

func (m Markdown) renderExecutable(content MaterializedContent) (string, error) {
	shell := content.Metadata["Shell"].(string)

	var builder strings.Builder

	m.renderBlockofCode(shell, content.Content, &builder)

	return builder.String(), nil
}

func (m Markdown) renderRemoteContent(content MaterializedContent) (string, error) {
	var builder strings.Builder
	builder.WriteString(content.Content)
	builder.WriteString("\n")

	return builder.String(), nil
}

func (m Markdown) renderContent(contentNode Contenter) (string, error) {
	content, err := contentNode.Materialize()
	if err != nil {
		return "", err
	}

	switch contentNode.Type() {
	case HeaderType:
		return m.renderHeader(content)
	case LinkType:
		return m.renderLink(content)
	case TextType:
		return m.renderText(content)
	case CodeType:
		return m.renderCode(content)
	case CodeBlockType:
		return m.renderCodeBlock(content)
	case BlockQuoteType:
		return m.renderBlockQuote(content)
	case ExecutableType:
		return m.renderExecutable(content)
	case RemoteType:
		return m.renderRemoteContent(content)
	}

	return "", errors.New("unknown content node type")
}

func (m Markdown) Render(node Node) (string, error) {
	switch node.Type() {
	case DocumentType, SectionType, ParagraphType, ListType:
		return m.renderStructureNode(node.(Structurer))
	default: // let the content renderer check through an error for invalid type
		return m.renderContent(node.(Contenter))
	}
}
