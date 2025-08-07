package doyoucompute

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Renderer defines a generic interface for converting Node elements into different output formats.
// The type parameter T allows for different renderers to produce different output types
// (e.g., string for markdown, []CommandPlan for execution planning).
type Renderer[T any] interface {
	// Render processes a node and converts it to the target format T.
	// Returns an error if the rendering process fails.
	Render(node Node) (T, error)
}

func getStringFromMetadata(metadata map[string]interface{}, key string) (string, error) {
	val, exists := metadata[key]
	if !exists {
		return "", fmt.Errorf("missing metadata key: %s", key)
	}
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("metadata key %s is not a string", key)
	}
	return str, nil
}

func getStringsFromMetadata(metadata map[string]interface{}, key string) ([]string, error) {
	val, exists := metadata[key]
	if !exists {
		return nil, fmt.Errorf("missing metadata key: %s", key)
	}

	// Handle []string
	if strs, ok := val.([]string); ok {
		return strs, nil
	}

	// Handle []interface{} and convert to []string
	if interfaces, ok := val.([]interface{}); ok {
		strs := make([]string, len(interfaces))
		for i, item := range interfaces {
			if str, ok := item.(string); ok {
				strs[i] = str
			} else {
				return nil, fmt.Errorf("metadata key %s contains non-string element at index %d: %T", key, i, item)
			}
		}
		return strs, nil
	}

	return nil, fmt.Errorf("metadata key %s expected []string, got %T", key, val)
}

// MARK: Tracking

// SectionInfo contains metadata about a section's position within the document hierarchy.
type SectionInfo struct {
	// Name is the identifier of the section
	Name string
	// Level indicates the nesting depth of the section (1 for top-level, 2 for subsection, etc.)
	Level int
}

// ContextPath represents a stack of section information that tracks the current
// position within a nested document structure during traversal or rendering.
type ContextPath []SectionInfo

// Push adds a new section to the context path and returns the updated path.
// The level is automatically calculated based on the current depth.
func (c ContextPath) Push(name string) ContextPath {
	level := len(c) + 1
	return append(c, SectionInfo{Name: name, Level: level})
}

// Current returns the SectionInfo for the current (most recent) section.
// Returns an empty SectionInfo if the path is empty.
func (c ContextPath) Current() SectionInfo {
	if len(c) == 0 {
		return SectionInfo{}
	}

	return c[len(c)-1]
}

// CurrentSection returns the name of the current section.
// Returns an empty string if no sections are in the path.
func (c ContextPath) CurrentSection() string {
	if len(c) == 0 {
		return ""
	}

	return c.Current().Name
}

// CurrentLevel returns the nesting level of the current section.
// Returns -1 if no sections are in the path.
func (c ContextPath) CurrentLevel() int {
	if len(c) == 0 {
		return -1
	}

	return c.Current().Level
}

// MARK: Markdown

// Markdown implements the Renderer interface to convert document nodes into markdown format.
// It handles hierarchical document structures and maintains proper heading levels during traversal.
type Markdown struct{}

// NewMarkdownRenderer creates a new Markdown renderer instance.
func NewMarkdownRenderer() Markdown {
	return Markdown{}
}

func (m Markdown) writeHeader(builder *strings.Builder, content string, level int) {
	fmt.Fprintf(builder, "%s %s\n\n", strings.Repeat("#", level), content)
}

func (m Markdown) renderChildren(children []Node, contextPath *ContextPath) ([]string, error) {
	if len(children) == 0 {
		return nil, nil
	}

	results := make([]string, len(children))

	for idx, leaf := range children {
		leafContent, err := m.renderWithTracking(leaf, contextPath)
		if err != nil {
			return nil, err
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

	return builder.String(), nil
}

func (m Markdown) renderDocument(d *Document, contextPath *ContextPath) (string, error) {
	ctxPath := contextPath.Push(d.Identifier())
	contextPath = &ctxPath // Update the context path so as we walk the tree we correctly track header level

	childContent, err := m.renderChildren(d.Children(), contextPath)
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	// Don't exceed an H5
	level := ctxPath.CurrentLevel()
	if ctxPath.CurrentLevel() > 5 {
		level = 5
	}

	if d.HasFrontmatter() {
		frontmatter, err := m.renderFrontmatter(d.Frontmatter)
		if err != nil {
			return "", err
		}

		builder.WriteString(frontmatter)
	}

	m.writeHeader(&builder, d.Identifier(), level)
	builder.WriteString(strings.Join(childContent, "\n\n"))

	// Final newline
	builder.WriteString("\n")

	return builder.String(), nil
}

func (m Markdown) renderSection(s Structurer, contextPath *ContextPath) (string, error) {
	ctxPath := contextPath.Push(s.Identifier())
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

	m.writeHeader(&builder, s.Identifier(), level)
	builder.WriteString(strings.Join(childContent, "\n\n"))

	return builder.String(), nil
}

func (m Markdown) renderTable(t *Table, contextPath *ContextPath) (string, error) {
	var builder strings.Builder

	joiner := strings.Join(t.Headers, " | ")

	// Header row
	builder.WriteString("| ")
	builder.WriteString(joiner)
	builder.WriteString(" |")
	builder.WriteString("\n")

	// Header row separator
	numSeparators := len(t.Headers) - 1
	numDividers := len(t.Headers)

	builder.WriteString("| ")

	for idx := range numDividers {
		builder.WriteString("----")

		if idx < numSeparators {
			builder.WriteString(" | ")
		}
	}

	builder.WriteString(" |")
	builder.WriteString("\n")

	// Children
	childContent, err := m.renderChildren(t.Children(), contextPath)
	if err != nil {
		return "", err
	}

	builder.WriteString(strings.Join(childContent, "\n"))

	return builder.String(), nil
}

func (m Markdown) renderList(l *List, contextPath *ContextPath) (string, error) {
	var builder strings.Builder

	childContent, err := m.renderChildren(l.Children(), contextPath)
	if err != nil {
		return "", err
	}

	for _, item := range childContent {
		builder.WriteString(l.TypeOfList.Prefix())
		builder.WriteString(" ")
		builder.WriteString(item)
		builder.WriteString("\n")
	}

	return builder.String(), nil
}

func (m Markdown) renderFrontmatter(f Frontmatter) (string, error) {
	var builder strings.Builder

	builder.WriteString("---\n")

	data, err := yaml.Marshal(f.Data)
	if err != nil {
		return "", err
	}
	builder.Write(data)
	builder.WriteString("\n")
	builder.WriteString("---\n\n")

	return builder.String(), nil
}

func (m Markdown) renderStructureNode(structureNode Structurer, contextPath *ContextPath) (string, error) {
	switch structureNode.Type() {
	case DocumentType:
		return m.renderDocument(structureNode.(*Document), contextPath)
	case SectionType:
		return m.renderSection(structureNode, contextPath)
	case ParagraphType:
		return m.renderParagraph(structureNode, contextPath)
	case ListType:
		return m.renderList(structureNode.(*List), contextPath)
	case TableType:
		return m.renderTable(structureNode.(*Table), contextPath)
	}

	return "", errors.New("unhandled structure node type")
}

func (m Markdown) renderHeader(content MaterializedContent, contextPath *ContextPath) (string, error) {
	var headerContent strings.Builder

	m.writeHeader(&headerContent, content.Content, contextPath.CurrentLevel())

	return headerContent.String(), nil
}

func (m Markdown) renderLink(content MaterializedContent) (string, error) {
	url, err := getStringFromMetadata(content.Metadata, "Url")
	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("[%s](%s)", content.Content, url), nil
}

func (m Markdown) renderText(content MaterializedContent) (string, error) {
	return content.Content, nil
}

func (m Markdown) renderCode(content MaterializedContent) (string, error) {
	return fmt.Sprintf("`%s`", content.Content), nil
}

func (m Markdown) renderBlockofCode(typeHint string, content string, builder *strings.Builder) {
	builder.WriteString("```")
	builder.WriteString(typeHint)
	builder.WriteString("\n")
	builder.WriteString(content)
	builder.WriteString("\n")
	builder.WriteString("```")
}

func (m Markdown) renderCodeBlock(content MaterializedContent) (string, error) {
	shell, err := getStringFromMetadata(content.Metadata, "BlockType")
	if err != nil {
		return "", nil
	}

	var builder strings.Builder

	m.renderBlockofCode(shell, content.Content, &builder)

	return builder.String(), nil
}

func (m Markdown) renderBlockQuote(content MaterializedContent) (string, error) {
	return fmt.Sprintf("> %s", content.Content), nil
}

func (m Markdown) renderExecutable(content MaterializedContent) (string, error) {
	shell, err := getStringFromMetadata(content.Metadata, "Shell")
	if err != nil {
		return "", nil
	}

	var builder strings.Builder

	m.renderBlockofCode(shell, content.Content, &builder)

	return builder.String(), nil
}

func (m Markdown) renderTableRow(content MaterializedContent) (string, error) {
	var builder strings.Builder

	items, err := getStringsFromMetadata(content.Metadata, "Items")
	if err != nil {
		return "", err
	}

	joiner := strings.Join(items, " | ")

	builder.WriteString("| ")
	builder.WriteString(joiner)
	builder.WriteString(" |")

	return builder.String(), nil
}

func (m Markdown) renderRemoteContent(content MaterializedContent) (string, error) {
	var builder strings.Builder
	builder.WriteString(content.Content)

	return builder.String(), nil
}

func (m Markdown) renderComment(content MaterializedContent) (string, error) {
	return fmt.Sprintf("<!-- %s -->", content.Content), nil
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
	case TableRowType:
		return m.renderTableRow(content)
	case RemoteType:
		return m.renderRemoteContent(content)
	case CommentType:
		return m.renderComment(content)
	}

	return "", errors.New("unknown content node type")
}

func (m Markdown) renderWithTracking(node Node, contextPath *ContextPath) (string, error) {
	switch node.Type() {
	case DocumentType, SectionType, ParagraphType, ListType, TableType, FrontmatterType:
		return m.renderStructureNode(node.(Structurer), contextPath)
	default: // let the content renderer check through an error for invalid type
		return m.renderContent(node.(Contenter), contextPath)
	}
}

// Render converts a document node into markdown format, starting with an empty context path.
// This is the main entry point for the Renderer interface implementation.
func (m Markdown) Render(node Node) (string, error) {
	return m.renderWithTracking(node, &ContextPath{})
}

// MARK: Executor

// CommandPlan represents a single executable command with its context information,
// used for planning and executing runnable documentation scripts.
type CommandPlan struct {
	// Shell specifies the shell or interpreter to use for execution
	Shell string
	// Args contains the command and its arguments to be executed
	Args []string
	// Context provides information about which section this command originated from
	Context SectionInfo
	// Environment variables that must be set for the command to be executed
	Environment []string
}

// Executioner implements the Renderer interface to extract executable commands
// from document nodes and create execution plans for runnable documentation.
type Executioner struct{}

// NewExecutionRenderer creates a new Executioner instance for building command execution plans.
func NewExecutionRenderer() Executioner {
	return Executioner{}
}

func (e Executioner) renderChildren(node Structurer, contextPath *ContextPath) ([]CommandPlan, error) {
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

func (e Executioner) renderStructureNode(node Structurer, contextPath *ContextPath) ([]CommandPlan, error) {
	ctxPath := contextPath.Push(node.Identifier())

	return e.renderChildren(node, &ctxPath)
}

func (e Executioner) renderExecutable(content MaterializedContent, contextPath *ContextPath) (CommandPlan, error) {
	shell, err := getStringFromMetadata(content.Metadata, "Shell")
	if err != nil {
		return CommandPlan{}, nil
	}

	args, err := getStringsFromMetadata(content.Metadata, "Command")
	if err != nil {
		return CommandPlan{}, err
	}

	envvars, err := getStringsFromMetadata(content.Metadata, "Environment")
	if err != nil {
		return CommandPlan{}, err
	}

	return CommandPlan{
		Shell:       shell,
		Args:        args,
		Context:     contextPath.Current(),
		Environment: envvars,
	}, nil
}

func (e Executioner) renderWithTracking(node Node, contextPath *ContextPath) ([]CommandPlan, error) {
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

// Render traverses a document node and extracts all executable commands,
// returning them as a slice of CommandPlan for execution planning.
// This is the main entry point for the Renderer interface implementation.
func (e Executioner) Render(node Node) ([]CommandPlan, error) {
	cmds, err := e.renderWithTracking(node, &ContextPath{})
	if err != nil {
		return make([]CommandPlan, 0), err
	}

	return cmds, nil
}
