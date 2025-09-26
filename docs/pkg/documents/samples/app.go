package samples

import (
	"os"

	"github.com/MoonMoon1919/doyoucompute"
	"github.com/MoonMoon1919/doyoucompute/pkg/app"
)

func main() {
	service, err := doyoucompute.DefaultService()
	if err != nil {
		panic(err)
	}

	// Create and run CLI app
	app := app.New(service)

	doc, err := doyoucompute.NewDocument("My Project")
	if err != nil {
		panic(err)
	}

	app.Register(doc)

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
