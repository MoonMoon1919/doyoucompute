package samples

import (
	"os"

	"github.com/MoonMoon1919/doyoucompute"
	"github.com/MoonMoon1919/doyoucompute/pkg/app"
)

func main() {
	// Create service with file repository and task runner
	repo := doyoucompute.NewFileRepository()
	runner := doyoucompute.NewTaskRunner(doyoucompute.DefaultSecureConfig())
	markdownRenderer := doyoucompute.NewMarkdownRenderer()
	executionRenderer := doyoucompute.NewExecutionRenderer()

	service := doyoucompute.NewService(repo, runner, markdownRenderer, executionRenderer)

	// Create and run CLI app
	app := app.New(&service)

	doc, err := doyoucompute.NewDocument("My Project")
	if err != nil {
		panic(err)
	}

	app.Register(doc)

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
