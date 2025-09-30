package doyoucompute

// Finalizer is a function that performs post-processing or validation after all options
// have been applied to a configuration object. Finalizers are useful when you need to:
//   - Validate relationships between multiple fields
//   - Compute derived state based on the final configuration
//   - Perform cleanup or normalization after all options are applied
//
// Example use case: validating that a GitHub issue's name and frontmatter are consistent.
type Finalizer[T any] func(p *T) error

// OptionBuilder is a function that modifies a configuration object and optionally returns
// a Finalizer to run after all options are applied. This two-phase approach allows:
//   - Immediate validation of individual option values
//   - Deferred validation of relationships between options
//   - Separation of concerns between option application and finalization
//
// If no finalizer is needed, return nil for the Finalizer.
//
// Example:
//
//	func WithName(name string) OptionBuilder[MyProps] {
//	    return func(p *MyProps) (Finalizer[MyProps], error) {
//	        if name == "" {
//	            return nil, errors.New("name cannot be empty")
//	        }
//	        p.name = name
//
//	        // Return a finalizer that validates name consistency
//	        return func(p *MyProps) error {
//	            if p.name != p.title {
//	                return errors.New("name and title must match")
//	            }
//	            return nil
//	        }, nil
//	    }
//	}
type OptionBuilder[T any] func(p *T) (Finalizer[T], error)

// SectionBuilder is a function that modifies a Section by adding content to it.
// Builders receive a pointer to the section and can add paragraphs, lists, code blocks,
// or any other content. They return an error if the modification fails.
//
// Example:
//
//	func AddDescription(s *Section) error {
//	    s.NewParagraph().Text("This is a description")
//	    return nil
//	}
type SectionBuilder func(s *Section) error

// DocumentApplier is a function that modifies a Document by adding sections or
// configuring document-level properties. Appliers receive a pointer to the document
// and return an error if the modification fails.
//
// Example:
//
//	func AddIntroSection(d *Document) error {
//	    intro := NewSection("Introduction")
//	    intro.NewParagraph().Text("Welcome to the project")
//	    d.AddSection(intro)
//	    return nil
//	}
type DocumentApplier func(d *Document) error

// SectionFactory creates a new Section with the given name and applies all provided
// builders to populate it with content. Builders are applied in order, and if any
// builder returns an error, the factory stops and returns that error.
//
// This is useful for creating reusable section templates with default content.
//
// Example:
//
//	section, err := SectionFactory("Setup",
//	    func(s *Section) error {
//	        s.NewParagraph().Text("Install dependencies:")
//	        return nil
//	    },
//	    func(s *Section) error {
//	        s.WriteCodeBlock("bash", []string{"npm install"}, true)
//	        return nil
//	    },
//	)
func SectionFactory(name string, contentFuncs ...SectionBuilder) (Section, error) {
	s := NewSection(name)

	for _, cFunc := range contentFuncs {
		if err := cFunc(&s); err != nil {
			return Section{}, err
		}
	}

	return s, nil
}

// DocumentFactory creates a new Document with the given name and applies all provided
// appliers to populate it with sections and content. Appliers are executed in order,
// and if any applier returns an error, the factory stops and returns that error.
//
// This is useful for creating document templates with predefined structure.
//
// Example:
//
//	doc, err := DocumentFactory("README",
//	    func(d *Document) error {
//	        d.WriteIntro().Text("Project overview")
//	        return nil
//	    },
//	    func(d *Document) error {
//	        section := NewSection("Setup")
//	        section.WriteCodeBlock("bash", []string{"make install"}, true)
//	        d.AddSection(section)
//	        return nil
//	    },
//	)
func DocumentFactory(name string, appliers ...DocumentApplier) (Document, error) {
	document, err := NewDocument(name)
	if err != nil {
		return Document{}, err
	}

	for _, applier := range appliers {
		if err := applier(&document); err != nil {
			return Document{}, err
		}
	}

	return document, nil
}

// ApplyOptions applies a sequence of option builders to a configuration object,
// then runs any finalizers returned by those options. This two-phase approach allows:
//
//  1. Options are applied in order, with immediate validation
//  2. Finalizers run after all options are applied, enabling cross-field validation
//
// If any option or finalizer returns an error, processing stops and that error is returned.
//
// This pattern is useful for building configurable document templates where some
// properties depend on others or require validation after all configuration is complete.
//
// Example:
//
//	type DocProps struct {
//	    title   string
//	    author  string
//	}
//
//	props := &DocProps{}
//	err := ApplyOptions(props,
//	    WithTitle("My Document"),
//	    WithAuthor("Jane Doe"),
//	)
func ApplyOptions[T any](props *T, opts ...OptionBuilder[T]) error {
	for _, opt := range opts {
		postEffect, err := opt(props)
		if err != nil {
			return err
		}

		if postEffect != nil {
			if err := postEffect(props); err != nil {
				return err
			}
		}
	}

	return nil
}
