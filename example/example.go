package example

import (
	"fmt"
	"os"

	"github.com/MoonMoon1919/doyoucompute"
)

func main() {
	file, err := os.Open("../testdocs/partial.md")
	if err != nil {
		panic(err)
	}

	section := doyoucompute.Section{
		Name: "Intro",
		Content: []doyoucompute.Node{
			doyoucompute.Paragraph{
				Items: []doyoucompute.Node{
					doyoucompute.Text("cool text bro"),
				},
			},
			doyoucompute.Remote{Reader: file},
			doyoucompute.Executable{
				Shell: "sh",
				Cmd:   []string{"echo", "hello", "world"},
			},
		},
	}

	fmt.Print(section)

	// rendered, _ := section.Materialize()

	// fmt.Print(rendered)
}
