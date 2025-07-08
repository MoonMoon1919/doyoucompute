package main

import (
	"fmt"
	"os"

	"github.com/MoonMoon1919/doyoucompute"
	"github.com/MoonMoon1919/doyoucompute/pkg/content"
)

func main() {
	file, err := os.Open("../testdocs/partial.md")
	if err != nil {
		panic(err)
	}

	section := doyoucompute.Section{
		Name: "Intro",
		Content: []content.Materializer{
			content.Paragraph("cool text bro"),
			content.Remote{Reader: file},
			content.Executable{
				Shell: "sh",
				Cmd:   []string{"echo", "hello", "world"},
			},
		},
	}

	rendered, _ := section.Materialize()

	fmt.Print(rendered)
}
