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

	readme := documents.Readme()
	contribution := documents.Contributing()
	bugreport := documents.BugReport()
	pullrequest := documents.PullRequest()

	app.Register(readme)
	app.Register(contribution)
	app.Register(bugreport)
	app.Register(pullrequest)

	app.Run(os.Args)
}
