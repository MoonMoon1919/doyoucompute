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

	readme := documents.Readme()
	contribution := documents.Contributing()
	bugreport := documents.BugReport()

	app := app.New([]*doyoucompute.Document{&readme, &contribution, &bugreport}, &svc)

	app.Run(os.Args)
}
