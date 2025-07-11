package doyoucompute

import (
	"errors"
	"log"
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

	// Iterate through each leaf and render it recursively
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

	// Separate each section, paragraph, and list by a new line
	// We don't need to append newlines at the end of the document
	if structureNode.Type() == ParagraphType {
		log.Print(structureNode.Type(), structureNode.Identifer())
		documentContent.WriteString("\n\n")
	}

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

func (m Markdown) renderLink(link Link) (string, error) {
	content, err := link.Materialize()
	if err != nil {
		return "", err
	}

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

func (m Markdown) renderText(text Text) (string, error) {
	content, err := text.Materialize()
	if err != nil {
		return "", err
	}

	return content.Content, nil
}

func (m Markdown) renderCode(code Code) (string, error) {
	content, err := code.Materialize()
	if err != nil {
		return "", err
	}

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

func (m Markdown) renderCodeBlock(codeBlock CodeBlock) (string, error) {
	content, err := codeBlock.Materialize()
	if err != nil {
		return "", err
	}

	shell := content.Metadata["BlockType"].(string)

	var builder strings.Builder

	m.renderBlockofCode(shell, content.Content, &builder)

	return builder.String(), nil
}

func (m Markdown) renderBlockQuote(blockQuote BlockQuote) (string, error) {
	content, err := blockQuote.Materialize()
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	builder.WriteString("> ")
	builder.WriteString(content.Content)
	builder.WriteString("\n\n")

	return builder.String(), nil
}

func (m Markdown) renderExecutable(executable Executable) (string, error) {
	content, err := executable.Materialize()
	if err != nil {
		return "", err
	}

	shell := content.Metadata["Shell"].(string)

	var builder strings.Builder

	m.renderBlockofCode(shell, content.Content, &builder)

	return builder.String(), nil
}

func (m Markdown) renderRemoteContent(remote Remote) (string, error) {
	content, err := remote.Materialize()
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString(content.Content)
	builder.WriteString("\n")

	return builder.String(), nil
}

func (m Markdown) renderContent(contentNode Contenter) (string, error) {

	switch contentNode.Type() {
	case HeaderType:
		return m.renderHeader(contentNode.(Header))
	case LinkType:
		return m.renderLink(contentNode.(Link))
	case TextType:
		return m.renderText(contentNode.(Text))
	case CodeType:
		return m.renderCode(contentNode.(Code))
	case CodeBlockType:
		return m.renderCodeBlock(contentNode.(CodeBlock))
	case BlockQuoteType:
		return m.renderBlockQuote(contentNode.(BlockQuote))
	case ExecutableType:
		return m.renderExecutable(contentNode.(Executable))
	case RemoteType:
		return m.renderRemoteContent(contentNode.(Remote))
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
