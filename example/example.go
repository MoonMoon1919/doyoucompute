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
	docIntro := doyoucompute.NewParagraph().Next(doyoucompute.Text("I am an introduction paragraph"))

	document.AddIntro(docIntro)

	// Build the section
	section := doyoucompute.NewSection("Intro")
	intro := doyoucompute.NewParagraph().Next(doyoucompute.Text("cool text bro")).Next(doyoucompute.Code("very cool code")).Next(doyoucompute.Link{
		Text: "Some Link",
		Url:  "https://example.com",
	})

	section.AddIntro(intro)
	section.AddBlockQuote("Here i am blockin' on my own")
	section.AddRemoteContent(doyoucompute.Remote{Reader: file})
	section.AddCodeBlock("sh", []string{"echo", "hello", "world"}, true)
	section.AddCodeBlock("json", []string{`{"key": "value"}`}, false)
	section.AddTable(
		[]string{"my", "cool", "table"},
		[]doyoucompute.Node{
			doyoucompute.TableRow{Values: []string{"some", "cool", "content"}},
			doyoucompute.TableRow{Values: []string{"more", "nice", "stuff"}},
			doyoucompute.TableRow{Values: []string{"very", "great", "table"}},
		},
	)
	section.AddList(doyoucompute.NUMBERED, []doyoucompute.Node{
		doyoucompute.Text("first item"),
		doyoucompute.Code("npm i"),
	})

	document.AddSection(section)

	return document
}

func main() {
	renderer := doyoucompute.Markdown{}
	document := builderRoute()
	rendered, _ := renderer.Render(document)

	manualRoute() // So go stops yelling at me

	fmt.Print(rendered)
}
