package main

import (
	"os"

	"github.com/MoonMoon1919/doyoucompute"
	"github.com/MoonMoon1919/doyoucompute/docs/pkg/documents"
	"github.com/MoonMoon1919/doyoucompute/pkg/app"
)

func main() {
	repo := doyoucompute.NewFileRepository()
	fileRenderer := doyoucompute.NewMarkdownRenderer()
	execRenderer := doyoucompute.NewExecutionRenderer()
	runner := doyoucompute.NewTaskRunner()
	svc := doyoucompute.NewService(repo, runner, fileRenderer, execRenderer)

	// Docs to register
	app := app.New(&svc)

	app.Register(documents.Readme())
	app.Register(documents.Contributing())
	app.Register(documents.BugReport())
	app.Register(documents.PullRequest())

	app.Run(os.Args)
}
