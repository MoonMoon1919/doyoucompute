package main

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
					doyoucompute.Code("very cool code"),
					doyoucompute.Link{
						Text: "Some Link",
						Url:  "https://example.com",
					},
				},
			},
			doyoucompute.Header{Content: "Things"},
			doyoucompute.BlockQuote("Here i am blockin' on my own"),
			doyoucompute.Remote{Reader: file},
			doyoucompute.Executable{
				Shell: "sh",
				Cmd:   []string{"echo", "hello", "world"},
			},
			doyoucompute.CodeBlock{
				BlockType: "json",
				Cmd:       []string{`{"key": "value"}`},
			},
			doyoucompute.Table{
				Headers: []string{"my", "cool", "table"},
				Items: []doyoucompute.Node{
					doyoucompute.TableRow{Values: []string{"some", "cool", "content"}},
					doyoucompute.TableRow{Values: []string{"more", "nice", "stuff"}},
					doyoucompute.TableRow{Values: []string{"very", "great", "table"}},
				},
			},
			doyoucompute.List{
				TypeOfList: doyoucompute.NUMBERED,
				Items: []doyoucompute.Node{
					doyoucompute.Text("first item"),
					doyoucompute.Code("npm i"),
				},
			},
		},
	}

	renderer := doyoucompute.Markdown{}

	rendered, _ := renderer.Render(section)

	fmt.Print(rendered)
}
