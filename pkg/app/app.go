package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/MoonMoon1919/doyoucompute"
	"github.com/urfave/cli/v3"
)

func cliBuilder(cliName string, service *doyoucompute.Service, documents map[string]doyoucompute.Document) *cli.Command {
	// helper function that looks up a document by name from the registered documents map.
	// returns an error if the document is not found.
	findDoc := func(documentName string) (doyoucompute.Document, error) {
		doc, ok := documents[documentName]

		if !ok {
			return doyoucompute.Document{}, errors.New("document not found")
		}

		return doc, nil
	}

	cmd := &cli.Command{
		Name:  cliName,
		Usage: "CLI for generating and running docs",
		Commands: []*cli.Command{
			{
				Name:  "render",
				Usage: "Render a document as markdown",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "path",
						Value: "README.md",
						Usage: "The path to which you want to write the document",
					},
					&cli.StringFlag{
						Name:  "doc-name",
						Usage: "The name of the document",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					outpath := c.String("path")
					name := c.String("doc-name")

					document, err := findDoc(name)
					if err != nil {
						return err
					}

					if err := service.RenderFile(&document, outpath); err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:  "compare",
				Usage: "Compares the content of a document with the content in a written file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "path",
						Value: "README.md",
						Usage: "The path to which you want to write the document",
					},
					&cli.StringFlag{
						Name:  "doc-name",
						Usage: "The name of the document",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					outpath := c.String("path")
					name := c.String("doc-name")

					document, err := findDoc(name)
					if err != nil {
						return err
					}

					if result, err := service.CompareFile(&document, outpath); err != nil {
						return err
					} else {

						if !result.Matches {
							return fmt.Errorf("Results do not match, file hash %s, content hash %s", result.FileHash, result.DocumentHash)
						}
					}

					return nil
				},
			},
			{
				Name:  "run",
				Usage: "Runs the document as a script",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "section",
						Value: doyoucompute.ALL_SECTIONS,
						Usage: "The specific section in a document you'd like to run",
					},
					&cli.StringFlag{
						Name:  "doc-name",
						Usage: "The name of the document",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					section := c.String("section")
					name := c.String("doc-name")

					document, err := findDoc(name)
					if err != nil {
						return err
					}

					results, err := service.ExecuteScript(&document, section)
					if err != nil {
						return fmt.Errorf("Failed to execute script: %w", err)
					}

					// Provide feedback on results
					var failedCount int

					for _, result := range results {
						if result.Status == doyoucompute.FAILED {
							failedCount++

							// Extract missing env vars from error message if it's an env validation error
							if strings.Contains(result.Error.Error(), "environment validation failed") {
								fmt.Printf("❌ Command failed in section '%s': %s\n", result.SectionName, result.Command)
								fmt.Printf("   Error: %v\n", result.Error)

								// Give helpful suggestion
								if strings.Contains(result.Error.Error(), "required environment variables not set") {
									fmt.Printf("   💡 Tip: Set the required environment variables and try again\n")
								}
							} else if strings.Contains(result.Error.Error(), "security validation failed") {
								fmt.Printf("❌ Command blocked for security in section '%s': %s\n", result.SectionName, result.Command)
								fmt.Printf("   Error: %v\n", result.Error)
							} else {
								fmt.Printf("❌ Command failed in section '%s': %s\n", result.SectionName, result.Command)
								fmt.Printf("   Error: %v\n", result.Error)
							}
							fmt.Println()
						} else {
							fmt.Printf("✅ Completed: %s (section: %s)\n", result.Command, result.SectionName)
						}
					}

					if failedCount > 0 {
						return fmt.Errorf("%d out of %d commands failed", failedCount, len(results))
					}

					fmt.Printf("🎉 All %d commands completed successfully!\n", len(results))

					return nil
				},
			},
			{
				Name:  "plan",
				Usage: "Shows the output of what would be run as a script for the document",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "section",
						Value: doyoucompute.ALL_SECTIONS,
						Usage: "The specific section in a document you'd like to run",
					},
					&cli.StringFlag{
						Name:  "doc-name",
						Usage: "The name of the document",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					section := c.String("section")
					name := c.String("doc-name")

					document, err := findDoc(name)
					if err != nil {
						return err
					}

					if results, err := service.PlanScriptExecution(&document, section); err != nil {
						return err
					} else {
						for _, result := range results {
							log.Printf("[Section: %s] - Command: '%s'", result.Context.Name, strings.Join(result.Args, " "))
						}
					}

					return nil
				},
			},
			{
				Name:  "list",
				Usage: "List all available docs",
				Action: func(ctx context.Context, c *cli.Command) error {
					if len(documents) == 0 {
						return errors.New("no documents found")
					}

					for docName := range documents {
						fmt.Print(docName)
						fmt.Print("\n")
					}

					return nil
				},
			},
		},
	}

	return cmd
}

type app struct {
	documents map[string]doyoucompute.Document
	service   *doyoucompute.Service
}

// New creates a new CLI application instance with the provided service and
// an empty document registry.
func New(service *doyoucompute.Service) *app {
	return &app{
		documents: map[string]doyoucompute.Document{},
		service:   service,
	}
}

// Register adds a document to the application's registry, making it available
// for CLI operations. The document is indexed by its Name field.
func (a *app) Register(document doyoucompute.Document) {
	a.documents[document.Name] = document
}

// Run executes the CLI application with the provided command-line arguments.
// This is the main entry point for the CLI functionality.
func (a *app) Run(args []string) error {
	cli := cliBuilder("dycoctl", a.service, a.documents)

	if err := cli.Run(context.Background(), args); err != nil {
		log.Fatal(err)
	}

	return nil
}
