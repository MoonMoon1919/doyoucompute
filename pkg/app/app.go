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

func cliBuilder(service *doyoucompute.Service, documents map[string]*doyoucompute.Document) *cli.Command {
	findDoc := func(documentName string) (*doyoucompute.Document, error) {
		doc, ok := documents[documentName]

		if !ok {
			return nil, errors.New("document not found")
		}

		return doc, nil
	}

	cmd := &cli.Command{
		Name:  "dycoctl",
		Usage: "CLI application for doyoucompute",
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

					if err := service.RenderFile(document, outpath); err != nil {
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

					if result, err := service.CompareFile(document, outpath); err != nil {
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

					if _, err := service.ExecuteScript(document, section); err != nil {
						return err
					}

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

					if results, err := service.PlanScriptExecution(document, section); err != nil {
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

					fmt.Print("Documents\n")
					fmt.Print("---------\n")
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
	documents map[string]*doyoucompute.Document
	service   *doyoucompute.Service
}

func New(documents []*doyoucompute.Document, service *doyoucompute.Service) *app {
	documentMap := map[string]*doyoucompute.Document{}

	for _, document := range documents {
		documentMap[document.Name] = document
	}

	return &app{
		documents: documentMap,
		service:   service,
	}
}

func (a *app) Run(args []string) error {
	cli := cliBuilder(a.service, a.documents)

	if err := cli.Run(context.Background(), args); err != nil {
		return err
	}

	return nil
}
