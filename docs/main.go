package main

import (
	"log"
	"os"

	"github.com/MoonMoon1919/doyoucompute"
	"github.com/MoonMoon1919/doyoucompute/docs/pkg/documents"
	"github.com/MoonMoon1919/doyoucompute/pkg/app"
)

func main() {
	repo := doyoucompute.NewFileRepository()
	fileRenderer := doyoucompute.NewMarkdownRenderer()
	execRenderer := doyoucompute.NewExecutionRenderer()
	runner := doyoucompute.NewTaskRunner(doyoucompute.DefaultSecureConfig())
	svc := doyoucompute.NewService(repo, runner, fileRenderer, execRenderer)

	// Docs to register
	app := app.New(&svc)

	readMe, err := documents.Readme()
	if err != nil {
		log.Fatal(err)
	}
	app.Register(readMe)

	contrib, err := documents.Contributing()
	if err != nil {
		log.Fatal(err)
	}
	app.Register(contrib)

	bugReport, err := documents.BugReport()
	if err != nil {
		log.Fatal(err)
	}
	app.Register(bugReport)

	prTemplate, err := documents.PullRequest()
	if err != nil {
		log.Fatal(err)
	}
	app.Register(prTemplate)

	app.Run(os.Args)
}
