package main

import (
	"fmt"
	"os"

	"github.com/MoonMoon1919/doyoucompute"
)

func manualRoute() doyoucompute.Document {
	file, err := os.Open("../testdocs/partial.md")
	if err != nil {
		panic(err)
	}

	document := doyoucompute.Document{
		Name: "MY DOC",
		Content: []doyoucompute.Node{
			doyoucompute.Paragraph{
				Items: []doyoucompute.Node{
					doyoucompute.Text("I am an introduction paragraph"),
				},
			},
			doyoucompute.Section{
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
						Items: []doyoucompute.TableRow{
							{Values: []string{"some", "cool", "content"}},
							{Values: []string{"more", "nice", "stuff"}},
							{Values: []string{"very", "great", "table"}},
						},
					},
					doyoucompute.List{
						TypeOfList: doyoucompute.NUMBERED,
						Items: []doyoucompute.Text{
							doyoucompute.Text("first item"),
							doyoucompute.Text("second item"),
						},
					},
				},
			},
		},
	}

	return document
}

func builderRoute() doyoucompute.Document {
	file, err := os.Open("../testdocs/partial.md")
	if err != nil {
		panic(err)
	}

	document := doyoucompute.NewDocument("MY DOC")
	document.NewIntroWriter().Text("I am an introduction paragraph")

	// Build the section
	section := document.NewSection("Intro")
	section.NewIntroWriter().Text("cool text bro").Code("very cool code").Link("Some Link", "https://example.com")

	section.WriteBlockQuote("Here i am blockin' on my own")
	section.WriteRemoteContent(doyoucompute.Remote{Reader: file})
	section.WriteCodeBlock("sh", []string{"echo", "hello", "world"}, doyoucompute.Exec)
	section.WriteCodeBlock("json", []string{`{"key": "value"}`}, doyoucompute.Static)

	// Table
	table := section.NewTable(
		[]string{"my", "cool", "table"},
	)
	table.AddRow(doyoucompute.TableRow{Values: []string{"some", "cool", "content"}})
	table.AddRow(doyoucompute.TableRow{Values: []string{"more", "nice", "stuff"}})
	table.AddRow(doyoucompute.TableRow{Values: []string{"very", "great", "table"}})

	// List
	list := section.NewList(doyoucompute.NUMBERED)
	list.Append("first item")
	list.Append("second item")

	return document
}

func main() {
	renderer := doyoucompute.Markdown{}
	document := builderRoute()
	rendered, _ := renderer.Render(document)

	manualRoute() // So go stops yelling at me

	fmt.Print(rendered)
}
