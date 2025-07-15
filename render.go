package doyoucompute

import (
	"errors"
	"strings"
)

type Renderer[T any] interface {
	Render(node Node) (T, error)
}

// MARK: Tracking
type SectionInfo struct {
	Name  string
	Level int
}

type ContextPath []SectionInfo

func (c ContextPath) Push(name string) ContextPath {
	level := len(c) + 1
	return append(c, SectionInfo{Name: name, Level: level})
}

func (c ContextPath) Current() SectionInfo {
	if len(c) == 0 {
		return SectionInfo{}
	}

	return c[len(c)-1]
}

func (c ContextPath) CurrentSection() string {
	if len(c) == 0 {
		return ""
	}

	return c.Current().Name
}

func (c ContextPath) CurrentLevel() int {
	if len(c) == 0 {
		return -1
	}

	return c.Current().Level
}

// MARK: Markdown
type Markdown struct{}

func (m Markdown) writeHeader(builder *strings.Builder, content string, level int) {
	builder.WriteString(strings.Repeat("#", level))
	builder.WriteString(" ")
	builder.WriteString(content)
	builder.WriteString("\n\n")
}

func (m Markdown) renderChildren(children []Node, contextPath *ContextPath) ([]string, error) {
	results := make([]string, len(children))

	for idx, leaf := range children {
		leafContent, err := m.renderWithTracking(leaf, contextPath)
		if err != nil {
			return []string{}, err
		}

		results[idx] = leafContent
	}

	return results, nil
}

func (m Markdown) renderParagraph(p Structurer, contextPath *ContextPath) (string, error) {
	childContent, err := m.renderChildren(p.Children(), contextPath)
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	builder.WriteString(strings.Join(childContent, " "))
	builder.WriteString("\n\n")

	return builder.String(), nil
}

func (m Markdown) renderHeaderedPortion(s Structurer, contextPath *ContextPath) (string, error) {
	ctxPath := contextPath.Push(s.Identifer())
	contextPath = &ctxPath // Update the context path so as we walk the tree we correctly track header level

	childContent, err := m.renderChildren(s.Children(), contextPath)
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	// Don't exceed an H5
	level := ctxPath.CurrentLevel()
	if ctxPath.CurrentLevel() > 5 {
		level = 5
	}

	m.writeHeader(&builder, s.Identifer(), level)
	builder.WriteString(strings.Join(childContent, ""))

	return builder.String(), nil
}

func (m Markdown) renderStructureNode(structureNode Structurer, contextPath *ContextPath) (string, error) {
	switch structureNode.Type() {
	case DocumentType:
		return m.renderHeaderedPortion(structureNode, contextPath)
	case SectionType:
		return m.renderHeaderedPortion(structureNode, contextPath)
	case ParagraphType:
		return m.renderParagraph(structureNode, contextPath)
	case ListType:
		// TODO: Add support for List type here..
		return "", errors.New("not implemented")
	case TableType:
		// TODO: Add support for table type here..
		return "", errors.New("not implemented")
	}

	return "", errors.New("unhandled structure node type")
}

func (m Markdown) renderHeader(content MaterializedContent, contextPath *ContextPath) (string, error) {
	var headerContent strings.Builder

	m.writeHeader(&headerContent, content.Content, contextPath.CurrentLevel())

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

func (m Markdown) renderContent(contentNode Contenter, contextPath *ContextPath) (string, error) {
	content, err := contentNode.Materialize()
	if err != nil {
		return "", err
	}

	switch contentNode.Type() {
	case HeaderType:
		return m.renderHeader(content, contextPath)
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
	case TableRowTable:
		// TODO
	case RemoteType:
		return m.renderRemoteContent(content)
	}

	return "", errors.New("unknown content node type")
}

func (m Markdown) renderWithTracking(node Node, contextPath *ContextPath) (string, error) {
	switch node.Type() {
	case DocumentType, SectionType, ParagraphType, ListType, TableType:
		return m.renderStructureNode(node.(Structurer), contextPath)
	default: // let the content renderer check through an error for invalid type
		return m.renderContent(node.(Contenter), contextPath)
	}
}

func (m Markdown) Render(node Node) (string, error) {
	return m.renderWithTracking(node, &ContextPath{})
}

// MARK: Executor
type CommandPlan struct {
	Shell   string
	Args    []string
	Context SectionInfo
}

type ExecutionPlan struct {
	Commands []CommandPlan
}

func (e ExecutionPlan) renderChildren(node Structurer, contextPath *ContextPath) ([]CommandPlan, error) {
	var commands []CommandPlan

	for _, leaf := range node.Children() {
		cmds, err := e.renderWithTracking(leaf, contextPath)
		if err != nil {
			return make([]CommandPlan, 0), err
		}

		commands = append(commands, cmds...)
	}

	return commands, nil
}

func (e ExecutionPlan) renderStructureNode(node Structurer, contextPath *ContextPath) ([]CommandPlan, error) {
	ctxPath := contextPath.Push(node.Identifer())

	return e.renderChildren(node, &ctxPath)
}

func (e ExecutionPlan) renderExecutable(content MaterializedContent, contextPath *ContextPath) (CommandPlan, error) {
	shell := content.Metadata["Shell"].(string)
	args := content.Metadata["Command"].([]string)

	return CommandPlan{
		Shell:   shell,
		Args:    args,
		Context: contextPath.Current(),
	}, nil
}

func (e ExecutionPlan) renderWithTracking(node Node, contextPath *ContextPath) ([]CommandPlan, error) {
	var commands []CommandPlan

	switch node.Type() {
	// Intentionally skip paragraphs and tables
	// Lists _could_ have executables as items
	case DocumentType, SectionType, ListType:
		cmds, err := e.renderStructureNode(node.(Structurer), contextPath)
		if err != nil {
			return []CommandPlan{}, nil
		}

		commands = append(commands, cmds...)

	// We only care to track executables for building execution plans
	case ExecutableType:
		content, err := node.(Contenter).Materialize()
		if err != nil {
			return []CommandPlan{}, err
		}

		cmd, err := e.renderExecutable(content, contextPath)
		if err != nil {
			return []CommandPlan{}, err
		}

		commands = append(commands, cmd)
	}

	return commands, nil
}

func (e *ExecutionPlan) Render(node Node) ([]CommandPlan, error) {
	cmds, err := e.renderWithTracking(node, &ContextPath{})
	if err != nil {
		return make([]CommandPlan, 0), err
	}

	e.Commands = append(e.Commands, cmds...)

	return e.Commands, nil
}
